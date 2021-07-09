data "azuredevops_project" "proj" {
  name = "My Project"
}

resource "azdoext_secure_file" "file" {
  project_id = data.azuredevops_project.proj.id
  name       = "hello_world.txt"
  content    = <<-EOT
    Hello World, I'm a very secure file.
  EOT
}