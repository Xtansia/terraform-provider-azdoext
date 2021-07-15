package provider

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"

	"github.com/Xtansia/terraform-provider-azdoext/internal/client"
	"github.com/Xtansia/terraform-provider-azdoext/internal/client/taskagent"
	"github.com/Xtansia/terraform-provider-azdoext/internal/utils"
)

const (
	sfProjectId     = "project_id"
	sfName          = "name"
	sfContent       = "content"
	sfContentBase64 = "content_base64"
	sfAllowAccess   = "allow_access"
	sfProperties    = "properties"
)

const (
	invalidSecureFileIdErrorMessageFormat = "Error parsing the secure file ID from the Terraform resource data: %v"
)

func resourceSecureFile() *schema.Resource {
	return &schema.Resource{
		Description: "Manages secure files within Azure DevOps.",

		CreateContext: resourceSecureFileCreate,
		ReadContext:   resourceSecureFileRead,
		UpdateContext: resourceSecureFileUpdate,
		DeleteContext: resourceSecureFileDelete,

		Schema: map[string]*schema.Schema{
			sfProjectId: {
				Description:  "The ID of the Azure DevOps project the secure file belongs to.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			sfName: {
				Description:  "The name of the secure file.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			sfContent: {
				Description:   "The plain-text content of the secure file. Use **" + sfContentBase64 + "** for binary content to avoid issues.",
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "",
				ForceNew:      true,
				Sensitive:     true,
				ConflictsWith: []string{sfContentBase64},
				StateFunc:     secureFileContentHash,
			},
			sfContentBase64: {
				Description:   "The base64 encoded content of the secure file.",
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "",
				ForceNew:      true,
				Sensitive:     true,
				ConflictsWith: []string{sfContent},
				ValidateFunc:  utils.StringIsBase64Encoded,
				StateFunc:     secureFileContentHash,
			},
			sfAllowAccess: {
				Description: "Whether to allow all pipelines access to this resource.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			sfProperties: {
				Description: "Properties assigned to the secure file.",
				Type:        schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"id": {
				Description: "The ID of this resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSecureFileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(*client.Clients)

	secureFile := expandSecureFile(d)
	projectId := d.Get(sfProjectId).(string)
	content := d.Get(sfContent).(string)
	contentBase64 := d.Get(sfContentBase64).(string)

	var data []byte
	if content != "" {
		data = []byte(content)
	} else {
		data, _ = base64.StdEncoding.DecodeString(contentBase64)
	}

	createdSecureFile, err := clients.TaskAgentClient.UploadSecureFile(
		ctx, taskagent.UploadSecureFileArgs{
			Project: &projectId,
			Name:    secureFile.Name,
			Content: &data,
		},
	)
	if err != nil {
		return diag.Errorf("Error creating secure file in Azure DevOps: %+v", err)
	}

	createdSecureFile, err = updateSecureFile(clients, ctx, &projectId, createdSecureFile.Id, secureFile)
	if err != nil {
		return diag.Errorf("Error updating properties on secure file in Azure DevOps: %+v", err)
	}

	flattenSecureFile(d, createdSecureFile, &projectId)

	definitionResources := expandAllowAccess(d, createdSecureFile)
	definitionResourceReferences, err := authorizeProjectReferences(clients, ctx, &projectId, definitionResources)
	if err != nil {
		return diag.Errorf("Error creating definitionResourceReference Azure DevOps object: %+v", err)
	}

	flattenAllowAccess(d, definitionResourceReferences)

	return resourceSecureFileRead(ctx, d, meta)
}

func resourceSecureFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(*client.Clients)

	secureFileId, projectId, err := parseSecureFileAndProjectIds(d)
	if err != nil {
		return diag.Errorf(invalidSecureFileIdErrorMessageFormat, err)
	}

	secureFile, err := clients.TaskAgentClient.GetSecureFile(
		ctx,
		taskagent.GetSecureFileArgs{
			Project:      projectId,
			SecureFileId: secureFileId,
		},
	)
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf(
			"Error looking up secure file given ID (%v) and project ID (%v): %v", secureFileId, projectId, err,
		)
	}
	if secureFile.Id == nil {
		d.SetId("")
		return nil
	}

	flattenSecureFile(d, secureFile, projectId)

	resourceRefType := "securefile"
	secFileId := secureFileId.String()

	projectResources, err := clients.BuildClient.GetProjectResources(
		ctx, build.GetProjectResourcesArgs{
			Project: projectId,
			Type:    &resourceRefType,
			Id:      &secFileId,
		},
	)
	if err != nil {
		return diag.Errorf(
			"Error looking up project resources given ID (%v) and project ID (%v): %v", secureFileId, projectId, err,
		)
	}

	flattenAllowAccess(d, projectResources)

	return nil
}

func resourceSecureFileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(*client.Clients)

	secureFileId, projectId, err := parseSecureFileAndProjectIds(d)
	if err != nil {
		return diag.Errorf(invalidSecureFileIdErrorMessageFormat, err)
	}

	secureFile := expandSecureFile(d)

	updatedSecureFile, err := updateSecureFile(clients, ctx, projectId, secureFileId, secureFile)
	if err != nil {
		return diag.Errorf("Error updating secure file in Azure DevOps: %+v", err)
	}

	flattenSecureFile(d, updatedSecureFile, projectId)

	definitionResources := expandAllowAccess(d, updatedSecureFile)
	definitionResourceReferences, err := authorizeProjectReferences(clients, ctx, projectId, definitionResources)
	if err != nil {
		return diag.Errorf("Error creating definitionResourceReference Azure DevOps object: %+v", err)
	}

	flattenAllowAccess(d, definitionResourceReferences)

	return resourceSecureFileRead(ctx, d, meta)
}

func resourceSecureFileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(*client.Clients)

	secureFileId, projectId, err := parseSecureFileAndProjectIds(d)
	if err != nil {
		return diag.Errorf(invalidSecureFileIdErrorMessageFormat, err)
	}

	resourceRefType := "securefile"
	secFileId := secureFileId.String()
	resourceRefName := ""
	authorized := false
	_, err = authorizeProjectReferences(
		clients, ctx, projectId, []build.DefinitionResourceReference{
			{
				Type:       &resourceRefType,
				Id:         &secFileId,
				Name:       &resourceRefName,
				Authorized: &authorized,
			},
		},
	)
	if err != nil {
		return diag.Errorf(
			"Error deleting the allow access definitionResource for secure file ID (%v) and project ID (%v): %v",
			secureFileId, projectId, err,
		)
	}

	err = clients.TaskAgentClient.DeleteSecureFile(
		ctx,
		taskagent.DeleteSecureFileArgs{
			Project:      projectId,
			SecureFileId: secureFileId,
		},
	)
	if err != nil {
		return diag.Errorf("Error deleting secure file in Azure DevOps: %+v", err)
	}

	return nil
}

func expandSecureFile(d *schema.ResourceData) taskagent.SecureFile {
	name := d.Get(sfName).(string)
	properties := map[string]string{}
	for k, v := range d.Get(sfProperties).(map[string]interface{}) {
		properties[k] = v.(string)
	}
	return taskagent.SecureFile{
		Name:       &name,
		Properties: &properties,
	}
}

func flattenSecureFile(d *schema.ResourceData, secureFile *taskagent.SecureFile, projectId *string) {
	d.SetId(secureFile.Id.String())
	_ = d.Set(sfName, *secureFile.Name)
	_ = d.Set(sfProjectId, projectId)
	_ = d.Set(sfProperties, secureFile.Properties)
}

func updateSecureFile(
	clients *client.Clients, ctx context.Context, projectId *string, secureFileId *uuid.UUID,
	secureFile taskagent.SecureFile,
) (*taskagent.SecureFile, error) {
	return clients.TaskAgentClient.UpdateSecureFile(
		ctx, taskagent.UpdateSecureFileArgs{
			Project:      projectId,
			SecureFileId: secureFileId,
			SecureFile:   &secureFile,
		},
	)
}

func parseSecureFileAndProjectIds(d *schema.ResourceData) (*uuid.UUID, *string, error) {
	secureFileId, err := uuid.Parse(d.Id())
	if err != nil {
		return nil, nil, err
	}

	projectId := d.Get(sfProjectId).(string)

	return &secureFileId, &projectId, nil
}

func expandAllowAccess(d *schema.ResourceData, secureFile *taskagent.SecureFile) []build.DefinitionResourceReference {
	resourceRefType := "securefile"
	secureFileId := secureFile.Id.String()
	authorized := d.Get(sfAllowAccess).(bool)

	return []build.DefinitionResourceReference{
		{
			Type:       &resourceRefType,
			Id:         &secureFileId,
			Name:       secureFile.Name,
			Authorized: &authorized,
		},
	}
}

func flattenAllowAccess(d *schema.ResourceData, definitionResources *[]build.DefinitionResourceReference) {
	secureFileId := d.Id()
	var allowAccess = false
	if definitionResources != nil {
		for _, resource := range *definitionResources {
			if secureFileId == *resource.Id {
				allowAccess = *resource.Authorized
			}
		}
	}
	_ = d.Set(sfAllowAccess, allowAccess)
}

func authorizeProjectReferences(
	clients *client.Clients, ctx context.Context, projectId *string,
	definitionResources []build.DefinitionResourceReference,
) (*[]build.DefinitionResourceReference, error) {
	return clients.BuildClient.AuthorizeProjectResources(
		ctx, build.AuthorizeProjectResourcesArgs{
			Project:   projectId,
			Resources: &definitionResources,
		},
	)
}

func secureFileContentHash(v interface{}) string {
	switch v := v.(type) {
	case string:
		data, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			data = []byte(v)
		}
		hash := sha256.Sum256(data)
		return hex.EncodeToString(hash[:])
	default:
		return ""
	}
}
