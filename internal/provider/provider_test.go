package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
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

func TestProvider_HasChildResources(t *testing.T) {
	expectedResources := []string{}

	resources := New("dev")().ResourcesMap
	require.Equal(t, len(expectedResources), len(resources), "There are an unexpected number of registered resources")

	for _, resource := range expectedResources {
		require.Contains(t, resources, resource, "An expected resource was not registered")
		require.NotNil(t, resources[resource], "A resource cannot have a nil schema")
	}
}

func TestProvider_HasChildDataSources(t *testing.T) {
	expectedDataSources := []string{}

	dataSources := New("dev")().DataSourcesMap
	require.Equal(t, len(expectedDataSources), len(dataSources), "There are an unexpected number of registered data sources")

	for _, resource := range expectedDataSources {
		require.Contains(t, dataSources, resource, "An expected data source was not registered")
		require.NotNil(t, dataSources[resource], "A data source cannot have a nil schema")
	}
}

func TestProvider_SchemaIsValid(t *testing.T) {
	type testParams struct {
		name          string
		required      bool
		defaultEnvVar string
		sensitive     bool
	}

	tests := []testParams{
		{"org_service_url", false, "AZDO_ORG_SERVICE_URL", false},
		{"personal_access_token", false, "AZDO_PERSONAL_ACCESS_TOKEN", true},
	}

	schema := New("dev")().Schema
	require.Equal(t, len(tests), len(schema), "There are an unexpected number of properties in the schema")

	for _, test := range tests {
		require.Contains(t, schema, test.name, "An expected property was not found in the schema")
		require.NotNil(t, schema[test.name], "A property in the schema cannot have a nil value")
		require.Equal(t, test.sensitive, schema[test.name].Sensitive, "A property in the schema has an incorrect sensitivity value")
		require.Equal(t, test.required, schema[test.name].Required, "A property in the schema has an incorrect required value")

		if test.defaultEnvVar != "" {
			expectedValue := os.Getenv(test.defaultEnvVar)

			actualValue, err := schema[test.name].DefaultFunc()
			if actualValue == nil {
				actualValue = ""
			}

			require.Nil(t, err, "An error occurred when getting the default value from the environment")
			require.Equal(t, expectedValue, actualValue, "The default value pulled from the environment has the wrong value")
		}
	}
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv(envOrgServiceUrl); err == "" {
		t.Fatal(envOrgServiceUrl + " must be set for acceptance tests")
	}
	if err := os.Getenv(envPersonalAccessToken); err == "" {
		t.Fatal(envPersonalAccessToken + " must be set for acceptance tests")
	}
}
