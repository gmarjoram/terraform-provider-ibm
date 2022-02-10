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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.ibm.com/org-ids/tekton-pipeline-go-sdk/continuousdeliverypipelinev2"
)

func ResourceIBMTektonPipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceIBMTektonPipelineCreate,
		ReadContext:   ResourceIBMTektonPipelineRead,
		UpdateContext: ResourceIBMTektonPipelineUpdate,
		DeleteContext: ResourceIBMTektonPipelineDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"integration_instance_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_tekton_pipeline", "integration_instance_id"),
				Description:  "UUID.",
			},
			"worker": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Worker object with just worker ID.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "String.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Pipeline status.",
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID.",
			},
			"toolchain": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Toolchain object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "UUID.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CRN for the toolchain that containing the tekton pipeline.",
						},
					},
				},
			},
			"definitions": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Definition list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scm_source": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "Scm source for tekton pipeline defintion.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "General href URL.",
									},
									"branch": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The branch of the repo.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The path to the definitions yaml files.",
									},
								},
							},
						},
						"service_instance_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "UUID.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "UUID.",
						},
					},
				},
			},
			"env_properties": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tekton pipeline level environment properties.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Property name.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "String format property value.",
						},
						"options": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Options for SINGLE_SELECT property type.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Property type.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "property path for INTEGRATION type properties.",
						},
					},
				},
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Standard RFC 3339 Date Time String.",
			},
			"created": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Standard RFC 3339 Date Time String.",
			},
			"pipeline_definition": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tekton pipeline definition document detail object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The state of pipeline definition status.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "UUID.",
						},
					},
				},
			},
			"triggers": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tekton pipeline triggers list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Trigger type.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Trigger name.",
						},
						"event_listener": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Event listener name.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "UUID.",
						},
						"properties": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Trigger properties.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Property name.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "String format property value.",
									},
									"options": &schema.Schema{
										Type:        schema.TypeMap,
										Optional:    true,
										Description: "Options for SINGLE_SELECT property type.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Property type.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "property path for INTEGRATION type properties.",
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "General href URL.",
									},
								},
							},
						},
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Trigger tags array.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"worker": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Trigger worker used to run the trigger, the trigger worker overrides the default pipeline worker.If not exist, this trigger uses default pipeline worker.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "worker name.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "worker type.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "ID.",
									},
								},
							},
						},
						"concurrency": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Concurrency object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_concurrent_runs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Defines the maximum number of concurrent runs for this trigger.",
									},
								},
							},
						},
						"disabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Defines if this trigger is disabled.",
						},
						"scm_source": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Scm source for git type tekton pipeline trigger.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Needed only for git trigger type. Repo URL that listening to.",
									},
									"branch": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Needed only for git trigger type. Branch name of the repo.",
									},
									"blind_connection": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Needed only for git trigger type. Branch name of the repo.",
									},
									"hook_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Webhook Id.",
									},
								},
							},
						},
						"events": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Needed only for git trigger type. Events object defines the events this git trigger listening to.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"push": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If true, the trigger will start when a 'push' event received.",
									},
									"pull_request_closed": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If true, the trigger will start when a pull request 'close' event received.",
									},
									"pull_request": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If true, the trigger will start when a pull request 'open' or 'update' event received.",
									},
								},
							},
						},
						"service_instance_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "UUID.",
						},
						"cron": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Needed only for timer trigger type. Cron expression for timer trigger.",
						},
						"timezone": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Needed only for timer trigger type. Timezones for timer trigger.",
						},
						"secret": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Needed only for generic trigger type. Secret used to execute generic trigger.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Secret type.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Secret value, not needed if secret type is \"internalValidation\".",
									},
									"source": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Secret location, not needed if secret type is \"internalValidation\".",
									},
									"key_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Secret name, not needed if type is \"internalValidation\".",
									},
									"algorithm": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Algorithm used for \"digestMatches\" secret type.",
									},
								},
							},
						},
					},
				},
			},
			"next_timers": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tekton pipeline timer for next timer type trigger.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "timer type.",
						},
						"created": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Standard RFC 3339 Date Time String.",
						},
						"updated_at": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Standard RFC 3339 Date Time String.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "UUID.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Timer name.",
						},
						"trigger_name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the trigger that created this timer.",
						},
						"trigger_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the trigger that created this timer.",
						},
						"timezone": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Time zones.",
						},
						"sub": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "User Email address that created this timer.",
						},
						"event_listener": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Event listener of the trigger that created this timer.",
						},
						"cron": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cron expression for timer.",
						},
						"disabled": &schema.Schema{
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Defines if this timer is disabled.",
						},
						"expiration": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The next tick for this timer.",
						},
					},
				},
			},
			"html_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Dashboard URL of this pipeline.",
			},
			"build_number": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Latest pipeline run build number. If this property is absent, the pipeline has not had any pipelineRuns.",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag whether this pipeline enabled.",
			},
		},
	}
}

func ResourceIBMTektonPipelineValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 1)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "integration_instance_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[-0-9a-z]+$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_tekton_pipeline", Schema: validateSchema}
	return &resourceValidator
}

func ResourceIBMTektonPipelineCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	continuousDeliveryPipelineClient, err := meta.(conns.ClientSession).ContinuousDeliveryPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createTektonPipelineOptions := &continuousdeliverypipelinev2.CreateTektonPipelineOptions{}

	if _, ok := d.GetOk("integration_instance_id"); ok {
		createTektonPipelineOptions.SetIntegrationInstanceID(d.Get("integration_instance_id").(string))
	}
	if _, ok := d.GetOk("worker"); ok {
		worker, err := ResourceIBMTektonPipelineMapToWorkerWithID(d.Get("worker.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createTektonPipelineOptions.SetWorker(worker)
	}

	tektonPipeline, response, err := continuousDeliveryPipelineClient.CreateTektonPipelineWithContext(context, createTektonPipelineOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateTektonPipelineWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateTektonPipelineWithContext failed %s\n%s", err, response))
	}

	d.SetId(*tektonPipeline.ID)

	return ResourceIBMTektonPipelineRead(context, d, meta)
}

func ResourceIBMTektonPipelineRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	continuousDeliveryPipelineClient, err := meta.(conns.ClientSession).ContinuousDeliveryPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelineOptions := &continuousdeliverypipelinev2.GetTektonPipelineOptions{}

	getTektonPipelineOptions.SetID(d.Id())

	tektonPipeline, response, err := continuousDeliveryPipelineClient.GetTektonPipelineWithContext(context, getTektonPipelineOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetTektonPipelineWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTektonPipelineWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("integration_instance_id", getTektonPipelineOptions.IntegrationInstanceID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting integration_instance_id: %s", err))
	}
	if tektonPipeline.Worker != nil {
		workerMap, err := ResourceIBMTektonPipelineWorkerWithIDToMap(tektonPipeline.Worker)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("worker", []map[string]interface{}{workerMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting worker: %s", err))
		}
	}
	if err = d.Set("name", tektonPipeline.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("status", tektonPipeline.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}
	if err = d.Set("resource_group_id", tektonPipeline.ResourceGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group_id: %s", err))
	}
	toolchainMap, err := ResourceIBMTektonPipelineToolchainToMap(tektonPipeline.Toolchain)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("toolchain", []map[string]interface{}{toolchainMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting toolchain: %s", err))
	}

	definitions := []map[string]interface{}{}
	for _, definitionsItem := range tektonPipeline.Definitions {
		definitionsItemMap, err := ResourceIBMTektonPipelineDefinitionToMap(&definitionsItem)
		if err != nil {
			return diag.FromErr(err)
		}
		definitions = append(definitions, definitionsItemMap)
	}
	if err = d.Set("definitions", definitions); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting definitions: %s", err))
	}

	envProperties := []map[string]interface{}{}
	for _, envPropertiesItem := range tektonPipeline.EnvProperties {
		envPropertiesItemMap, err := ResourceIBMTektonPipelinePropertyToMap(&envPropertiesItem)
		if err != nil {
			return diag.FromErr(err)
		}
		envProperties = append(envProperties, envPropertiesItemMap)
	}
	if err = d.Set("env_properties", envProperties); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting env_properties: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(tektonPipeline.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("created", flex.DateTimeToString(tektonPipeline.Created)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created: %s", err))
	}
	if tektonPipeline.PipelineDefinition != nil {
		pipelineDefinitionMap, err := ResourceIBMTektonPipelineTektonPipelinePipelineDefinitionToMap(tektonPipeline.PipelineDefinition)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("pipeline_definition", []map[string]interface{}{pipelineDefinitionMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting pipeline_definition: %s", err))
		}
	}

	triggers := []map[string]interface{}{}
	for _, triggersItem := range tektonPipeline.Triggers {
		triggersItemMap, err := ResourceIBMTektonPipelineTriggerToMap(triggersItem)
		if err != nil {
			return diag.FromErr(err)
		}
		triggers = append(triggers, triggersItemMap)
	}
	if err = d.Set("triggers", triggers); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting triggers: %s", err))
	}

	nextTimers := []map[string]interface{}{}
	if tektonPipeline.NextTimers != nil {
		for _, nextTimersItem := range tektonPipeline.NextTimers {
			nextTimersItemMap, err := ResourceIBMTektonPipelineTimerToMap(&nextTimersItem)
			if err != nil {
				return diag.FromErr(err)
			}
			nextTimers = append(nextTimers, nextTimersItemMap)
		}
	}
	if err = d.Set("next_timers", nextTimers); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting next_timers: %s", err))
	}
	if err = d.Set("html_url", tektonPipeline.HTMLURL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting html_url: %s", err))
	}
	if err = d.Set("build_number", flex.IntValue(tektonPipeline.BuildNumber)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting build_number: %s", err))
	}
	if err = d.Set("enabled", tektonPipeline.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting enabled: %s", err))
	}

	return nil
}

func ResourceIBMTektonPipelineUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	continuousDeliveryPipelineClient, err := meta.(conns.ClientSession).ContinuousDeliveryPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	updateTektonPipelineOptions := &continuousdeliverypipelinev2.UpdateTektonPipelineOptions{}

	updateTektonPipelineOptions.SetID(d.Id())

	hasChange := false

	if d.HasChange("integration_instance_id") {
		updateTektonPipelineOptions.SetIntegrationInstanceID(d.Get("integration_instance_id").(string))
		hasChange = true
	}
	if d.HasChange("worker") {
		worker, err := ResourceIBMTektonPipelineMapToWorkerWithID(d.Get("worker.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateTektonPipelineOptions.SetWorker(worker)
		hasChange = true
	}

	if hasChange {
		_, response, err := continuousDeliveryPipelineClient.UpdateTektonPipelineWithContext(context, updateTektonPipelineOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateTektonPipelineWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateTektonPipelineWithContext failed %s\n%s", err, response))
		}
	}

	return ResourceIBMTektonPipelineRead(context, d, meta)
}

func ResourceIBMTektonPipelineDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	continuousDeliveryPipelineClient, err := meta.(conns.ClientSession).ContinuousDeliveryPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTektonPipelineOptions := &continuousdeliverypipelinev2.DeleteTektonPipelineOptions{}

	deleteTektonPipelineOptions.SetID(d.Id())

	response, err := continuousDeliveryPipelineClient.DeleteTektonPipelineWithContext(context, deleteTektonPipelineOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteTektonPipelineWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteTektonPipelineWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func ResourceIBMTektonPipelineMapToWorkerWithID(modelMap map[string]interface{}) (*continuousdeliverypipelinev2.WorkerWithID, error) {
	model := &continuousdeliverypipelinev2.WorkerWithID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMTektonPipelineWorkerWithIDToMap(model *continuousdeliverypipelinev2.WorkerWithID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	return modelMap, nil
}

func ResourceIBMTektonPipelineToolchainToMap(model *continuousdeliverypipelinev2.Toolchain) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["crn"] = model.CRN
	return modelMap, nil
}

func ResourceIBMTektonPipelineDefinitionToMap(model *continuousdeliverypipelinev2.Definition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scmSourceMap, err := ResourceIBMTektonPipelineDefinitionScmSourceToMap(model.ScmSource)
	if err != nil {
		return modelMap, err
	}
	modelMap["scm_source"] = []map[string]interface{}{scmSourceMap}
	modelMap["service_instance_id"] = model.ServiceInstanceID
	modelMap["id"] = model.ID
	return modelMap, nil
}

func ResourceIBMTektonPipelineDefinitionScmSourceToMap(model *continuousdeliverypipelinev2.DefinitionScmSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["url"] = model.URL
	modelMap["branch"] = model.Branch
	modelMap["path"] = model.Path
	return modelMap, nil
}

func ResourceIBMTektonPipelinePropertyToMap(model *continuousdeliverypipelinev2.Property) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	if model.Options != nil {
		modelMap["options"] = model.Options
	}
	modelMap["type"] = model.Type
	if model.Path != nil {
		modelMap["path"] = model.Path
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTektonPipelinePipelineDefinitionToMap(model *continuousdeliverypipelinev2.TektonPipelinePipelineDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Status != nil {
		modelMap["status"] = model.Status
	}
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerToMap(model continuousdeliverypipelinev2.TriggerIntf) (map[string]interface{}, error) {
	if _, ok := model.(*continuousdeliverypipelinev2.TriggerManualTrigger); ok {
		return ResourceIBMTektonPipelineTriggerManualTriggerToMap(model.(*continuousdeliverypipelinev2.TriggerManualTrigger))
	} else if _, ok := model.(*continuousdeliverypipelinev2.TriggerScmTrigger); ok {
		return ResourceIBMTektonPipelineTriggerScmTriggerToMap(model.(*continuousdeliverypipelinev2.TriggerScmTrigger))
	} else if _, ok := model.(*continuousdeliverypipelinev2.TriggerTimerTrigger); ok {
		return ResourceIBMTektonPipelineTriggerTimerTriggerToMap(model.(*continuousdeliverypipelinev2.TriggerTimerTrigger))
	} else if _, ok := model.(*continuousdeliverypipelinev2.TriggerGenericTrigger); ok {
		return ResourceIBMTektonPipelineTriggerGenericTriggerToMap(model.(*continuousdeliverypipelinev2.TriggerGenericTrigger))
	} else if _, ok := model.(*continuousdeliverypipelinev2.Trigger); ok {
		modelMap := make(map[string]interface{})
		model := model.(*continuousdeliverypipelinev2.Trigger)
		if model.Type != nil {
			modelMap["type"] = model.Type
		}
		if model.Name != nil {
			modelMap["name"] = model.Name
		}
		if model.EventListener != nil {
			modelMap["event_listener"] = model.EventListener
		}
		if model.ID != nil {
			modelMap["id"] = model.ID
		}
		if model.Properties != nil {
			properties := []map[string]interface{}{}
			for _, propertiesItem := range model.Properties {
				propertiesItemMap, err := ResourceIBMTektonPipelineTriggerPropertiesItemToMap(&propertiesItem)
				if err != nil {
					return modelMap, err
				}
				properties = append(properties, propertiesItemMap)
			}
			modelMap["properties"] = properties
		}
		if model.Tags != nil {
			modelMap["tags"] = model.Tags
		}
		if model.Worker != nil {
			workerMap, err := ResourceIBMTektonPipelineWorkerToMap(model.Worker)
			if err != nil {
				return modelMap, err
			}
			modelMap["worker"] = []map[string]interface{}{workerMap}
		}
		if model.Concurrency != nil {
			concurrencyMap, err := ResourceIBMTektonPipelineTriggerConcurrencyToMap(model.Concurrency)
			if err != nil {
				return modelMap, err
			}
			modelMap["concurrency"] = []map[string]interface{}{concurrencyMap}
		}
		if model.Disabled != nil {
			modelMap["disabled"] = model.Disabled
		}
		if model.ScmSource != nil {
			scmSourceMap, err := ResourceIBMTektonPipelineTriggerScmSourceToMap(model.ScmSource)
			if err != nil {
				return modelMap, err
			}
			modelMap["scm_source"] = []map[string]interface{}{scmSourceMap}
		}
		if model.Events != nil {
			eventsMap, err := ResourceIBMTektonPipelineEventsToMap(model.Events)
			if err != nil {
				return modelMap, err
			}
			modelMap["events"] = []map[string]interface{}{eventsMap}
		}
		if model.ServiceInstanceID != nil {
			modelMap["service_instance_id"] = model.ServiceInstanceID
		}
		if model.Cron != nil {
			modelMap["cron"] = model.Cron
		}
		if model.Timezone != nil {
			modelMap["timezone"] = model.Timezone
		}
		if model.Secret != nil {
			secretMap, err := ResourceIBMTektonPipelineGenericSecretToMap(model.Secret)
			if err != nil {
				return modelMap, err
			}
			modelMap["secret"] = []map[string]interface{}{secretMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized continuousdeliverypipelinev2.TriggerIntf subtype encountered")
	}
}

func ResourceIBMTektonPipelineTriggerPropertiesItemToMap(model *continuousdeliverypipelinev2.TriggerPropertiesItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	if model.Options != nil {
		modelMap["options"] = model.Options
	}
	modelMap["type"] = model.Type
	if model.Path != nil {
		modelMap["path"] = model.Path
	}
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineWorkerToMap(model *continuousdeliverypipelinev2.Worker) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	modelMap["id"] = model.ID
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerConcurrencyToMap(model *continuousdeliverypipelinev2.TriggerConcurrency) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerScmSourceToMap(model *continuousdeliverypipelinev2.TriggerScmSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = model.URL
	}
	if model.Branch != nil {
		modelMap["branch"] = model.Branch
	}
	if model.BlindConnection != nil {
		modelMap["blind_connection"] = model.BlindConnection
	}
	if model.HookID != nil {
		modelMap["hook_id"] = model.HookID
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineEventsToMap(model *continuousdeliverypipelinev2.Events) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Push != nil {
		modelMap["push"] = model.Push
	}
	if model.PullRequestClosed != nil {
		modelMap["pull_request_closed"] = model.PullRequestClosed
	}
	if model.PullRequest != nil {
		modelMap["pull_request"] = model.PullRequest
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineGenericSecretToMap(model *continuousdeliverypipelinev2.GenericSecret) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	if model.Source != nil {
		modelMap["source"] = model.Source
	}
	if model.KeyName != nil {
		modelMap["key_name"] = model.KeyName
	}
	if model.Algorithm != nil {
		modelMap["algorithm"] = model.Algorithm
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerManualTriggerToMap(model *continuousdeliverypipelinev2.TriggerManualTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["name"] = model.Name
	modelMap["event_listener"] = model.EventListener
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := ResourceIBMTektonPipelineTriggerManualTriggerPropertiesItemToMap(&propertiesItem)
			if err != nil {
				return modelMap, err
			}
			properties = append(properties, propertiesItemMap)
		}
		modelMap["properties"] = properties
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.Worker != nil {
		workerMap, err := ResourceIBMTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.Concurrency != nil {
		concurrencyMap, err := ResourceIBMTektonPipelineTriggerManualTriggerConcurrencyToMap(model.Concurrency)
		if err != nil {
			return modelMap, err
		}
		modelMap["concurrency"] = []map[string]interface{}{concurrencyMap}
	}
	modelMap["disabled"] = model.Disabled
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerManualTriggerPropertiesItemToMap(model *continuousdeliverypipelinev2.TriggerManualTriggerPropertiesItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	if model.Options != nil {
		modelMap["options"] = model.Options
	}
	modelMap["type"] = model.Type
	if model.Path != nil {
		modelMap["path"] = model.Path
	}
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerManualTriggerConcurrencyToMap(model *continuousdeliverypipelinev2.TriggerManualTriggerConcurrency) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerScmTriggerToMap(model *continuousdeliverypipelinev2.TriggerScmTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["name"] = model.Name
	modelMap["event_listener"] = model.EventListener
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := ResourceIBMTektonPipelineTriggerScmTriggerPropertiesItemToMap(&propertiesItem)
			if err != nil {
				return modelMap, err
			}
			properties = append(properties, propertiesItemMap)
		}
		modelMap["properties"] = properties
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.Worker != nil {
		workerMap, err := ResourceIBMTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.Concurrency != nil {
		concurrencyMap, err := ResourceIBMTektonPipelineTriggerScmTriggerConcurrencyToMap(model.Concurrency)
		if err != nil {
			return modelMap, err
		}
		modelMap["concurrency"] = []map[string]interface{}{concurrencyMap}
	}
	modelMap["disabled"] = model.Disabled
	if model.ScmSource != nil {
		scmSourceMap, err := ResourceIBMTektonPipelineTriggerScmSourceToMap(model.ScmSource)
		if err != nil {
			return modelMap, err
		}
		modelMap["scm_source"] = []map[string]interface{}{scmSourceMap}
	}
	if model.Events != nil {
		eventsMap, err := ResourceIBMTektonPipelineEventsToMap(model.Events)
		if err != nil {
			return modelMap, err
		}
		modelMap["events"] = []map[string]interface{}{eventsMap}
	}
	if model.ServiceInstanceID != nil {
		modelMap["service_instance_id"] = model.ServiceInstanceID
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerScmTriggerPropertiesItemToMap(model *continuousdeliverypipelinev2.TriggerScmTriggerPropertiesItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	if model.Options != nil {
		modelMap["options"] = model.Options
	}
	modelMap["type"] = model.Type
	if model.Path != nil {
		modelMap["path"] = model.Path
	}
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerScmTriggerConcurrencyToMap(model *continuousdeliverypipelinev2.TriggerScmTriggerConcurrency) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerTimerTriggerToMap(model *continuousdeliverypipelinev2.TriggerTimerTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["name"] = model.Name
	modelMap["event_listener"] = model.EventListener
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := ResourceIBMTektonPipelineTriggerTimerTriggerPropertiesItemToMap(&propertiesItem)
			if err != nil {
				return modelMap, err
			}
			properties = append(properties, propertiesItemMap)
		}
		modelMap["properties"] = properties
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.Worker != nil {
		workerMap, err := ResourceIBMTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.Concurrency != nil {
		concurrencyMap, err := ResourceIBMTektonPipelineTriggerTimerTriggerConcurrencyToMap(model.Concurrency)
		if err != nil {
			return modelMap, err
		}
		modelMap["concurrency"] = []map[string]interface{}{concurrencyMap}
	}
	modelMap["disabled"] = model.Disabled
	if model.Cron != nil {
		modelMap["cron"] = model.Cron
	}
	if model.Timezone != nil {
		modelMap["timezone"] = model.Timezone
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerTimerTriggerPropertiesItemToMap(model *continuousdeliverypipelinev2.TriggerTimerTriggerPropertiesItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	if model.Options != nil {
		modelMap["options"] = model.Options
	}
	modelMap["type"] = model.Type
	if model.Path != nil {
		modelMap["path"] = model.Path
	}
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerTimerTriggerConcurrencyToMap(model *continuousdeliverypipelinev2.TriggerTimerTriggerConcurrency) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerGenericTriggerToMap(model *continuousdeliverypipelinev2.TriggerGenericTrigger) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["name"] = model.Name
	modelMap["event_listener"] = model.EventListener
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Properties != nil {
		properties := []map[string]interface{}{}
		for _, propertiesItem := range model.Properties {
			propertiesItemMap, err := ResourceIBMTektonPipelineTriggerGenericTriggerPropertiesItemToMap(&propertiesItem)
			if err != nil {
				return modelMap, err
			}
			properties = append(properties, propertiesItemMap)
		}
		modelMap["properties"] = properties
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.Worker != nil {
		workerMap, err := ResourceIBMTektonPipelineWorkerToMap(model.Worker)
		if err != nil {
			return modelMap, err
		}
		modelMap["worker"] = []map[string]interface{}{workerMap}
	}
	if model.Concurrency != nil {
		concurrencyMap, err := ResourceIBMTektonPipelineTriggerGenericTriggerConcurrencyToMap(model.Concurrency)
		if err != nil {
			return modelMap, err
		}
		modelMap["concurrency"] = []map[string]interface{}{concurrencyMap}
	}
	modelMap["disabled"] = model.Disabled
	if model.Secret != nil {
		secretMap, err := ResourceIBMTektonPipelineGenericSecretToMap(model.Secret)
		if err != nil {
			return modelMap, err
		}
		modelMap["secret"] = []map[string]interface{}{secretMap}
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerGenericTriggerPropertiesItemToMap(model *continuousdeliverypipelinev2.TriggerGenericTriggerPropertiesItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	if model.Options != nil {
		modelMap["options"] = model.Options
	}
	modelMap["type"] = model.Type
	if model.Path != nil {
		modelMap["path"] = model.Path
	}
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTriggerGenericTriggerConcurrencyToMap(model *continuousdeliverypipelinev2.TriggerGenericTriggerConcurrency) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxConcurrentRuns != nil {
		modelMap["max_concurrent_runs"] = flex.IntValue(model.MaxConcurrentRuns)
	}
	return modelMap, nil
}

func ResourceIBMTektonPipelineTimerToMap(model *continuousdeliverypipelinev2.Timer) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	modelMap["created"] = model.Created.String()
	modelMap["updated_at"] = model.UpdatedAt.String()
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["trigger_name"] = model.TriggerName
	modelMap["trigger_id"] = model.TriggerID
	modelMap["timezone"] = model.Timezone
	modelMap["sub"] = model.Sub
	modelMap["event_listener"] = model.EventListener
	modelMap["cron"] = model.Cron
	modelMap["disabled"] = model.Disabled
	modelMap["expiration"] = model.Expiration.String()
	return modelMap, nil
}
