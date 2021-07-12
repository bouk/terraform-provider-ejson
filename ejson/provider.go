package ejson

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"keydir": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("EJSON_KEYDIR", "/opt/ejson/keys"),
				Description: `Directory to read private keys from. Defaults to $EJSON_KEYDIR or /opt/ejson/keys if not set.`,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ejson_keypair": resourceEjsonKeypair(),
			"ejson_file":    resourceEjsonFile(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ejson_file": dataSourceEjsonFile(),
		},
		ConfigureContextFunc: func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
			keydir := data.Get("keydir").(string)
			return &Config{
				Keydir: keydir,
			}, nil
		},
	}
}

// Config for the provider
type Config struct {
	Keydir string
}
