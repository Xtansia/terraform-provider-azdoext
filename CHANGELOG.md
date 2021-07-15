# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html), and is generated
by [Changie](https://github.com/miniscruff/changie).

## v0.3.0 - 2021-07-15

### Added

* Implemented `properties` on `azdoext_secure_file`

### Changed

* Allow `content` and `content_base64` to be empty on `azdoext_secure_file`
* `azdoext_secure_file` now stores hash of content in state rather than full plaintext content

## v0.2.0 - 2021-07-12

### Added

* Implemented `allow_access` on the `azdoext_secure_file` resource

## v0.1.1 - 2021-07-12

### Changed

* Reduced required Azure DevOps API version to `5.0`

## v0.1.0 - 2021-07-10

### Added

* Initial implementation of `azdoext_secure_file` resource