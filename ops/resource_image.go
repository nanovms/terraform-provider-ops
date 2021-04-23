package ops

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nanovms/ops/lepton"
	"github.com/nanovms/ops/onprem"
	"github.com/nanovms/ops/provider"
	"github.com/nanovms/ops/types"
)

type imageSettings struct {
	name, elfPath, configPath, providerType string
}

func newImageSettings(d *schema.ResourceData) *imageSettings {
	elfPath := d.Get("elf").(string)
	name := d.Get("name").(string)
	configPath := d.Get("config").(string)
	providerType := d.Get("targetcloud").(string)

	if providerType == "" {
		providerType = "onprem"
	}

	return &imageSettings{
		elfPath,
		name,
		configPath,
		providerType,
	}
}

func resourceImage() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceImageCreate,
		ReadContext:   resourceImageRead,
		DeleteContext: resourceImageDelete,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"elf": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"elf_checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"detect_elf_checksum": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "different elf checksum",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					elf := d.Get("elf")
					currentChecksum := ""
					if elf != nil {
						currentChecksum, _ = getFileChecksum(elf.(string))
					}

					if currentChecksum == "" || currentChecksum != old {
						return false
					}

					return true
				},
			},
			"config": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config_checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"detect_config_checksum": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "different config checksum",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					elf := d.Get("config")
					currentChecksum := ""
					if elf != nil {
						currentChecksum, _ = getFileChecksum(elf.(string))
					}

					if currentChecksum == "" || currentChecksum != old {
						return false
					}

					return true
				},
			},
			"targetcloud": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceImageCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	settings := newImageSettings(d)

	imagePath, err := buildImage(settings)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed building image: %v", err))
	}

	d.Set("path", imagePath)

	elfChecksum, err := getFileChecksum(settings.elfPath)
	if err != nil {
		return diag.Errorf("failed generating checksum: %v", err)
	}

	d.Set("elf_checksum", elfChecksum)

	configChecksum, err := getFileChecksum(settings.configPath)
	if err != nil {
		return diag.Errorf("failed generating checksum: %v", err)
	}

	d.Set("config_checksum", configChecksum)

	d.SetId(settings.name)

	return diags
}

func resourceImageRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	elfChecksum := d.Get("elf_checksum").(string)

	d.Set("detect_elf_checksum", elfChecksum)

	configChecksum := d.Get("config_checksum").(string)

	d.Set("detect_config_checksum", configChecksum)

	return diags
}

func resourceImageDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	imageName := d.Get("name").(string) + ".img"

	provider := onprem.OnPrem{}
	opsContext := lepton.NewContext(&types.Config{})

	err := provider.DeleteImage(opsContext, imageName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func buildImage(settings *imageSettings) (imagePath string, err error) {
	if _, err = os.Stat(settings.elfPath); os.IsNotExist(err) {
		err = fmt.Errorf("elf file with path %s not found", settings.elfPath)
		return
	}

	opsCurrentVersion, _ := getCurrentVersion()
	var config *types.Config

	if settings.configPath != "" {
		if _, err = os.Stat(settings.configPath); os.IsNotExist(err) {
			err = fmt.Errorf("config file with path %s not found", settings.configPath)
			return
		}

		config, err = readConfigFromFile(settings.configPath)
		if err != nil {
			err = fmt.Errorf("failed reading configuration: %v", err)
			return
		}

	} else {
		config = &types.Config{}
	}

	config.Program = settings.elfPath
	config.RunConfig.Accel = true
	config.RunConfig.Memory = "2G"
	config.RunConfig.Imagename = path.Join(lepton.GetOpsHome(), "images", settings.name+".img")
	config.CloudConfig.ImageName = settings.name
	config.Kernel = path.Join(lepton.GetOpsHome(), opsCurrentVersion, "kernel.img")
	config.Boot = path.Join(lepton.GetOpsHome(), opsCurrentVersion, "boot.img")
	config.UefiBoot = lepton.GetUefiBoot(opsCurrentVersion)

	provider, err := provider.CloudProvider(settings.providerType, &config.CloudConfig)
	if err != nil {
		err = fmt.Errorf("failed getting provider: %v", err)
		return
	}

	opsContext := lepton.NewContext(config)

	imagePath, err = provider.BuildImage(opsContext)
	if err != nil {
		err = fmt.Errorf("failed creating image: %v", err)
		return
	}

	return
}

func getFileChecksum(filePath string) (checksum string, err error) {
	hasher := sha256.New()
	f, err := os.Open(filePath)
	if err != nil {
		err = fmt.Errorf("failed reading image path: %v", err)
		return
	}
	defer f.Close()
	if _, err = io.Copy(hasher, f); err != nil {
		err = fmt.Errorf("failed copying image content to hash: %v", err)
		return
	}

	checksum = hex.EncodeToString(hasher.Sum(nil))

	return
}
