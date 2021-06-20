package ops

import (
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/types"
	"gotest.tools/assert"
)

func TestResourceImageCreate(t *testing.T) {
	dataSource := ResourceImage()
	instanceState := &terraform.InstanceState{}
	data := dataSource.Data(instanceState)

	elfPath := BuildBasicProgram()
	defer os.Remove(elfPath)

	configPath := WriteConfigFile(&types.Config{})
	defer os.Remove(configPath)

	data.Set("name", "lorem")
	data.Set("elf", elfPath)
	data.Set("config", configPath)

	diags := resourceImageCreate(context.TODO(), data, nil)

	fmt.Println(diags)

	assert.Assert(t, !diags.HasError())

	imagePath := data.Get("path")

	assert.Equal(t, imagePath, path.Join(lepton.GetOpsHome(), "images", "lorem.img"))
}

func TestResourceImageRead(t *testing.T) {
	dataSource := ResourceImage()
	instanceState := &terraform.InstanceState{}
	data := dataSource.Data(instanceState)

	data.Set("path", "/home/xyz/.ops/images/lorem.img")

	diags := resourceImageRead(context.TODO(), data, nil)

	assert.Assert(t, !diags.HasError())
}
