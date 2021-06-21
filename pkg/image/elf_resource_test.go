package image

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/types"
	"github.com/nanovms/terraform-provider-ops/pkg/testutil"
	"gotest.tools/assert"
)

func TestExecCreateImage(t *testing.T) {
	dataSource := NewFromExecutable()
	instanceState := &terraform.InstanceState{}
	data := dataSource.Data(instanceState)

	elfPath := testutil.BuildBasicProgram()
	defer os.Remove(elfPath)

	configPath := testutil.WriteConfigFile(&types.Config{})
	defer os.Remove(configPath)

	data.Set("name", "lorem")
	data.Set("elf", elfPath)
	data.Set("config", configPath)

	diags := execCreateImage(context.TODO(), data, nil)

	fmt.Println(diags)

	assert.Assert(t, !diags.HasError())

	imagePath := data.Get("path")

	assert.Equal(t, imagePath, path.Join(lepton.GetOpsHome(), "images", "lorem.img"))
}

func TestExecReadImage(t *testing.T) {
	dataSource := NewFromExecutable()
	instanceState := &terraform.InstanceState{}
	data := dataSource.Data(instanceState)

	data.Set("path", "/home/xyz/.ops/images/lorem.img")

	diags := execReadImage(context.TODO(), data, nil)

	assert.Assert(t, !diags.HasError())
}
