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
		Schema: map[string]*schema.Schema{
			"data": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
			"file": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ReadContext: dataSourceEjsonFileRead,
	}
}

func dataSourceEjsonFileRead(ctx context.Context, file *schema.ResourceData, i interface{}) diag.Diagnostics {
	config := i.(*Config)
	filePath := file.Get("file").(string)
	data, err := ejson.DecryptFile(filePath, config.Keydir, "")
	if err != nil {
		return diag.FromErr(err)
	}
	var d map[string]interface{}
	if err := json.Unmarshal(data, &d); err != nil {
		return diag.FromErr(err)
	}
	if err := file.Set("data", d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
