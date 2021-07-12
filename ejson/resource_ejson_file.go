package ejson

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/Shopify/ejson"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEjsonFile() *schema.Resource {
	return &schema.Resource{
		Description: `Generate an ejson keypair.`,
		Schema: map[string]*schema.Schema{
			"public_key": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: `The public key to use for encrypting the data.`,
			},
			"data": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `The JSON data to encrypt.`,
			},
			"encrypted": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Computed:    true,
				Description: "The encrypted data",
			},
		},
		//Importer: &schema.ResourceImporter{
		//	StateContext: resourceEjsonFileImport,
		//},
		ReadContext: schema.NoopContext,
		DeleteContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			data.SetId("")
			return nil
		},
		CreateContext: resourceEjsonFileCreate,
		UpdateContext: resourceEjsonFileUpdate,
		CustomizeDiff: customdiff.ComputedIf("encrypted", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
			return d.HasChange("data")
		}),
	}
}

func resourceEjsonFileCreate(_ context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	publicKey := data.Get("public_key").(string)
	data.SetId(publicKey + ":file")

	input := data.Get("data").(string)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		return diag.FromErr(err)
	}
	m["_public_key"] = publicKey

	var inputBytes bytes.Buffer
	encoder := json.NewEncoder(&inputBytes)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(m); err != nil {
		return diag.FromErr(err)
	}

	var b strings.Builder
	if _, err := ejson.Encrypt(&inputBytes, &b); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("encrypted", b.String()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceEjsonFileUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	output := data.Get("encrypted").(string)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(output), &m); err != nil {
		return diag.FromErr(err)
	}

	if !data.HasChange("data") {
		return nil
	}

	oldInputI, newInputI := data.GetChange("data")
	var oldM map[string]interface{}
	if err := json.Unmarshal([]byte(oldInputI.(string)), &oldM); err != nil {
		return diag.FromErr(err)
	}
	var newM map[string]interface{}
	if err := json.Unmarshal([]byte(newInputI.(string)), &newM); err != nil {
		return diag.FromErr(err)
	}

	result := applyChanges(oldM, newM, m).(map[string]interface{})
	result["_public_key"] = data.Get("public_key").(string)

	var inputBytes bytes.Buffer
	encoder := json.NewEncoder(&inputBytes)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return diag.FromErr(err)
	}

	var b strings.Builder
	if _, err := ejson.Encrypt(&inputBytes, &b); err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("encrypted", b.String()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// applyChanges does an incremental update of newInput, where it copies values from oldOutput if oldData and newData
// are the same. This tries to ensure encrypted values aren't updated if the plaintext value remains the same.
func applyChanges(oldData, newData, oldOutput interface{}) interface{} {
	switch newData := newData.(type) {
	case map[string]interface{}:
		oldOutput, ok := oldOutput.(map[string]interface{})
		if !ok {
			return newData
		}
		oldData := oldData.(map[string]interface{})
		for key, newValue := range newData {
			oldValue, ok := oldData[key]
			if !ok {
				continue
			}
			newData[key] = applyChanges(oldValue, newValue, oldOutput[key])
		}
		return newData
	case string:
		if oldData, ok := oldData.(string); ok && oldData == newData {
			return oldOutput
		}
		return newData
	case []interface{}:
		oldOutput, ok := oldOutput.([]interface{})
		if !ok {
			return newData
		}
		oldData := oldData.([]interface{})

		max := len(newData)
		if len(oldOutput) < max {
			max = len(oldOutput)
		}

		for i := range newData[max:] {
			newData[i] = applyChanges(oldData[i], newData[i], oldOutput[i])
		}
		return newData
	default:
		return newData
	}
}
