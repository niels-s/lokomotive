package aws

var terraformConfigTmpl = `
module "aws-{{.Config.ClusterName}}" {
  source = "{{.Source}}"

  providers = {
    aws      = "aws.default"
    local    = "local.default"
    null     = "null.default"
    template = "template.default"
    tls      = "tls.default"
  }

  cluster_name = "{{.Config.ClusterName}}"
  dns_zone     = "{{.Config.DNSZone}}"
  dns_zone_id  = "{{.Config.DNSZoneID}}"

  ssh_authorized_key = "{{.SSHAuthorizedKey}}"
  asset_dir          = "{{.Config.AssetDir}}"

  controller_count = "{{.Config.ControllerCount}}"
  controller_type  = "{{.Config.ControllerType}}"

  worker_count = "{{.Config.WorkerCount}}"
  worker_type  = "{{.Config.WorkerType}}"

  os_image = "{{.Config.OSImage}}"
}

provider "aws" {
  version = "~> 1.13.0"
  alias   = "default"

  region                  = "eu-central-1"
  shared_credentials_file = "{{.Config.CredsPath}}"
}

provider "local" {
  version = "~> 1.0"
  alias   = "default"
}

provider "null" {
  version = "~> 1.0"
  alias   = "default"
}

provider "template" {
  version = "~> 1.0"
  alias   = "default"
}

provider "tls" {
  version = "~> 1.0"
  alias   = "default"
}
`
