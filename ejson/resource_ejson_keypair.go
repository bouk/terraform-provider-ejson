package ejson

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/Shopify/ejson"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEjsonKeypair() *schema.Resource {
	return &schema.Resource{
		Description: `Generate an ejson keypair.`,
		Schema: map[string]*schema.Schema{
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The public part of the key pair, for embedding into the ejson file.`,
			},
			"private_key": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Computed:    true,
				Description: `The private part of the key pair, for decrypting the ejson file.`,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: resourceEjsonKeypairImport,
		},
		ReadContext: schema.NoopContext,
		DeleteContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			data.SetId("")
			return nil
		},
		CreateContext: resourceEjsonKeypairCreate,
	}
}

func resourceEjsonKeypairCreate(_ context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	pub, priv, err := ejson.GenerateKeypair()
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(pub)
	if err := data.Set("public_key", pub); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("private_key", priv); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceEjsonKeypairImport(_ context.Context, data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	config := i.(*Config)
	pub := data.Id()
	keyPath := filepath.Join(config.Keydir, pub)
	private, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	data.SetId(pub)
	priv := string(private)
	priv = strings.TrimSpace(priv)

	// Verify that we can encrypt and decrypt a message with these keys
	var b bytes.Buffer
	if _, err := ejson.Encrypt(strings.NewReader(fmt.Sprintf(`{"_public_key":%q, "verify": "keys"}`, pub)), &b); err != nil {
		return nil, fmt.Errorf("verifying public key: %w", err)
	}
	if err := ejson.Decrypt(&b, ioutil.Discard, "", priv); err != nil {
		return nil, fmt.Errorf("verifying keypair: %w", err)
	}
	if err := data.Set("public_key", pub); err != nil {
		return nil, err
	}
	if err := data.Set("private_key", priv); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{data}, nil
}
