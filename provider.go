package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDCENTER_USERNAME", nil),
				Description: "Username used to access Cisco Cloudcenter",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDCENTER_PASSWORD", nil),
				Description: "Password used to access Cisco Cloudcenter",
			},
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDCENTER_URL", nil),
				Description: "URL to the CloudCenter Manager",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cloudcenter_user": resourceUser(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := &Config{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Base_url: d.Get("base_url").(string),
	}

	return config.Client(), nil
}
