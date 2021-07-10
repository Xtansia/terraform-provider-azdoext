package taskagent

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/url"
    "strconv"

    "github.com/google/uuid"
    "github.com/microsoft/azure-devops-go-api/azuredevops"
    "github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
)

const (
    ApiVersion61Preview1            = "6.1-preview.1"
    MediaTypeApplicationOctetStream = "application/octet-stream"
)

var (
    SecureFilesLocationId, _ = uuid.Parse("adcfd8bc-b184-43ba-bd84-7c8c6a2ff421")
)

type Client interface {
    DeleteSecureFile(context.Context, DeleteSecureFileArgs) error
    GetSecureFile(context.Context, GetSecureFileArgs) (*taskagent.SecureFile, error)
    UpdateSecureFile(context.Context, UpdateSecureFileArgs) (*taskagent.SecureFile, error)
    UploadSecureFile(context.Context, UploadSecureFileArgs) (*taskagent.SecureFile, error)
}

type ClientImpl taskagent.ClientImpl

func NewClient(ctx context.Context, connection *azuredevops.Connection) (Client, error) {
    internalClient, err := taskagent.NewClient(ctx, connection)
    if err != nil {
        return nil, err
    }
    clientImpl := ClientImpl(*internalClient.(*taskagent.ClientImpl))
    return &clientImpl, nil
}

func (client *ClientImpl) DeleteSecureFile(ctx context.Context, args DeleteSecureFileArgs) error {
    if args.Project == nil || *args.Project == "" {
        return &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
    }
    if args.SecureFileId == nil {
        return &azuredevops.ArgumentNilError{ArgumentName: "args.SecureFileId"}
    }

    routeValues := make(map[string]string)
    routeValues["project"] = *args.Project
    routeValues["secureFileId"] = (*args.SecureFileId).String()

    _, err := client.Client.Send(
        ctx, http.MethodDelete, SecureFilesLocationId, ApiVersion61Preview1, routeValues, nil, nil, "",
        azuredevops.MediaTypeApplicationJson, nil,
    )

    return err
}

type DeleteSecureFileArgs struct {
    Project      *string
    SecureFileId *uuid.UUID
}

func (client *ClientImpl) GetSecureFile(ctx context.Context, args GetSecureFileArgs) (*taskagent.SecureFile, error) {
    if args.Project == nil || *args.Project == "" {
        return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
    }
    if args.SecureFileId == nil {
        return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SecureFileId"}
    }

    routeValues := make(map[string]string)
    routeValues["project"] = *args.Project
    routeValues["secureFileId"] = (*args.SecureFileId).String()

    queryParams := url.Values{}
    if args.IncludeDownloadTicket != nil {
        queryParams.Add("includeDownloadTicket", strconv.FormatBool(*args.IncludeDownloadTicket))
    }
    if args.ActionFilter != nil {
        queryParams.Add("actionFilter", (string)(*args.ActionFilter))
    }

    resp, err := client.Client.Send(
        ctx, http.MethodGet, SecureFilesLocationId, ApiVersion61Preview1, routeValues, queryParams, nil, "",
        azuredevops.MediaTypeApplicationJson, nil,
    )

    if err != nil {
        return nil, err
    }

    var responseValue taskagent.SecureFile
    err = client.Client.UnmarshalBody(resp, &responseValue)
    return &responseValue, err
}

type GetSecureFileArgs struct {
    Project               *string
    SecureFileId          *uuid.UUID
    IncludeDownloadTicket *bool
    ActionFilter          *taskagent.SecureFileActionFilter
}

func (client *ClientImpl) UpdateSecureFile(ctx context.Context, args UpdateSecureFileArgs) (*taskagent.SecureFile, error) {
    if args.Project == nil || *args.Project == "" {
        return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
    }
    if args.SecureFileId == nil {
        return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SecureFileId"}
    }
    if args.SecureFile == nil {
        return nil, &azuredevops.ArgumentNilError{ArgumentName: "args.SecureFile"}
    }

    routeValues := make(map[string]string)
    routeValues["project"] = *args.Project
    routeValues["secureFileId"] = (*args.SecureFileId).String()

    body, err := json.Marshal(args.SecureFile)
    if err != nil {
        return nil, err
    }

    resp, err := client.Client.Send(
        ctx, http.MethodPatch, SecureFilesLocationId, ApiVersion61Preview1, routeValues, nil, bytes.NewReader(body), azuredevops.MediaTypeApplicationJson,
        azuredevops.MediaTypeApplicationJson, nil,
    )

    if err != nil {
        return nil, err
    }

    var responseValue taskagent.SecureFile
    err = client.Client.UnmarshalBody(resp, &responseValue)
    return &responseValue, err
}

type UpdateSecureFileArgs struct {
    Project *string
    SecureFileId *uuid.UUID
    SecureFile *taskagent.SecureFile
}

func (client *ClientImpl) UploadSecureFile(ctx context.Context, args UploadSecureFileArgs) (
    *taskagent.SecureFile, error,
) {
    if args.Name == nil || *args.Name == "" {
        return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Name"}
    }
    if args.Project == nil || *args.Project == "" {
        return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Project"}
    }
    if args.Content == nil || len(*args.Content) == 0 {
        return nil, &azuredevops.ArgumentNilOrEmptyError{ArgumentName: "args.Content"}
    }

    routeValues := make(map[string]string)
    routeValues["project"] = *args.Project

    queryParams := url.Values{}
    queryParams.Add("name", *args.Name)
    if args.AuthorizePipelines != nil {
        queryParams.Add("authorizePipelines", strconv.FormatBool(*args.AuthorizePipelines))
    }

    resp, err := client.Client.Send(
        ctx, http.MethodPost, SecureFilesLocationId, ApiVersion61Preview1, routeValues, queryParams,
        bytes.NewReader(*args.Content), MediaTypeApplicationOctetStream, azuredevops.MediaTypeApplicationJson, nil,
    )
    if err != nil {
        return nil, err
    }

    var responseValue taskagent.SecureFile
    err = client.Client.UnmarshalBody(resp, &responseValue)
    return &responseValue, err
}

type UploadSecureFileArgs struct {
    Name               *string
    Project            *string
    AuthorizePipelines *bool
    Content            *[]byte
}