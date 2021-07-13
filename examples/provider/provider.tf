terraform {
  required_providers {
    azdoext = {
      source  = "Xtansia/azdoext"
      version = ">= 0.2.0"
    }
  }
}

provider "azdoext" {
  org_service_url       = "https://dev.azure.com/..."
  personal_access_token = "f6g..."
}