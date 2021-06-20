package image

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nanovms/ops/cmd"
	"github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/types"
	"github.com/nanovms/terraform-provider-ops/ops"
	"github.com/nanovms/terraform-provider-ops/pkg/file"
	"github.com/spf13/pflag"
)

func pkgCreateImage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	settings := newPkgSettings(d)
	d.SetId(settings.Name)

	imagePath, err := pkgBuildImage(settings)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to build image: %v", err))
	}
	d.Set("path", imagePath)

	if settings.ConfigPath != "" {
		configChecksum, err := file.Checksum(settings.ConfigPath)
		if err != nil {
			return diag.Errorf("failed to generate config checksum: %v", err)
		}
		d.Set("config_checksum", configChecksum)
	}

	return diags
}

func pkgBuildImage(st *pkgSettings) (string, error) {
	opsCurrentVersion, err := ops.CurrentVersion()
	if err != nil {
		log.Println("[DEBUG] buildImage: ", err)
		return "", err
	}

	config := &types.Config{}
	if st.ConfigPath != "" {
		if _, err = os.Stat(st.ConfigPath); os.IsNotExist(err) {
			return "", fmt.Errorf("config file with path %s not found", st.ConfigPath)
		}

		config, err = ops.ReadConfigFromFile(st.ConfigPath)
		if err != nil {
			return "", fmt.Errorf("failed to read configuration: %v", err)
		}
	}

	flagset := pflag.NewFlagSet("loadPackage", pflag.ContinueOnError)
	flagset.Bool("local", st.LocalPackage, "")
	flagset.String("package", st.PackageName, "")
	flags := cmd.NewPkgCommandFlags(flagset)
	if err := flags.MergeToConfig(config); err != nil {
		return "", err
	}

	opsHome := lepton.GetOpsHome()

	config.RunConfig.Accel = true
	config.RunConfig.Memory = "2G"
	config.RunConfig.Imagename = path.Join(opsHome, "images", st.Name+".img")
	config.CloudConfig.ImageName = st.Name
	config.Kernel = path.Join(opsHome, opsCurrentVersion, "kernel.img")
	config.Boot = path.Join(opsHome, opsCurrentVersion, "boot.img")
	config.UefiBoot = lepton.GetUefiBoot(opsCurrentVersion)
	config.Args = st.Arguments

	if _, err := lepton.DownloadPackage(st.PackageName, config); err != nil {
		return "", err
	}

	provider, err := ops.ProviderByType(st.ProviderType, config)
	if err != nil {
		return "", fmt.Errorf("failed getting provider: %v", err)
	}

	pkgPath := filepath.Join(opsHome, "packages", st.PackageName)
	imagePath, err := provider.BuildImageWithPackage(lepton.NewContext(config), pkgPath)
	if err != nil {
		return "", fmt.Errorf("failed creating image: %v", err)
	}
	return imagePath, nil
}

func pkgReadImage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	configChecksum := d.Get("config_checksum").(string)
	d.Set("detect_config_checksum", configChecksum)

	return diags
}

func pkgDeleteImage(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	settings := newPkgSettings(d)
	provider, err := ops.ProviderByType(settings.ProviderType, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	opsCtx := lepton.NewContext(&types.Config{})
	if err = provider.DeleteImage(opsCtx, settings.Name+".img"); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diags
}
