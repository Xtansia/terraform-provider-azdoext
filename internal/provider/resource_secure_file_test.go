package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func preCheckProject(t *testing.T) {
	projectId := os.Getenv("AZDO_TEST_PROJECT_ID")
	if projectId == "" {
		t.Skipf("AZDO_TEST_PROJECT_ID not set")
	}
}

func TestAccResourceSecureFile(t *testing.T) {
	projectId := os.Getenv("AZDO_TEST_PROJECT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckProject(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecureFileConfig(projectId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"azdoext_secure_file.foo", "name", regexp.MustCompile("^foo")),
				),
			},
		},
	})
}

func testAccResourceSecureFileConfig(projectId string) string {
	return fmt.Sprintf(`
resource "azdoext_secure_file" "foo" {
  project_id = "%s"
  name = "foobar"
  content = "Hello World"
}
`, projectId)
}
