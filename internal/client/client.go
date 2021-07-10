package client

import (
    "context"
    "fmt"
    "log"
    "strings"

    "github.com/microsoft/azure-devops-go-api/azuredevops"

    "github.com/Xtansia/terraform-provider-azdoext/internal/client/taskagent"
)

type ClientOptions struct {
    OrganisationUrl string
    PersonalAccessToken string
    ProviderVersion string
    TerraformVersion string
}

type Clients struct {
    TaskAgentClient taskagent.Client
}

func (o *ClientOptions) Clients(ctx context.Context) (*Clients, error) {
    if strings.EqualFold(o.OrganisationUrl, "") {
        return nil, fmt.Errorf("url of the Azure DevOps is required")
    }
    if strings.EqualFold(o.PersonalAccessToken, "") {
        return nil, fmt.Errorf("personal access token is required")
    }

    connection := azuredevops.NewPatConnection(o.OrganisationUrl, o.PersonalAccessToken)
    o.setUserAgent(connection)

    taskAgentClient, err := taskagent.NewClient(ctx, connection)
    if err != nil {
        return nil, err
    }

    return &Clients{
        TaskAgentClient: taskAgentClient,
    }, nil
}

func (o *ClientOptions) setUserAgent(connection *azuredevops.Connection) {
    parts := []string {
        connection.UserAgent,
        fmt.Sprintf("terraform-provider-azdoext/%s (+https://registry.terraform.io/providers/Xtansia/azdoext)", o.ProviderVersion),
        fmt.Sprintf("Terraform/%s (+https://www.terraform.io)", o.TerraformVersion),
    }

    connection.UserAgent = strings.TrimSpace(strings.Join(parts, " "))

    log.Printf("[DEBUG] Azure DevOps Client User Agent: %s\n", connection.UserAgent)
}