/**
 * (C) Copyright IBM Corp. 2021.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func resourceIbmIamAccountSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIamAccountSettingsCreate,
		ReadContext:   resourceIbmIamAccountSettingsRead,
		UpdateContext: resourceIbmIamAccountSettingsUpdate,
		DeleteContext: resourceIbmIamAccountSettingsDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"include_history": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Defines if the entity history is included in the response.",
			},
			"context": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Context with key properties for problem determination.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transaction_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The transaction ID of the inbound REST request.",
						},
						"operation": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The operation of the inbound REST request.",
						},
						"user_agent": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The user agent of the inbound REST request.",
						},
						"url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL of that cluster.",
						},
						"instance_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The instance ID of the server instance processing the request.",
						},
						"thread_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The thread ID of the server instance processing the request.",
						},
						"host": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The host of the server instance processing the request.",
						},
						"start_time": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The start time of the request.",
						},
						"end_time": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The finish time of the request.",
						},
						"elapsed_time": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The elapsed time in msec.",
						},
						"cluster_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The cluster name.",
						},
					},
				},
			},
			"restrict_create_service_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines whether or not creating a Service Id is access controlled. Valid values:  * RESTRICTED - to apply access control  * NOT_RESTRICTED - to remove access control  * NOT_SET - to 'unset' a previous set value.",
			},
			"restrict_create_platform_apikey": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines whether or not creating platform API keys is access controlled. Valid values:  * RESTRICTED - to apply access control  * NOT_RESTRICTED - to remove access control  * NOT_SET - to 'unset' a previous set value.",
			},
			"allowed_ip_addresses": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the IP addresses and subnets from which IAM tokens can be created for the account.",
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Version of the account settings.",
			},
			"mfa": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the MFA trait for the account. Valid values:  * NONE - No MFA trait set  * TOTP - For all non-federated IBMId users  * TOTP4ALL - For all users  * LEVEL1 - Email-based MFA for all users  * LEVEL2 - TOTP-based MFA for all users  * LEVEL3 - U2F MFA for all users.",
			},
			"if_match": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Version of the account settings to be updated. Specify the version that you retrieved as entity_tag (ETag header) when reading the account. This value helps identifying parallel usage of this API. Pass * to indicate to update any version available. This might result in stale updates.",
				Default:     "*",
			},
			"history": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "History of the Account Settings.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Timestamp when the action was triggered.",
						},
						"iam_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "IAM ID of the identity which triggered the action.",
						},
						"iam_id_account": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Account of the identity which triggered the action.",
						},
						"action": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Action of the history entry.",
						},
						"params": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "Params of the history entry.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Message which summarizes the executed action.",
						},
					},
				},
			},
			"session_expiration_in_seconds": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the session expiration in seconds for the account. Valid values:  * Any whole number between between '900' and '86400'  * NOT_SET - To unset account setting and use service default.",
			},
			"session_invalidation_in_seconds": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Defines the period of time in seconds in which a session will be invalidated due  to inactivity. Valid values:   * Any whole number between '900' and '7200'   * NOT_SET - To unset account setting and use service default.",
			},
		},
	}
}

func resourceIbmIamAccountSettingsCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(ClientSession).IamIdentityV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getAccountSettingsOptions := &iamidentityv1.GetAccountSettingsOptions{}

	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}
	getAccountSettingsOptions.SetAccountID(userDetails.userAccount)
	if _, ok := d.GetOk("include_history"); ok {
		getAccountSettingsOptions.SetIncludeHistory(d.Get("include_history").(bool))
	}

	accountSettingsResponse, response, err := iamIdentityClient.GetAccountSettingsWithContext(context, getAccountSettingsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetAccountSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", *getAccountSettingsOptions.AccountID, *accountSettingsResponse.AccountID))

	return resourceIbmIamAccountSettingsUpdate(context, d, meta)
}

func resourceIbmIamAccountSettingsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(ClientSession).IamIdentityV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getAccountSettingsOptions := &iamidentityv1.GetAccountSettingsOptions{}

	parts, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getAccountSettingsOptions.SetAccountID(parts[0])
	getAccountSettingsOptions.SetAccountID(parts[1])
	getAccountSettingsOptions.SetIncludeHistory(d.Get("include_history").(bool))

	accountSettingsResponse, response, err := iamIdentityClient.GetAccountSettingsWithContext(context, getAccountSettingsOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetAccountSettingsWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	if accountSettingsResponse.Context != nil {
		contextMap := resourceIbmIamAccountSettingsResponseContextToMap(*accountSettingsResponse.Context)
		if err = d.Set("context", []map[string]interface{}{contextMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting context: %s", err))
		}
	}
	if err = d.Set("restrict_create_service_id", accountSettingsResponse.RestrictCreateServiceID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting restrict_create_service_id: %s", err))
	}
	if err = d.Set("restrict_create_platform_apikey", accountSettingsResponse.RestrictCreatePlatformApikey); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting restrict_create_platform_apikey: %s", err))
	}
	if err = d.Set("allowed_ip_addresses", accountSettingsResponse.AllowedIPAddresses); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting allowed_ip_addresses: %s", err))
	}
	if err = d.Set("entity_tag", accountSettingsResponse.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}
	if err = d.Set("mfa", accountSettingsResponse.Mfa); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting mfa: %s", err))
	}
	if accountSettingsResponse.History != nil {
		history := []map[string]interface{}{}
		for _, historyItem := range accountSettingsResponse.History {
			historyItemMap := resourceIbmIamAccountSettingsEnityHistoryRecordToMap(historyItem)
			history = append(history, historyItemMap)
		}
		if err = d.Set("history", history); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting history: %s", err))
		}
	}
	if err = d.Set("session_expiration_in_seconds", accountSettingsResponse.SessionExpirationInSeconds); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting session_expiration_in_seconds: %s", err))
	}
	if err = d.Set("session_invalidation_in_seconds", accountSettingsResponse.SessionInvalidationInSeconds); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting session_invalidation_in_seconds: %s", err))
	}

	return nil
}

func resourceIbmIamAccountSettingsResponseContextToMap(responseContext iamidentityv1.ResponseContext) map[string]interface{} {
	responseContextMap := map[string]interface{}{}

	responseContextMap["transaction_id"] = responseContext.TransactionID
	responseContextMap["operation"] = responseContext.Operation
	responseContextMap["user_agent"] = responseContext.UserAgent
	responseContextMap["url"] = responseContext.URL
	responseContextMap["instance_id"] = responseContext.InstanceID
	responseContextMap["thread_id"] = responseContext.ThreadID
	responseContextMap["host"] = responseContext.Host
	responseContextMap["start_time"] = responseContext.StartTime
	responseContextMap["end_time"] = responseContext.EndTime
	responseContextMap["elapsed_time"] = responseContext.ElapsedTime
	responseContextMap["cluster_name"] = responseContext.ClusterName

	return responseContextMap
}

func resourceIbmIamAccountSettingsEnityHistoryRecordToMap(enityHistoryRecord iamidentityv1.EnityHistoryRecord) map[string]interface{} {
	enityHistoryRecordMap := map[string]interface{}{}

	enityHistoryRecordMap["timestamp"] = enityHistoryRecord.Timestamp
	enityHistoryRecordMap["iam_id"] = enityHistoryRecord.IamID
	enityHistoryRecordMap["iam_id_account"] = enityHistoryRecord.IamIDAccount
	enityHistoryRecordMap["action"] = enityHistoryRecord.Action
	enityHistoryRecordMap["params"] = enityHistoryRecord.Params
	enityHistoryRecordMap["message"] = enityHistoryRecord.Message

	return enityHistoryRecordMap
}

func resourceIbmIamAccountSettingsUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(ClientSession).IamIdentityV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateAccountSettingsOptions := &iamidentityv1.UpdateAccountSettingsOptions{}

	parts, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	updateAccountSettingsOptions.SetAccountID(parts[0])
	updateAccountSettingsOptions.SetAccountID(parts[1])
	updateAccountSettingsOptions.SetIfMatch(d.Get("if_match").(string))

	hasChange := false

	if d.HasChange("restrict_create_service_id") {
		restrict_create_service_id_str := d.Get("restrict_create_service_id").(string)
		updateAccountSettingsOptions.SetRestrictCreateServiceID(restrict_create_service_id_str)
		hasChange = true
	}

	if d.HasChange("restrict_create_platform_apikey") {
		restrict_create_platform_apikey_str := d.Get("restrict_create_platform_apikey").(string)
		updateAccountSettingsOptions.SetRestrictCreatePlatformApikey(restrict_create_platform_apikey_str)
		hasChange = true
	}

	if d.HasChange("mfa") {
		mfa_str := d.Get("mfa").(string)
		updateAccountSettingsOptions.SetMfa(mfa_str)
		hasChange = true
	}

	if d.HasChange("session_expiration_in_seconds") {
		session_expiration_in_seconds_str := d.Get("session_expiration_in_seconds").(string)
		updateAccountSettingsOptions.SetSessionExpirationInSeconds(session_expiration_in_seconds_str)
		hasChange = true
	}

	if d.HasChange("session_invalidation_in_seconds") {
		session_invalidation_in_seconds_str := d.Get("session_invalidation_in_seconds").(string)
		updateAccountSettingsOptions.SetSessionInvalidationInSeconds(session_invalidation_in_seconds_str)
		hasChange = true
	}

	if hasChange {
		_, response, err := iamIdentityClient.UpdateAccountSettingsWithContext(context, updateAccountSettingsOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateAccountSettingsWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
	}

	return resourceIbmIamAccountSettingsRead(context, d, meta)
}

func resourceIbmIamAccountSettingsDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	// DELETE NOT SUPPORTED
	d.SetId("")

	return nil
}
