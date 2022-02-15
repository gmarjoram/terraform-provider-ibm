// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package continuousdeliverypipeline

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.ibm.com/org-ids/tekton-pipeline-go-sdk/continuousdeliverypipelinev2"
)

func ResourceIBMTektonPipelineProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceIBMTektonPipelinePropertyCreate,
		ReadContext:   ResourceIBMTektonPipelinePropertyRead,
		UpdateContext: ResourceIBMTektonPipelinePropertyUpdate,
		DeleteContext: ResourceIBMTektonPipelinePropertyDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_tekton_pipeline_property", "pipeline_id"),
				Description:  "The tekton pipeline ID.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_tekton_pipeline_property", "name"),
				Description:  "Property name.",
			},
			"value": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_tekton_pipeline_property", "value"),
				Description:  "String format property value.",
			},
			"options": &schema.Schema{
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Options for SINGLE_SELECT property type.",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_tekton_pipeline_property", "type"),
				Description:  "Property type.",
			},
			"path": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_tekton_pipeline_property", "path"),
				Description:  "property path for INTEGRATION type properties.",
			},
		},
	}
}

func ResourceIBMTektonPipelinePropertyValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 1)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "pipeline_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z]+$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[-0-9a-zA-Z_.]{1,234}$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "value",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `.`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "INTEGRATION, SECURE, SINGLE_SELECT, TEXT",
		},
		validate.ValidateSchema{
			Identifier:                 "path",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `.`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_tekton_pipeline_property", Schema: validateSchema}
	return &resourceValidator
}

func ResourceIBMTektonPipelinePropertyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	continuousDeliveryPipelineClient, err := meta.(conns.ClientSession).ContinuousDeliveryPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createTektonPipelinePropertiesOptions := &continuousdeliverypipelinev2.CreateTektonPipelinePropertiesOptions{}

	createTektonPipelinePropertiesOptions.SetPipelineID(d.Get("pipeline_id").(string))
	if _, ok := d.GetOk("name"); ok {
		createTektonPipelinePropertiesOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("value"); ok {
		createTektonPipelinePropertiesOptions.SetValue(d.Get("value").(string))
	}
	if _, ok := d.GetOk("options"); ok {
		createTektonPipelinePropertiesOptions.SetOptions(d.Get("options").(interface{}))
	}
	if _, ok := d.GetOk("type"); ok {
		createTektonPipelinePropertiesOptions.SetType(d.Get("type").(string))
	}
	if _, ok := d.GetOk("path"); ok {
		createTektonPipelinePropertiesOptions.SetPath(d.Get("path").(string))
	}

	property, response, err := continuousDeliveryPipelineClient.CreateTektonPipelinePropertiesWithContext(context, createTektonPipelinePropertiesOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateTektonPipelinePropertiesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateTektonPipelinePropertiesWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createTektonPipelinePropertiesOptions.PipelineID, *property.Name))

	return ResourceIBMTektonPipelinePropertyRead(context, d, meta)
}

func ResourceIBMTektonPipelinePropertyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	continuousDeliveryPipelineClient, err := meta.(conns.ClientSession).ContinuousDeliveryPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelinePropertyOptions := &continuousdeliverypipelinev2.GetTektonPipelinePropertyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelinePropertyOptions.SetPipelineID(parts[0])
	getTektonPipelinePropertyOptions.SetPropertyName(parts[1])

	property, response, err := continuousDeliveryPipelineClient.GetTektonPipelinePropertyWithContext(context, getTektonPipelinePropertyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetTektonPipelinePropertyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTektonPipelinePropertyWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("pipeline_id", getTektonPipelinePropertyOptions.PipelineID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting pipeline_id: %s", err))
	}
	if err = d.Set("name", property.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("value", property.Value); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting value: %s", err))
	}
	if err = d.Set("options", property.Options); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting options: %s", err))
	}
	if err = d.Set("type", property.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}
	if err = d.Set("path", property.Path); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting path: %s", err))
	}

	return nil
}

func ResourceIBMTektonPipelinePropertyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	continuousDeliveryPipelineClient, err := meta.(conns.ClientSession).ContinuousDeliveryPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	replaceTektonPipelinePropertyOptions := &continuousdeliverypipelinev2.ReplaceTektonPipelinePropertyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	replaceTektonPipelinePropertyOptions.SetPipelineID(parts[0])
	replaceTektonPipelinePropertyOptions.SetPropertyName(parts[1])
	replaceTektonPipelinePropertyOptions.SetName(d.Get("name").(string))
	replaceTektonPipelinePropertyOptions.SetType(d.Get("type").(string))

	hasChange := false

	if d.HasChange("pipeline_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "pipeline_id"))
	}

	if d.Get("type").(string) == "SINGLE_SELECT" && d.HasChange("options") {
		replaceTektonPipelinePropertyOptions.SetOptions(d.Get("options").(interface{}))
		hasChange = true
	} else if d.HasChange("value") {
		replaceTektonPipelinePropertyOptions.SetValue(d.Get("value").(string))
		hasChange = true
	}

	if d.HasChange("path") {
		replaceTektonPipelinePropertyOptions.SetPath(d.Get("path").(string))
		hasChange = true
	}

	if hasChange {
		_, response, err := continuousDeliveryPipelineClient.ReplaceTektonPipelinePropertyWithContext(context, replaceTektonPipelinePropertyOptions)
		if err != nil {
			log.Printf("[DEBUG] ReplaceTektonPipelinePropertyWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ReplaceTektonPipelinePropertyWithContext failed %s\n%s", err, response))
		}
	}

	return ResourceIBMTektonPipelinePropertyRead(context, d, meta)
}

func ResourceIBMTektonPipelinePropertyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	continuousDeliveryPipelineClient, err := meta.(conns.ClientSession).ContinuousDeliveryPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTektonPipelinePropertyOptions := &continuousdeliverypipelinev2.DeleteTektonPipelinePropertyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTektonPipelinePropertyOptions.SetPipelineID(parts[0])
	deleteTektonPipelinePropertyOptions.SetPropertyName(parts[1])

	response, err := continuousDeliveryPipelineClient.DeleteTektonPipelinePropertyWithContext(context, deleteTektonPipelinePropertyOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteTektonPipelinePropertyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteTektonPipelinePropertyWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
