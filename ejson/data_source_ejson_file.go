package ejson

import (
	"context"
	"encoding/json"

	"github.com/Shopify/ejson"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEjsonFile() *schema.Resource {
	return &schema.Resource{
		Description: `Decrypt an ejson file and access its contents.`,
		Schema: map[string]*schema.Schema{
			"file": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `ejson file to decrypt.`,
			},
			"private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `Private key to use for decryption. The provider-level config keydir is used to find a key by default.`,
			},
			"data": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: `Decrypted contents of ejson file. Use jsondecode to get an object from the JSON blob.`,
			},
			"map": {
				Type:        schema.TypeMap,
				Computed:    true,
				Sensitive:   true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Mapping of decrypted keys to values, only top-level string values are included. The public key is stripped out and underscore prefixes are removed.`,
			},
		},
		ReadContext: dataSourceEjsonFileRead,
	}
}

func dataSourceEjsonFileRead(_ context.Context, file *schema.ResourceData, i interface{}) diag.Diagnostics {
	config := i.(*Config)
	filePath := file.Get("file").(string)
	file.SetId(filePath)
	privateKey := file.Get("private_key").(string)
	data, err := ejson.DecryptFile(filePath, config.Keydir, privateKey)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
				Detail:   "While trying to decrypt ejson file.",
			},
		}
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return diag.FromErr(err)
	}
	mapping := make(map[string]interface{}, len(m))
	for k, v := range m {
		if len(k) > 0 && k[0] == '_' {
			if k == "_public_key" {
				continue
			}
			k = k[1:]
		}
		if str, ok := v.(string); ok {
			mapping[k] = str
		}
	}
	if err := file.Set("data", string(data)); err != nil {
		return diag.FromErr(err)
	}
	if err := file.Set("map", mapping); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
