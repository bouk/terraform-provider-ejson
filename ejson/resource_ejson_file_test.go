package ejson

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Shopify/ejson"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

// ejsonPrivateKey is the private part of 797ff0645ca1d8177f8546b417e212c450b3c5c1ab0a9e947621d19348ef0946
const ejsonPrivateKey = `638ab6ad0df6f17e8cdcaa265e60e0371b8f1fb7dadec3f8e2a715ff3d7329fc`

func Test_resourceEjsonFile(t *testing.T) {
	var (
		step1Result map[string]interface{}
		step2Result map[string]interface{}
		step3Result map[string]interface{}
		step4Result map[string]interface{}
	)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"ejson": func() (*schema.Provider, error) {
				return Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
resource "ejson_file" "secrets" {
  public_key = "797ff0645ca1d8177f8546b417e212c450b3c5c1ab0a9e947621d19348ef0946"
  data = jsonencode({
    a = "Hello"
    b = "Goodbye"
  })
}
`,
				ResourceName: `ejson_file.secrets`,
				Check: func(state *terraform.State) error {
					t.Log(state)
					attrs := state.RootModule().Resources["ejson_file.secrets"].Primary.Attributes
					return json.Unmarshal([]byte(attrs["encrypted"]), &step1Result)
				},
			},
			{
				Config: `
resource "ejson_file" "secrets" {
  public_key = "797ff0645ca1d8177f8546b417e212c450b3c5c1ab0a9e947621d19348ef0946"
  data = jsonencode({
    a = "Hello"
    c = "Nice"
  })
}
`,
				ResourceName: `ejson_file.secrets`,
				Check: func(state *terraform.State) error {
					t.Log(state)
					attrs := state.RootModule().Resources["ejson_file.secrets"].Primary.Attributes
					if err := json.Unmarshal([]byte(attrs["encrypted"]), &step2Result); err != nil {
						return err
					}
					require.Equal(t, step1Result["a"], step2Result["a"])
					return nil
				},
			},
			{
				Config: `
resource "ejson_file" "secrets" {
  public_key = "797ff0645ca1d8177f8546b417e212c450b3c5c1ab0a9e947621d19348ef0946"
  data = jsonencode({
    a = "Hello"
    b = "Goodbye"
    c = "Nice"
  })
}
`,
				ResourceName: `ejson_file.secrets`,
				Check: func(state *terraform.State) error {
					t.Log(state)
					attrs := state.RootModule().Resources["ejson_file.secrets"].Primary.Attributes
					if err := json.Unmarshal([]byte(attrs["encrypted"]), &step3Result); err != nil {
						return err
					}
					require.Equal(t, step1Result["a"], step3Result["a"])
					require.NotEqual(t, step1Result["b"], step3Result["b"])
					require.Equal(t, step2Result["c"], step3Result["c"])
					return nil
				},
			},
			{ // no change
				Config: `
resource "ejson_file" "secrets" {
  public_key = "797ff0645ca1d8177f8546b417e212c450b3c5c1ab0a9e947621d19348ef0946"
  data = jsonencode({
    a = "Hello"
    b = "Goodbye"
    c = "Nice"

	array = ["wow", "cool"]
  })
}
`,
				ResourceName: `ejson_file.secrets`,
				Check: func(state *terraform.State) error {
					t.Log(state)
					attrs := state.RootModule().Resources["ejson_file.secrets"].Primary.Attributes
					if err := json.Unmarshal([]byte(attrs["encrypted"]), &step4Result); err != nil {
						return err
					}
					require.Equal(t, step3Result["a"], step4Result["a"])
					require.Equal(t, step3Result["b"], step4Result["b"])
					require.Equal(t, step3Result["c"], step4Result["c"])

					// Check it decodes properly
					var output bytes.Buffer
					if err := ejson.Decrypt(strings.NewReader(attrs["encrypted"]), &output, "", ejsonPrivateKey); err != nil {
						return err
					}
					var result map[string]interface{}
					if err := json.Unmarshal(output.Bytes(), &result); err != nil {
						return err
					}

					require.Equal(t, "Hello", result["a"])
					require.Equal(t, "Goodbye", result["b"])
					require.Equal(t, "Nice", result["c"])
					require.Equal(t, []interface{}{"wow", "cool"}, result["array"])

					return nil
				},
			},
		},
	})
}
