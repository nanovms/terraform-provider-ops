package ops

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/onprem"
	"github.com/nanovms/ops/provider"
	"github.com/nanovms/ops/types"
)

func resourceImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImageCreate,
		ReadContext:   resourceImageRead,
		UpdateContext: resourceImageUpdate,
		DeleteContext: resourceImageDelete,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"elf": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config": {
				Type:     schema.TypeString,
				Required: true,
			},
			"targetcloud": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceImageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	elfPath := d.Get("elf").(string)
	name := d.Get("name").(string)
	configPath := d.Get("config").(string)
	providerType := d.Get("targetcloud").(string)

	if providerType == "" {
		providerType = "onprem"
	}

	if _, err := os.Stat(elfPath); os.IsNotExist(err) {
		return diag.FromErr(fmt.Errorf("elf file with path %s not found", elfPath))
	}

	opsCurrentVersion, _ := getCurrentVersion()
	var config *types.Config
	var err error

	if configPath != "" {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return diag.Errorf("config file with path %s not found", configPath)
		}

		config, err = readConfigFromFile(configPath)
		if err != nil {
			return diag.Errorf("failed reading configuration: %v", err)
		}

	} else {
		config = &types.Config{}
	}

	config.Program = elfPath
	config.RunConfig.Accel = true
	config.RunConfig.Memory = "2G"
	config.RunConfig.Imagename = path.Join(lepton.GetOpsHome(), "images", name+".img")
	config.CloudConfig.ImageName = name
	config.Kernel = path.Join(lepton.GetOpsHome(), opsCurrentVersion, "kernel.img")
	config.Boot = path.Join(lepton.GetOpsHome(), opsCurrentVersion, "boot.img")
	config.UefiBoot = lepton.GetUefiBoot(opsCurrentVersion)

	provider, err := provider.CloudProvider(providerType, &config.CloudConfig)
	if err != nil {
		return diag.Errorf("failed getting provider: %v", err)
	}

	opsContext := lepton.NewContext(config)

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Creating images resource",
		Detail:   fmt.Sprintf("creating image %v", name),
	})

	imagePath, err := provider.BuildImage(opsContext)
	if err != nil {
		return diag.Errorf("failed creating image: %v", err)
	}

	d.Set("path", imagePath)

	d.SetId(imagePath)

	return diags
}

func resourceImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	path := d.Get("path").(string)

	provider := onprem.OnPrem{}
	opsContext := lepton.NewContext(&types.Config{})

	images, err := provider.GetImages(opsContext)
	if err != nil {
		return diag.Errorf("failed getting images: %v", err)
	}

	var image *lepton.CloudImage

	for _, i := range images {
		if i.Path == path {
			image = &i
			break
		}
	}

	if image != nil {
		nameParts := strings.Split(image.Name, ".")
		d.Set("name", nameParts[0])
	}

	return diags
}

func resourceImageUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Updating images resource",
		Detail:   fmt.Sprintf("Provider will update resource with images %v", d.Get("images")),
	})

	if d.HasChanges("elf", "config") {
		createDiags := resourceImageCreate(ctx, d, m)

		diags = append(diags, createDiags...)
	}

	return diags

}

func resourceImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	imageName := d.Get("name").(string) + ".img"

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Deleting images resource",
		Detail:   fmt.Sprintf("Provider will delete image with name %s", imageName),
	})

	provider := onprem.OnPrem{}
	opsContext := lepton.NewContext(&types.Config{})

	err := provider.DeleteImage(opsContext, imageName)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
