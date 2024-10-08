package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackArmsAlertContact() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackArmsAlertContactCreate,
		Read:   resourceAlibabacloudStackArmsAlertContactRead,
		Update: resourceAlibabacloudStackArmsAlertContactUpdate,
		Delete: resourceAlibabacloudStackArmsAlertContactDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"alert_contact_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ding_robot_webhook_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"phone_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_noc": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceAlibabacloudStackArmsAlertContactCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateAlertContact"
	request := make(map[string]interface{})
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("alert_contact_name"); ok {
		request["ContactName"] = v
	}
	if v, ok := d.GetOk("ding_robot_webhook_url"); ok {
		request["DingRobotWebhookUrl"] = v
	} else if v, ok := d.GetOk("email"); ok && v.(string) == "" {
		if v, ok := d.GetOk("phone_num"); ok && v.(string) == "" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "ding_robot_webhook_url", "email", d.Get("email"), "phone_num", d.Get("phone_num")))
		}
	}
	if v, ok := d.GetOk("email"); ok {
		request["Email"] = v
	} else if v, ok := d.GetOk("ding_robot_webhook_url"); ok && v.(string) == "" {
		if v, ok := d.GetOk("phone_num"); ok && v.(string) == "" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "email", "ding_robot_webhook_url", d.Get("ding_robot_webhook_url"), "phone_num", d.Get("phone_num")))
		}
	}
	if v, ok := d.GetOk("phone_num"); ok {
		request["PhoneNum"] = v
	} else if v, ok := d.GetOk("ding_robot_webhook_url"); ok && v.(string) == "" {
		if v, ok := d.GetOk("email"); ok && v.(string) == "" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "phone_num", "ding_robot_webhook_url", d.Get("ding_robot_webhook_url"), "email", d.Get("email")))
		}
	}
	request["RegionId"] = client.RegionId
	request["Product"] = "ARMS"
	request["OrganizationId"] = client.Department
	//request["product"] = "ARMS"

	if v, ok := d.GetOkExists("system_noc"); ok {
		request["SystemNoc"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_arms_alert_contact", action, AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["ContactId"]))

	return resourceAlibabacloudStackArmsAlertContactRead(d, meta)
}
func resourceAlibabacloudStackArmsAlertContactRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	armsService := ArmsService{client}
	object, err := armsService.DescribeArmsAlertContact(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_arms_alert_contact armsService.DescribeArmsAlertContact Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("alert_contact_name", object["ContactName"])
	d.Set("ding_robot_webhook_url", object["DingRobot"])
	d.Set("email", object["Email"])
	d.Set("phone_num", object["Phone"])
	d.Set("system_noc", object["SystemNoc"])
	return nil
}
func resourceAlibabacloudStackArmsAlertContactUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"ContactId": d.Id(),
	}
	request["RegionId"] = client.RegionId

	request["Product"] = "ARMS"
	request["OrganizationId"] = client.Department
	if d.HasChange("alert_contact_name") {
		update = true
	}
	if v, ok := d.GetOk("alert_contact_name"); ok {
		request["ContactName"] = v
	}
	if d.HasChange("ding_robot_webhook_url") {
		update = true
	}
	if v, ok := d.GetOk("ding_robot_webhook_url"); ok {
		request["DingRobotWebhookUrl"] = v
	} else if v, ok := d.GetOk("email"); ok && v.(string) == "" {
		if v, ok := d.GetOk("phone_num"); ok && v.(string) == "" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "ding_robot_webhook_url", "email", d.Get("email"), "phone_num", d.Get("phone_num")))
		}
	}
	if d.HasChange("email") {
		update = true
	}
	if v, ok := d.GetOk("email"); ok {
		request["Email"] = v
	} else if v, ok := d.GetOk("ding_robot_webhook_url"); ok && v.(string) == "" {
		if v, ok := d.GetOk("phone_num"); ok && v.(string) == "" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "email", "ding_robot_webhook_url", d.Get("ding_robot_webhook_url"), "phone_num", d.Get("phone_num")))
		}
	}
	if d.HasChange("phone_num") {
		update = true
	}
	if v, ok := d.GetOk("phone_num"); ok {
		request["PhoneNum"] = v
	} else if v, ok := d.GetOk("ding_robot_webhook_url"); ok && v.(string) == "" {
		if v, ok := d.GetOk("email"); ok && v.(string) == "" {
			return WrapError(fmt.Errorf("attribute '%s' is required when '%s' is %v and '%s' is %v ", "phone_num", "ding_robot_webhook_url", d.Get("ding_robot_webhook_url"), "email", d.Get("email")))
		}
	}
	if d.HasChange("system_noc") || d.IsNewResource() {
		update = true
	}
	if v, ok := d.GetOkExists("system_noc"); ok {
		request["SystemNoc"] = v
	}
	if update {
		action := "UpdateAlertContact"
		conn, err := client.NewArmsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
		}
	}
	return resourceAlibabacloudStackArmsAlertContactRead(d, meta)
}
func resourceAlibabacloudStackArmsAlertContactDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteAlertContact"
	var response map[string]interface{}
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ContactId": d.Id(),
	}

	request["RegionId"] = client.RegionId

	request["Product"] = "ARMS"
	request["OrganizationId"] = client.Department
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	return nil
}
