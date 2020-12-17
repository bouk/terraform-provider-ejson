package ejson

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func Test_dataSourceEjsonFile(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"ejson": func() (*schema.Provider, error) {
				return Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
data "ejson_file" "config" {
  file = "../examples/secrets.ejson"
  private_key = "b9d31c0c787e56f1d29b60e813b388abab496808731616cc1a16887fa7179ca0"
}
`,
				Check: func(state *terraform.State) error {
					attr := state.RootModule().Resources["data.ejson_file.config"].Primary.Attributes
					require.Equal(t, attr["map.hi"], "what")
					return nil
				},
			},
		},
	})
}
