package image

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/types"
	"github.com/nanovms/terraform-provider-ops/pkg/testutil"
	"gotest.tools/assert"
)

func TestImageFromPackage(t *testing.T) {
	dataSource := NewFromPackage()
	instanceState := &terraform.InstanceState{}
	data := dataSource.Data(instanceState)

	data.Set("name", "test-pkg")
	data.Set("package_name", "node_v14.2.0")
	data.Set("arguments", []string{"hello.js"})

	configPath := testutil.WriteConfigFile(&types.Config{})
	defer os.Remove(configPath)
	data.Set("config", configPath)

	diags := pkgCreateImage(context.TODO(), data, nil)

	assert.Assert(t, !diags.HasError())
	imagePath := data.Get("path")
	assert.Equal(t, imagePath, path.Join(lepton.GetOpsHome(), "images", "test-pkg.img"))
}
