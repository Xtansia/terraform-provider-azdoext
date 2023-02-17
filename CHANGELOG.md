# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html), and is generated
by [Changie](https://github.com/miniscruff/changie).## v0.4.0 - 2023-02-17

### Changed

* Updated `github.com/hashicorp/terraform-plugin-sdk/v2` from `v2.10.1` to `v2.13.0`
* Updated `github.com/hashicorp/terraform-plugin-log` from `v0.2.1` to `v0.3.0`
* Updated `github.com/hashicorp/terraform-plugin-docs` from `v0.5.1` to `v0.7.0`
* Updated `github.com/stretchr/testify` from `v1.7.0` to `v1.7.1`
* Updated Go to 1.18
* Updated `github.com/hashicorp/terraform-plugin-docs` from `v0.7.0` to `v0.13.0`
* Updated `github.com/hashicorp/terraform-plugin-log` from `v0.3.0` to `v0.8.0`
* Updated `github.com/hashicorp/terraform-plugin-sdk/v2` from `v2.13.0` to `v2.25.0`
* Updated `github.com/stretchr/testify` from `v1.7.1` to `v1.8.1`## v0.3.2 - 2022-02-16

### Changed

* Updated `github.com/microsoft/azure-devops-go-api/azuredevops` from `v1.0.0-b5` to `v6.0.1`
* Updated `github.com/hashicorp/terraform-plugin-sdk/v2` from `v2.7.0` to `v2.10.1`
* Updated `github.com/hashicorp/terraform-plugin-docs` from `v0.4.0` to `v0.5.1`
* Use tflog instead of log## v0.3.1 - 2022-02-16

### Changed

* Updated Go to 1.17 to allow building darwin_arm64## v0.3.0 - 2021-07-15

### Added

* Implemented `properties` on `azdoext_secure_file`

### Changed

* Allow `content` and `content_base64` to be empty on `azdoext_secure_file`
* `azdoext_secure_file` now stores hash of content in state rather than full plaintext content## v0.2.0 - 2021-07-12

### Added

* Implemented `allow_access` on the `azdoext_secure_file` resource## v0.1.1 - 2021-07-12

### Changed

* Reduced required Azure DevOps API version to `5.0`## v0.1.0 - 2021-07-10

### Added

* Initial implementation of `azdoext_secure_file` resource