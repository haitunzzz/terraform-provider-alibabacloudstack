package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackLogMachineGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackLogMachineGroupCreate,
		Read:   resourceAlibabacloudStackLogMachineGroupRead,
		Update: resourceAlibabacloudStackLogMachineGroupUpdate,
		Delete: resourceAlibabacloudStackLogMachineGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"identify_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      sls.MachineIDTypeIP,
				ValidateFunc: validation.StringInSlice([]string{sls.MachineIDTypeIP, sls.MachineIDTypeUserDefined}, false),
			},
			"topic": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"identify_list": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MinItems: 1,
			},
		},
	}
}

func resourceAlibabacloudStackLogMachineGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	params := &sls.MachineGroup{
		Name:          d.Get("name").(string),
		MachineIDType: d.Get("identify_type").(string),
		MachineIDList: expandStringList(d.Get("identify_list").(*schema.Set).List()),
		Attribute: sls.MachinGroupAttribute{
			TopicName: d.Get("topic").(string),
		},
	}
	var requestInfo *sls.Client
	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.CreateMachineGroup(d.Get("project").(string), params)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("CreateMachineGroup", raw, requestInfo, map[string]interface{}{
				"project":      d.Get("project").(string),
				"MachineGroup": params,
			})
		}
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_log_machine_group", "CreateMachineGroup", AlibabacloudStackLogGoSdkERROR)
	}

	d.SetId(fmt.Sprintf("%s%s%s", d.Get("project").(string), COLON_SEPARATED, d.Get("name").(string)))

	return resourceAlibabacloudStackLogMachineGroupRead(d, meta)
}

func resourceAlibabacloudStackLogMachineGroupRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	object, err := logService.DescribeLogMachineGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("project", parts[0])
	d.Set("name", object.Name)
	d.Set("identify_type", object.MachineIDType)
	d.Set("identify_list", object.MachineIDList)
	d.Set("topic", object.Attribute.TopicName)

	return nil
}

func resourceAlibabacloudStackLogMachineGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("identify_type") || d.HasChange("identify_list") || d.HasChange("topic") {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return WrapError(err)
		}

		client := meta.(*connectivity.AlibabacloudStackClient)
		var requestInfo *sls.Client
		params := &sls.MachineGroup{
			Name:          parts[1],
			MachineIDType: d.Get("identify_type").(string),
			MachineIDList: expandStringList(d.Get("identify_list").(*schema.Set).List()),
			Attribute: sls.MachinGroupAttribute{
				TopicName: d.Get("topic").(string),
			},
		}
		if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
			raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
				requestInfo = slsClient
				return nil, slsClient.UpdateMachineGroup(parts[0], params)
			})
			if err != nil {
				if IsExpectedErrors(err, []string{LogClientTimeout}) {
					time.Sleep(5 * time.Second)
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			if debugOn() {
				addDebug("UpdateMachineGroup", raw, requestInfo, map[string]interface{}{
					"project":      parts[0],
					"MachineGroup": params,
				})
			}
			return nil
		}); err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateMachineGroup", AlibabacloudStackLogGoSdkERROR)
		}
	}

	return resourceAlibabacloudStackLogMachineGroupRead(d, meta)
}

func resourceAlibabacloudStackLogMachineGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	var requestInfo *sls.Client
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithLogClient(func(slsClient *sls.Client) (interface{}, error) {
			requestInfo = slsClient
			return nil, slsClient.DeleteMachineGroup(parts[0], parts[1])
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			addDebug("DeleteMachineGroup", raw, requestInfo, map[string]interface{}{
				"project":      parts[0],
				"machineGroup": parts[1],
			})
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_log_store", "ListShards", AlibabacloudStackLogGoSdkERROR)
	}
	return WrapError(logService.WaitForLogMachineGroup(d.Id(), Deleted, DefaultTimeout))
}
