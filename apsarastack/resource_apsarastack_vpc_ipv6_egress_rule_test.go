package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackVPCIpv6EgressRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "apsarastack_vpc_ipv6_egress_rule.default"
	ra := resourceAttrInit(resourceId, ApsaraStackVPCIpv6EgressRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}, "DescribeVpcIpv6EgressRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcipv6egressrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, ApsaraStackVPCIpv6EgressRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ECS_WITH_IPV6_ADDRESS")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ipv6_egress_rule_name": "${var.name}",
					"ipv6_gateway_id":       "${data.apsarastack_vpc_ipv6_addresses.default.addresses.0.ipv6_gateway_id}",
					"instance_id":           "${data.apsarastack_vpc_ipv6_addresses.default.ids.0}",
					"instance_type":         "Ipv6Address",
					"description":           "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipv6_egress_rule_name": name,
						"ipv6_gateway_id":       CHECKSET,
						"instance_id":           CHECKSET,
						"instance_type":         "Ipv6Address",
						"description":           name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var ApsaraStackVPCIpv6EgressRuleMap0 = map[string]string{
	"instance_type": CHECKSET,
	"status":        CHECKSET,
}

func ApsaraStackVPCIpv6EgressRuleBasicDependence0(name string) string {
	return fmt.Sprintf(` 
provider "apsarastack" {
	assume_role {}
}
variable "name" {
  default = "%s"
}

data "apsarastack_instances" "default" {
  name_regex = "no-deleteing-ipv6-address"
  status     = "Running"
}

data "apsarastack_vpc_ipv6_addresses" "default" {
  associated_instance_id = data.apsarastack_instances.default.instances.0.id
  status                 = "Available"
}

`, name)
}
