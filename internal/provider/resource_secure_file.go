package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	sfProjectId     = "project_id"
	sfName          = "name"
	sfContent       = "content"
	sfContentBase64 = "content_base64"
)

func resourceSecureFile() *schema.Resource {
	return &schema.Resource{
		Description: "Manages secure files within Azure DevOps.",

		CreateContext: resourceSecureFileCreate,
		ReadContext:   resourceSecureFileRead,
		UpdateContext: resourceSecureFileUpdate,
		DeleteContext: resourceSecureFileDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSecureFileImport,
		},

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
			},
			sfContentBase64: {
				Description:   "The base64 encoded content of the secure file.",
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Sensitive:     true,
				ConflictsWith: []string{sfContent},
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
	return diag.Errorf("not yet implemented")
}

func resourceSecureFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("not yet implemented")
}

func resourceSecureFileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("not yet implemented")
}

func resourceSecureFileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("not yet implemented")
}

func resourceSecureFileImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, fmt.Errorf("not yet implemented")
}
