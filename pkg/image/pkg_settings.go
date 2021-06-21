package image

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type pkgSettings struct {
	Name         string
	PackageName  string
	LocalPackage bool
	ConfigPath   string
	Arguments    []string
	ProviderType string
}

func newPkgSettings(d *schema.ResourceData) *pkgSettings {
	st := &pkgSettings{
		Name:         d.Get("name").(string),
		PackageName:  d.Get("package_name").(string),
		LocalPackage: d.Get("local_package").(bool),
		ConfigPath:   d.Get("config").(string),
		ProviderType: d.Get("targetcloud").(string),
	}

	args, ok := d.Get("arguments").([]interface{})
	if ok {
		for _, a := range args {
			st.Arguments = append(st.Arguments, fmt.Sprintf("%v", a))
		}
	}
	return st
}
