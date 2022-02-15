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
	"github.ibm.com/org-ids/tekton-pipeline-go-sdk/continuousdeliverypipelinev2"
)

func DataSourceIBMTektonPipelineTriggerProperty() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMTektonPipelineTriggerPropertyRead,

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The tekton pipeline ID.",
			},
			"trigger_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The trigger ID.",
			},
			"property_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The property's name.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Property name.",
			},
			"value": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "String format property value.",
			},
			"options": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Options for SINGLE_SELECT property type.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Property type.",
			},
			"path": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "property path for INTEGRATION type properties.",
			},
		},
	}
}

func DataSourceIBMTektonPipelineTriggerPropertyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	continuousDeliveryPipelineClient, err := meta.(conns.ClientSession).ContinuousDeliveryPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelineTriggerPropertyOptions := &continuousdeliverypipelinev2.GetTektonPipelineTriggerPropertyOptions{}

	getTektonPipelineTriggerPropertyOptions.SetPipelineID(d.Get("pipeline_id").(string))
	getTektonPipelineTriggerPropertyOptions.SetTriggerID(d.Get("trigger_id").(string))
	getTektonPipelineTriggerPropertyOptions.SetPropertyName(d.Get("property_name").(string))

	triggerProperty, response, err := continuousDeliveryPipelineClient.GetTektonPipelineTriggerPropertyWithContext(context, getTektonPipelineTriggerPropertyOptions)
	if err != nil {
		log.Printf("[DEBUG] GetTektonPipelineTriggerPropertyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTektonPipelineTriggerPropertyWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *getTektonPipelineTriggerPropertyOptions.PipelineID, *getTektonPipelineTriggerPropertyOptions.TriggerID, *getTektonPipelineTriggerPropertyOptions.PropertyName))

	if err = d.Set("name", triggerProperty.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("value", triggerProperty.Value); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting value: %s", err))
	}

	if err = d.Set("options", triggerProperty.Options); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting options: %s", err))
	}

	if err = d.Set("type", triggerProperty.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}

	if err = d.Set("path", triggerProperty.Path); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting path: %s", err))
	}

	return nil
}
