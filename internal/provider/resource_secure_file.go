package provider

import (
	"context"
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/Xtansia/terraform-provider-azdoext/internal/client"
	"github.com/Xtansia/terraform-provider-azdoext/internal/client/taskagent"
	"github.com/Xtansia/terraform-provider-azdoext/internal/utils"
)

const (
	sfProjectId     = "project_id"
	sfName          = "name"
	sfContent       = "content"
	sfContentBase64 = "content_base64"
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
				ForceNew:      true,
				Sensitive:     true,
				ConflictsWith: []string{sfContentBase64},
				ValidateFunc:  validation.StringIsNotEmpty,
			},
			sfContentBase64: {
				Description:   "The base64 encoded content of the secure file.",
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Sensitive:     true,
				ConflictsWith: []string{sfContent},
				ValidateFunc:  utils.StringIsBase64EncodedAndNotEmpty,
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

	projectId := d.Get(sfProjectId).(string)
	name := d.Get(sfName).(string)
	content := d.Get(sfContent).(string)
	contentBase64 := d.Get(sfContentBase64).(string)
	var data []byte

	if content != "" {
		data = []byte(content)
	} else if contentBase64 != "" {
		data, _ = base64.StdEncoding.DecodeString(contentBase64)
	} else {
		return diag.Errorf("one of %q or %q must be set", sfContent, sfContentBase64)
	}

	createdSecureFile, err := clients.TaskAgentClient.UploadSecureFile(ctx, taskagent.UploadSecureFileArgs{
		Project: &projectId,
		Name:    &name,
		Content: &data,
	})
	if err != nil {
		return diag.Errorf("Error creating secure file in Azure DevOps: %+v", err)
	}

	flattenSecureFile(d, createdSecureFile, &projectId)

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
		})
	if err != nil {
		if utils.ResponseWasNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error looking up secure file given ID (%v) and project ID (%v): %v", secureFileId, projectId, err)
	}
	if secureFile.Id == nil {
		d.SetId("")
		return nil
	}

	flattenSecureFile(d, secureFile, projectId)

	return nil
}

func resourceSecureFileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(*client.Clients)

	secureFileId, projectId, err := parseSecureFileAndProjectIds(d)
	if err != nil {
		return diag.Errorf(invalidSecureFileIdErrorMessageFormat, err)
	}

	name := d.Get(sfName).(string)

	updatedSecureFile, err := clients.TaskAgentClient.UpdateSecureFile(
		ctx,
		taskagent.UpdateSecureFileArgs{
			Project:      projectId,
			SecureFileId: secureFileId,
			SecureFile: &taskagent.SecureFile{
				Name: &name,
			},
		})

	if err != nil {
		return diag.Errorf("Error updating secure file in Azure DevOps: %+v", err)
	}

	flattenSecureFile(d, updatedSecureFile, projectId)

	return resourceSecureFileRead(ctx, d, meta)
}

func resourceSecureFileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(*client.Clients)

	secureFileId, projectId, err := parseSecureFileAndProjectIds(d)
	if err != nil {
		return diag.Errorf(invalidSecureFileIdErrorMessageFormat, err)
	}

	err = clients.TaskAgentClient.DeleteSecureFile(
		ctx,
		taskagent.DeleteSecureFileArgs{
			Project:      projectId,
			SecureFileId: secureFileId,
		})
	if err != nil {
		return diag.Errorf("Error deleting secure file in Azure DevOps: %+v", err)
	}

	return nil
}

func flattenSecureFile(d *schema.ResourceData, secureFile *taskagent.SecureFile, projectId *string) {
	d.SetId(secureFile.Id.String())
	d.Set(sfName, *secureFile.Name)
	d.Set(sfProjectId, projectId)
}

func parseSecureFileAndProjectIds(d *schema.ResourceData) (*uuid.UUID, *string, error) {
	secureFileId, err := uuid.Parse(d.Id())
	if err != nil {
		return nil, nil, err
	}

	projectId := d.Get(sfProjectId).(string)

	return &secureFileId, &projectId, nil
}
