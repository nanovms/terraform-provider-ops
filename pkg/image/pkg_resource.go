package image

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nanovms/terraform-provider-ops/pkg/file"
)

// NewFromPackage creates image from a package.
func NewFromPackage() *schema.Resource {
	return &schema.Resource{
		CreateContext: pkgCreateImage,
		ReadContext:   pkgReadImage,
		DeleteContext: pkgDeleteImage,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"package_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"local_package": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"arguments": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"config": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"targetcloud": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
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
						currentChecksum, _ = file.Checksum(elf.(string))
					}

					if currentChecksum == "" || currentChecksum != old {
						return false
					}

					return true
				},
			},
		},
	}
}
