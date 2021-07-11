package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/google/uuid"
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
	id, _ := uuid.NewRandom()
	fileName := fmt.Sprintf("%s.txt", id)

	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckProject(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecureFileConfig(projectId, fileName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"azdoext_secure_file.foo", "name", regexp.MustCompile("^"+id.String())),
				),
			},
			{
				Config: testAccResourceSecureFileConfig(projectId, "foo-" + fileName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"azdoext_secure_file.foo", "name", regexp.MustCompile("^foo-"))),
			},
		},
	})
}

func testAccResourceSecureFileConfig(projectId string, fileName string) string {
	return fmt.Sprintf(`
resource "azdoext_secure_file" "foo" {
  project_id = "%s"
  name = "%s"
  content = "Hello World"
}
`, projectId, fileName)
}
