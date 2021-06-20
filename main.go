package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/nanovms/terraform-provider-ops/pkg/image"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return &schema.Provider{
				ResourcesMap: map[string]*schema.Resource{
					"ops_images":        image.ResourceImage(),
					"ops_image_package": image.NewFromPackage(),
				},
				DataSourcesMap: map[string]*schema.Resource{},
			}
		},
	})
}
