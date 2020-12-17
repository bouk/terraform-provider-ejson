package ejson

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func Test_resourceEjsonKeypair(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"ejson": func() (*schema.Provider, error) {
				return Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
provider "ejson" {
  keydir = "../examples/keys"
}
resource "ejson_keypair" "key" {}
`,
				ResourceName:  `ejson_keypair.key`,
				ImportState:   true,
				ImportStateId: `965aed709d63e22fd7cb6b4ee8f317bb0ad99a07fa68bc24bdcc349d0b2af130`,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					state := states[0]
					require.Equal(t, state.ID, "965aed709d63e22fd7cb6b4ee8f317bb0ad99a07fa68bc24bdcc349d0b2af130")
					attr := state.Attributes
					require.Equal(t, attr["public_key"], "965aed709d63e22fd7cb6b4ee8f317bb0ad99a07fa68bc24bdcc349d0b2af130")
					require.Equal(t, attr["private_key"], "b9d31c0c787e56f1d29b60e813b388abab496808731616cc1a16887fa7179ca0")
					return nil
				},
			},
		},
	})
}
