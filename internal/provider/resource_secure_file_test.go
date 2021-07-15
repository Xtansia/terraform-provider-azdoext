package provider

import (
	"fmt"
	"os"
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
	content := "Hello World"
	contentHash := "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e" // SHA256 hash of "Hello World"

	resource.UnitTest(
		t, resource.TestCase{
			PreCheck: func() {
				preCheck(t)
				preCheckProject(t)
			},
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceSecureFileConfig(projectId, fileName, content, false, false),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("azdoext_secure_file.foo", "name", id.String()+".txt"),
						resource.TestCheckResourceAttr("azdoext_secure_file.foo", "content", contentHash),
						resource.TestCheckNoResourceAttr("azdoext_secure_file.foo", "properties.foo"),
					),
				},
				{
					Config: testAccResourceSecureFileConfig(projectId, fileName, content, true, false),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("azdoext_secure_file.foo", "name", id.String()+".txt"),
						resource.TestCheckResourceAttr("azdoext_secure_file.foo", "content_base64", contentHash),
						resource.TestCheckNoResourceAttr("azdoext_secure_file.foo", "properties.foo"),
					),
				},
				{
					Config: testAccResourceSecureFileConfig(projectId, "foo-"+fileName, content, true, true),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("azdoext_secure_file.foo", "name", "foo-"+id.String()+".txt"),
						resource.TestCheckResourceAttr("azdoext_secure_file.foo", "content_base64", contentHash),
						resource.TestCheckNoResourceAttr("azdoext_secure_file.foo", "properties.foo"),
					),
				},
				{
					Config: testAccResourceSecureFileConfigWithProps(projectId, "foo-"+fileName, content, true, true),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("azdoext_secure_file.foo", "name", "foo-"+id.String()+".txt"),
						resource.TestCheckResourceAttr("azdoext_secure_file.foo", "content_base64", contentHash),
						resource.TestCheckResourceAttr("azdoext_secure_file.foo", "properties.foo", "bar"),
					),
				},
				{
					Config: testAccResourceSecureFileConfig(projectId, "foo-"+fileName, content, true, true),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("azdoext_secure_file.foo", "name", "foo-"+id.String()+".txt"),
						resource.TestCheckResourceAttr("azdoext_secure_file.foo", "content_base64", contentHash),
						resource.TestCheckNoResourceAttr("azdoext_secure_file.foo", "properties.foo"),
					),
				},
			},
		},
	)
}

func testAccResourceSecureFileConfig(
	projectId string, fileName string, content string, base64Encoded bool, allowAccess bool,
) string {
	if base64Encoded {
		content = fmt.Sprintf("content_base64 = base64encode(%q)", content)
	} else {
		content = fmt.Sprintf("content = %q", content)
	}

	return fmt.Sprintf(
		`
resource "azdoext_secure_file" "foo" {
  project_id = %q
  name = %q
  %s
  allow_access = %v
}
`, projectId, fileName, content, allowAccess,
	)
}

func testAccResourceSecureFileConfigWithProps(
	projectId string, fileName string, content string, base64Encoded bool, allowAccess bool,
) string {
	if base64Encoded {
		content = fmt.Sprintf("content_base64 = base64encode(%q)", content)
	} else {
		content = fmt.Sprintf("content = %q", content)
	}

	return fmt.Sprintf(
		`
resource "azdoext_secure_file" "foo" {
  project_id = %q
  name = %q
  %s
  allow_access = %v
  properties = {
    foo = "bar"
  }
}
`, projectId, fileName, content, allowAccess,
	)
}
