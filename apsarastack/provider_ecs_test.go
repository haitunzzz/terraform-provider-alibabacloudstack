package apsarastack

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackProviderEcs(t *testing.T) {
	var v ecs.Instance

	resourceId := "apsarastack_instance.default"
	ra := resourceAttrInit(resourceId, testAccInstanceCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAcc%sEcsInstanceConfigVpc%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return providerCommon + resourceInstanceVpcConfigDependence(name)
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":        "${data.apsarastack_images.default.images.0.id}",
					"security_groups": []string{"${apsarastack_security_group.default.0.id}"},
					"instance_type":   "${local.instance_type_id}",

					"availability_zone":             "${data.apsarastack_instance_types.default.instance_types.0.availability_zones.0}",
					"system_disk_category":          "cloud_efficiency",
					"instance_name":                 "${var.name}",
					"key_name":                      "${apsarastack_key_pair.default.key_name}",
					"security_enhancement_strategy": "Active",
					"user_data":                     "I_am_user_data",

					"vswitch_id": "${apsarastack_vswitch.default.id}",
					"role_name":  "${apsarastack_ram_role.default.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"key_name":      name,
						"role_name":     name,
					}),
				),
			},
		},
	})
}
