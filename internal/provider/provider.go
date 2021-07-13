package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Xtansia/terraform-provider-azdoext/internal/client"
	"github.com/Xtansia/terraform-provider-azdoext/internal/utils"
)

const (
	argOrgServiceUrl       = "org_service_url"
	envOrgServiceUrl       = "AZDO_ORG_SERVICE_URL"
	argPersonalAccessToken = "personal_access_token"
	envPersonalAccessToken = "AZDO_PERSONAL_ACCESS_TOKEN"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description

		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}

		if s.ConflictsWith != nil && len(s.ConflictsWith) > 0 {
			conflictFields := utils.Map(s.ConflictsWith, func(c string) string {
				return fmt.Sprintf("**%s**", c)
			})
			desc += fmt.Sprintf(" Conflicts with %s.", utils.HumaniseList(conflictFields))
		}

		return strings.TrimSpace(desc)
	}
}

func New(version string) func() *schema.Provider {
	rn := func(resourceName string) string {
		return fmt.Sprintf("azdoext_%s", resourceName)
	}

	return func() *schema.Provider {
		p := &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{},
			ResourcesMap: map[string]*schema.Resource{
				rn("secure_file"): resourceSecureFile(),
			},
			Schema: map[string]*schema.Schema{
				argOrgServiceUrl: {
					Description: "The url of the Azure DevOps instance which should be used. Can also be set via the `" + envOrgServiceUrl + "` environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc(envOrgServiceUrl, nil),
				},
				argPersonalAccessToken: {
					Description: "The personal access token which should be used. Can also be set via the `" + envPersonalAccessToken + "` environment variable.",
					Type:        schema.TypeString,
					Sensitive:   true,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc(envPersonalAccessToken, nil),
				},
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		orgServiceUrl := d.Get(argOrgServiceUrl).(string)
		personalAccessToken := d.Get(argPersonalAccessToken).(string)

		var diags diag.Diagnostics

		if strings.EqualFold(orgServiceUrl, "") {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Organisation service URL not set",
				Detail:   "The Azure DevOps organisation service URL must be set",
			})
		}
		if strings.EqualFold(personalAccessToken, "") {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Personal access token not set",
				Detail:   "The Azure DevOps personal access token must be set",
			})
		}

		if len(diags) > 0 {
			return nil, diags
		}

		options := client.Options{
			OrganisationUrl:     orgServiceUrl,
			PersonalAccessToken: personalAccessToken,
			ProviderVersion:     version,
			TerraformVersion:    p.TerraformVersion,
		}

		clients, err := options.Clients(ctx)

		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error initialising Azure DevOps clients",
			})
			return nil, diags
		}

		return clients, nil
	}
}
