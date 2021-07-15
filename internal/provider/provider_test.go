package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"azdoext": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func preCheck(t *testing.T) {
	if err := os.Getenv(envOrgServiceUrl); err == "" {
		t.Fatal(envOrgServiceUrl + " must be set for acceptance tests")
	}
	if err := os.Getenv(envPersonalAccessToken); err == "" {
		t.Fatal(envPersonalAccessToken + " must be set for acceptance tests")
	}
}
