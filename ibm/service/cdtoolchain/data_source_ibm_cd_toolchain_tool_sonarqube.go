// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cdtoolchain

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/continuous-delivery-go-sdk/cdtoolchainv2"
)

func DataSourceIBMCdToolchainToolSonarqube() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCdToolchainToolSonarqubeRead,

		Schema: map[string]*schema.Schema{
			"toolchain_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the toolchain.",
			},
			"tool_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the tool bound to the toolchain.",
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group where tool can be found.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tool CRN.",
			},
			"toolchain_crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of toolchain which the tool is bound to.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI representing the tool.",
			},
			"referent": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information on URIs to access this resource through the UI or API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ui_href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URI representing the this resource through the UI.",
						},
						"api_href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URI representing the this resource through an API.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tool name.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Latest tool update timestamp.",
			},
			"parameters": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Unique key-value pairs representing parameters to be used to create the tool.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type a name for this tool integration, for example: my-sonarqube. This name displays on your toolchain.",
						},
						"dashboard_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type the URL of the SonarQube instance that you want to open when you click the SonarQube card in your toolchain.",
						},
						"user_login": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "If you are using an authentication token, leave this field empty.",
						},
						"user_password": &schema.Schema{
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"blind_connection": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Select this checkbox only if the server is not addressable on the public internet. IBM Cloud will not be able to validate the connection details you provide.",
						},
					},
				},
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current configuration state of the tool.",
			},
		},
	}
}

func dataSourceIBMCdToolchainToolSonarqubeRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getToolByIDOptions := &cdtoolchainv2.GetToolByIDOptions{}

	getToolByIDOptions.SetToolchainID(d.Get("toolchain_id").(string))
	getToolByIDOptions.SetToolID(d.Get("tool_id").(string))

	toolchainTool, response, err := cdToolchainClient.GetToolByIDWithContext(context, getToolByIDOptions)
	if err != nil {
		log.Printf("[DEBUG] GetToolByIDWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetToolByIDWithContext failed %s\n%s", err, response))
	}

	if *toolchainTool.ToolTypeID != "sonarqube" {
		return diag.FromErr(fmt.Errorf("Retrieved tool is not the correct type: %s", *toolchainTool.ToolTypeID))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getToolByIDOptions.ToolchainID, *getToolByIDOptions.ToolID))

	if err = d.Set("resource_group_id", toolchainTool.ResourceGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group_id: %s", err))
	}

	if err = d.Set("crn", toolchainTool.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}

	if err = d.Set("toolchain_crn", toolchainTool.ToolchainCRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting toolchain_crn: %s", err))
	}

	if err = d.Set("href", toolchainTool.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	referent := []map[string]interface{}{}
	if toolchainTool.Referent != nil {
		modelMap, err := dataSourceIBMCdToolchainToolSonarqubeToolModelReferentToMap(toolchainTool.Referent)
		if err != nil {
			return diag.FromErr(err)
		}
		referent = append(referent, modelMap)
	}
	if err = d.Set("referent", referent); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting referent %s", err))
	}

	if err = d.Set("name", toolchainTool.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("updated_at", flex.DateTimeToString(toolchainTool.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}

	parameters := []map[string]interface{}{}
	if toolchainTool.Parameters != nil {
		modelMap := GetParametersFromRead(toolchainTool.Parameters, DataSourceIBMCdToolchainToolSonarqube(), nil)
		parameters = append(parameters, modelMap)
	}
	if err = d.Set("parameters", parameters); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting parameters %s", err))
	}

	if err = d.Set("state", toolchainTool.State); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}

	return nil
}

func dataSourceIBMCdToolchainToolSonarqubeToolModelReferentToMap(model *cdtoolchainv2.ToolModelReferent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UIHref != nil {
		modelMap["ui_href"] = *model.UIHref
	}
	if model.APIHref != nil {
		modelMap["api_href"] = *model.APIHref
	}
	return modelMap, nil
}
