// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package toolchain

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.ibm.com/org-ids/toolchain-go-sdk/toolchainv2"
)

func ResourceIbmToolchainToolSecretsmanager() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceIbmToolchainToolSecretsmanagerCreate,
		ReadContext:   ResourceIbmToolchainToolSecretsmanagerRead,
		UpdateContext: ResourceIbmToolchainToolSecretsmanagerUpdate,
		DeleteContext: ResourceIbmToolchainToolSecretsmanagerDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"toolchain_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_toolchain_tool_secretsmanager", "toolchain_id"),
				Description:  "ID of the toolchain to bind integration to.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of tool integration.",
			},
			"parameters": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"resource_group": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"instance_name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"parameters_references": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Decoded values used on provision in the broker that reference fields in the parameters.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"resource_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"crn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"toolchain_crn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"href": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"referent": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ui_href": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"api_href": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"updated_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"get_integration_by_id_response_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIbmToolchainToolSecretsmanagerValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 1)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "toolchain_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_toolchain_tool_secretsmanager", Schema: validateSchema}
	return &resourceValidator
}

func ResourceIbmToolchainToolSecretsmanagerCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	toolchainClient, err := meta.(conns.ClientSession).ToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	postIntegrationOptions := &toolchainv2.PostIntegrationOptions{}

	postIntegrationOptions.SetToolchainID(d.Get("toolchain_id").(string))
	postIntegrationOptions.SetToolID("secretsmanager")
	if _, ok := d.GetOk("name"); ok {
		postIntegrationOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("parameters"); ok {
		parametersModel, err := ResourceIbmToolchainToolSecretsmanagerMapToParameters(d.Get("parameters.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		postIntegrationOptions.SetParameters(parametersModel)
	}
	if _, ok := d.GetOk("parameters_references"); ok {
		// TODO: Add code to handle map container: ParametersReferences
	}

	postIntegrationResponse, response, err := toolchainClient.PostIntegrationWithContext(context, postIntegrationOptions)
	if err != nil {
		log.Printf("[DEBUG] PostIntegrationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("PostIntegrationWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *postIntegrationOptions.ToolchainID, *postIntegrationResponse.ID))

	return ResourceIbmToolchainToolSecretsmanagerRead(context, d, meta)
}

func ResourceIbmToolchainToolSecretsmanagerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	toolchainClient, err := meta.(conns.ClientSession).ToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getIntegrationByIdOptions := &toolchainv2.GetIntegrationByIdOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getIntegrationByIdOptions.SetToolchainID(parts[0])
	getIntegrationByIdOptions.SetIntegrationID(parts[1])

	getIntegrationByIdResponse, response, err := toolchainClient.GetIntegrationByIDWithContext(context, getIntegrationByIdOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetIntegrationByIDWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetIntegrationByIDWithContext failed %s\n%s", err, response))
	}

	// TODO: handle argument of type map[string]interface{}
	if err = d.Set("toolchain_id", getIntegrationByIdResponse.ToolchainID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting toolchain_id: %s", err))
	}
	if err = d.Set("name", getIntegrationByIdResponse.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if getIntegrationByIdResponse.Parameters != nil {
		parametersMap, err := ResourceIbmToolchainToolSecretsmanagerParametersToMap(getIntegrationByIdResponse.Parameters)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("parameters", []map[string]interface{}{parametersMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting parameters: %s", err))
		}
	}
	if err = d.Set("resource_group_id", getIntegrationByIdResponse.ResourceGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group_id: %s", err))
	}
	if err = d.Set("crn", getIntegrationByIdResponse.Crn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("toolchain_crn", getIntegrationByIdResponse.ToolchainCrn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting toolchain_crn: %s", err))
	}
	if err = d.Set("href", getIntegrationByIdResponse.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	referentMap, err := ResourceIbmToolchainToolSecretsmanagerGetIntegrationByIdResponseReferentToMap(getIntegrationByIdResponse.Referent)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("referent", []map[string]interface{}{referentMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting referent: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(getIntegrationByIdResponse.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("state", getIntegrationByIdResponse.State); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}
	if err = d.Set("get_integration_by_id_response_id", getIntegrationByIdResponse.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting get_integration_by_id_response_id: %s", err))
	}

	return nil
}

func ResourceIbmToolchainToolSecretsmanagerUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	toolchainClient, err := meta.(conns.ClientSession).ToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	patchToolIntegrationOptions := &toolchainv2.PatchToolIntegrationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	patchToolIntegrationOptions.SetToolchainID(parts[0])
	patchToolIntegrationOptions.SetIntegrationID(parts[1])
	patchToolIntegrationOptions.SetToolID("secretsmanager")

	hasChange := false

	if d.HasChange("toolchain_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "toolchain_id"))
	}
	if d.HasChange("name") {
		patchToolIntegrationOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("parameters") {
		parameters, err := ResourceIbmToolchainToolSecretsmanagerMapToParameters(d.Get("parameters.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		patchToolIntegrationOptions.SetParameters(parameters)
		hasChange = true
	}
	if d.HasChange("parameters_references") {
		// TODO: handle ParametersReferences of type TypeMap -- not primitive, not model
		hasChange = true
	}

	if hasChange {
		_, response, err := toolchainClient.PatchToolIntegrationWithContext(context, patchToolIntegrationOptions)
		if err != nil {
			log.Printf("[DEBUG] PatchToolIntegrationWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("PatchToolIntegrationWithContext failed %s\n%s", err, response))
		}
	}

	return ResourceIbmToolchainToolSecretsmanagerRead(context, d, meta)
}

func ResourceIbmToolchainToolSecretsmanagerDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	toolchainClient, err := meta.(conns.ClientSession).ToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteToolIntegrationOptions := &toolchainv2.DeleteToolIntegrationOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteToolIntegrationOptions.SetToolchainID(parts[0])
	deleteToolIntegrationOptions.SetIntegrationID(parts[1])

	response, err := toolchainClient.DeleteToolIntegrationWithContext(context, deleteToolIntegrationOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteToolIntegrationWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteToolIntegrationWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func ResourceIbmToolchainToolSecretsmanagerMapToParameters(modelMap map[string]interface{}) (map[string]interface{}, error) {
	model := make(map[string]interface{})
	if modelMap["name"] != nil {
		model["name"] = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["region"] != nil {
		model["region"] = core.StringPtr(modelMap["region"].(string))
	}
	if modelMap["resource_group"] != nil {
		model["resource-group"] = core.StringPtr(modelMap["resource_group"].(string))
	}
	if modelMap["instance_name"] != nil {
		model["instance-name"] = core.StringPtr(modelMap["instance_name"].(string))
	}
	return model, nil
}

func ResourceIbmToolchainToolSecretsmanagerParametersToMap(model map[string]interface{}) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model["name"] != nil {
		modelMap["name"] = model["name"]
	}
	if model["region"] != nil {
		modelMap["region"] = model["region"]
	}
	if model["resource-group"] != nil {
		modelMap["resource_group"] = model["resource-group"]
	}
	if model["instance-name"] != nil {
		modelMap["instance_name"] = model["instance-name"]
	}
	return modelMap, nil
}

func ResourceIbmToolchainToolSecretsmanagerGetIntegrationByIdResponseReferentToMap(model *toolchainv2.GetIntegrationByIdResponseReferent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UiHref != nil {
		modelMap["ui_href"] = model.UiHref
	}
	if model.ApiHref != nil {
		modelMap["api_href"] = model.ApiHref
	}
	return modelMap, nil
}
