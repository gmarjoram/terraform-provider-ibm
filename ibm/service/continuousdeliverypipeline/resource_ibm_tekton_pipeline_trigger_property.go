// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package continuousdeliverypipeline

import (
	"context"
	"crypto/hmac"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/google/go-cmp/cmp"
	"github.ibm.com/org-ids/tekton-pipeline-go-sdk/continuousdeliverypipelinev2"
)

func ResourceIBMTektonPipelineTriggerProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceIBMTektonPipelineTriggerPropertyCreate,
		ReadContext:   ResourceIBMTektonPipelineTriggerPropertyRead,
		UpdateContext: ResourceIBMTektonPipelineTriggerPropertyUpdate,
		DeleteContext: ResourceIBMTektonPipelineTriggerPropertyDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_tekton_pipeline_trigger_property", "pipeline_id"),
				Description:  "The tekton pipeline ID.",
			},
			"trigger_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_tekton_pipeline_trigger_property", "trigger_id"),
				Description:  "The trigger ID.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_tekton_pipeline_trigger_property", "name"),
				Description:  "Property name.",
			},
			"value": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_tekton_pipeline_trigger_property", "value"),
				Description:  "String format property value.",
			},
			"enum": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Options for SINGLE_SELECT property type.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"default": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default option for SINGLE_SELECT property type.",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_tekton_pipeline_trigger_property", "type"),
				Description:  "Property type.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("type").(string) == "SECURE" {
						parts, _ := flex.SepIdParts(d.Id(), "/")
						segs := []string{parts[0], parts[1], d.Get("name").(string)}
						secret := strings.Join(segs, ".")
						mac := hmac.New(sha3.New512, []byte(secret))
						mac.Write([]byte(new))
						secureHmac := hex.EncodeToString(mac.Sum(nil))
						hasEnvChange := !cmp.Equal(strings.Join([]string{"hash", "SHA3-512", secureHmac}, ":"), old)
						if hasEnvChange {
							return false
						}
						return true
					} else {
						if old == new {
							return true
						}
						return false
					}
				},
			},
			"path": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_tekton_pipeline_trigger_property", "path"),
				Description:  "property path for INTEGRATION type properties.",
			},
		},
	}
}

func ResourceIBMTektonPipelineTriggerPropertyValidator() *validate.ResourceValidator {
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
			Identifier:                 "trigger_id",
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
			AllowedValues:              "INTEGRATION, SECURE, SINGLE_SELECT, TEXT, APPCONFIG",
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_tekton_pipeline_trigger_property", Schema: validateSchema}
	return &resourceValidator
}

func ResourceIBMTektonPipelineTriggerPropertyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	continuousDeliveryPipelineClient, err := meta.(conns.ClientSession).ContinuousDeliveryPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createTektonPipelineTriggerPropertiesOptions := &continuousdeliverypipelinev2.CreateTektonPipelineTriggerPropertiesOptions{}

	createTektonPipelineTriggerPropertiesOptions.SetPipelineID(d.Get("pipeline_id").(string))
	createTektonPipelineTriggerPropertiesOptions.SetTriggerID(d.Get("trigger_id").(string))
	if _, ok := d.GetOk("name"); ok {
		createTektonPipelineTriggerPropertiesOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("value"); ok {
		createTektonPipelineTriggerPropertiesOptions.SetValue(d.Get("value").(string))
	}
	if _, ok := d.GetOk("enum"); ok {
		createTektonPipelineTriggerPropertiesOptions.SetEnum(d.Get("enum").([]string))
	}
	if _, ok := d.GetOk("default"); ok {
		createTektonPipelineTriggerPropertiesOptions.SetDefault(d.Get("default").(string))
	}
	if _, ok := d.GetOk("type"); ok {
		createTektonPipelineTriggerPropertiesOptions.SetType(d.Get("type").(string))
	}
	if _, ok := d.GetOk("path"); ok {
		createTektonPipelineTriggerPropertiesOptions.SetPath(d.Get("path").(string))
	}

	triggerProperty, response, err := continuousDeliveryPipelineClient.CreateTektonPipelineTriggerPropertiesWithContext(context, createTektonPipelineTriggerPropertiesOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateTektonPipelineTriggerPropertiesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateTektonPipelineTriggerPropertiesWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *createTektonPipelineTriggerPropertiesOptions.PipelineID, *createTektonPipelineTriggerPropertiesOptions.TriggerID, *triggerProperty.Name))

	return ResourceIBMTektonPipelineTriggerPropertyRead(context, d, meta)
}

func ResourceIBMTektonPipelineTriggerPropertyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	continuousDeliveryPipelineClient, err := meta.(conns.ClientSession).ContinuousDeliveryPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelineTriggerPropertyOptions := &continuousdeliverypipelinev2.GetTektonPipelineTriggerPropertyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelineTriggerPropertyOptions.SetPipelineID(parts[0])
	getTektonPipelineTriggerPropertyOptions.SetTriggerID(parts[1])
	getTektonPipelineTriggerPropertyOptions.SetPropertyName(parts[2])

	triggerProperty, response, err := continuousDeliveryPipelineClient.GetTektonPipelineTriggerPropertyWithContext(context, getTektonPipelineTriggerPropertyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetTektonPipelineTriggerPropertyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTektonPipelineTriggerPropertyWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("pipeline_id", getTektonPipelineTriggerPropertyOptions.PipelineID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting pipeline_id: %s", err))
	}
	if err = d.Set("trigger_id", getTektonPipelineTriggerPropertyOptions.TriggerID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting trigger_id: %s", err))
	}
	if err = d.Set("name", triggerProperty.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("value", triggerProperty.Value); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting value: %s", err))
	}

	if triggerProperty.Enum != nil {
		if err = d.Set("enum", triggerProperty.Enum); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting enum: %s", err))
		}
	}
	if err = d.Set("default", triggerProperty.Default); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting default: %s", err))
	}
	if err = d.Set("type", triggerProperty.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}
	if err = d.Set("path", triggerProperty.Path); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting path: %s", err))
	}

	return nil
}

func ResourceIBMTektonPipelineTriggerPropertyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	continuousDeliveryPipelineClient, err := meta.(conns.ClientSession).ContinuousDeliveryPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	replaceTektonPipelineTriggerPropertyOptions := &continuousdeliverypipelinev2.ReplaceTektonPipelineTriggerPropertyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	replaceTektonPipelineTriggerPropertyOptions.SetPipelineID(parts[0])
	replaceTektonPipelineTriggerPropertyOptions.SetTriggerID(parts[1])
	replaceTektonPipelineTriggerPropertyOptions.SetPropertyName(parts[2])
	replaceTektonPipelineTriggerPropertyOptions.SetName(d.Get("name").(string))
	replaceTektonPipelineTriggerPropertyOptions.SetType(d.Get("type").(string))

	hasChange := false

	if d.HasChange("pipeline_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "pipeline_id"))
	}
	if d.HasChange("trigger_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "trigger_id"))
	}
	if d.HasChange("enum") {
		replaceTektonPipelineTriggerPropertyOptions.SetEnum(d.Get("enum").([]string))
		hasChange = true
	}
	if d.HasChange("default") {
		replaceTektonPipelineTriggerPropertyOptions.SetDefault(d.Get("default").(string))
		hasChange = true
	}
	if d.HasChange("path") {
		replaceTektonPipelineTriggerPropertyOptions.SetPath(d.Get("path").(string))
		hasChange = true
	}

	if hasChange {
		_, response, err := continuousDeliveryPipelineClient.ReplaceTektonPipelineTriggerPropertyWithContext(context, replaceTektonPipelineTriggerPropertyOptions)
		if err != nil {
			log.Printf("[DEBUG] ReplaceTektonPipelineTriggerPropertyWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ReplaceTektonPipelineTriggerPropertyWithContext failed %s\n%s", err, response))
		}
	}

	return ResourceIBMTektonPipelineTriggerPropertyRead(context, d, meta)
}

func ResourceIBMTektonPipelineTriggerPropertyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	continuousDeliveryPipelineClient, err := meta.(conns.ClientSession).ContinuousDeliveryPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTektonPipelineTriggerPropertyOptions := &continuousdeliverypipelinev2.DeleteTektonPipelineTriggerPropertyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTektonPipelineTriggerPropertyOptions.SetPipelineID(parts[0])
	deleteTektonPipelineTriggerPropertyOptions.SetTriggerID(parts[1])
	deleteTektonPipelineTriggerPropertyOptions.SetPropertyName(parts[2])

	response, err := continuousDeliveryPipelineClient.DeleteTektonPipelineTriggerPropertyWithContext(context, deleteTektonPipelineTriggerPropertyOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteTektonPipelineTriggerPropertyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteTektonPipelineTriggerPropertyWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
