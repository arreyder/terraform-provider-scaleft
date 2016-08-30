package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"scaleft_server": resourceServer(),
		},
		Schema: map[string]*schema.Schema{
			"key_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCALEFT_KEY_ID", nil),
				Description: "The key id for the scaleft service user",
			},
			"key_secret": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCALEFT_KEY_SECRET", nil),
				Description: "The key secret for the scaleft service user",
			},
			"team": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCALEFT_TEAM", nil),
				Description: "The team for the scaleft service user",
			},
			"project": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCALEFT_PROJECT", nil),
				Description: "The project for the scaleft servers",
			},
		},
	}
}
