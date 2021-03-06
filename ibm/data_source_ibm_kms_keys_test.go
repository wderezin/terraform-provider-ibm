package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIBMKmsDataSource_basic(t *testing.T) {
	instanceName := fmt.Sprintf("kms_%d", acctest.RandIntRange(10, 100))
	// bucketName := fmt.Sprintf("bucket", acctest.RandIntRange(10, 100))
	keyName := fmt.Sprintf("key_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMKmsDataSourceConfig(instanceName, keyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_kms_key.test", "key_name", keyName),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMKmsDataSourceHpcsConfig(hpcsInstanceID, keyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_kms_key.test", "key_name", keyName),
				),
			},
		},
	})
}

func testAccCheckIBMKmsDataSourceConfig(instanceName, keyName string) string {
	return fmt.Sprintf(`
	resource "ibm_resource_instance" "kms_instance" {
		name              = "%s"
		service           = "kms"
		plan              = "tiered-pricing"
		location          = "us-south"
	  }
	  resource "ibm_kms_key" "test" {
		instance_id = "${ibm_resource_instance.kms_instance.guid}"
		key_name = "%s"
		standard_key =  true
		force_delete = true
	}
	data "ibm_kms_keys" "test" {
		instance_id = "${ibm_kms_key.test.instance_id}" 
	}
`, instanceName, keyName)
}

func testAccCheckIBMKmsDataSourceHpcsConfig(hpcsInstanceID, KeyName string) string {
	return fmt.Sprintf(`
	  resource "ibm_kms_key" "test" {
		instance_id = "%s"
		key_name = "%s"
		standard_key =  true
		force_delete = true
	}
	data "ibm_kms_keys" "test" {
		instance_id = "${ibm_kms_key.test.instance_id}" 
	}
	
`, hpcsInstanceID, KeyName)
}
