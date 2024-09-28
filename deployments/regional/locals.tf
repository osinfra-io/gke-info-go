# Local Values
# https://www.terraform.io/docs/language/values/locals.html

locals {
  datadog_mci_synthetic_url          = var.environment == "production" ? "https://gcp.osinfra.io/${local.datadog_synthetic_service}/metadata/cluster-name" : "https://${local.env}.gcp.osinfra.io/${local.datadog_synthetic_service}/metadata/cluster-name"
  datadog_synthetic_message_critical = var.environment == "production" ? "@hangouts-Platform-CriticalHighPriority" : ""
  datadog_synthetic_message_medium   = var.environment == "production" ? "@hangouts-Platform-MediumLowInfoPriority" : ""
  datadog_synthetic_name             = "GKE Info"
  datadog_synthetic_service          = "gke-info-go"

  datadog_synthetic_tests = var.region == "us-east1" || var.zone == "b" ? {
    "mci" = {
      locations = [
        "aws:ca-central-1",
        "aws:us-west-1",
        "aws:us-east-1",
        "aws:eu-central-1",
        "aws:eu-south-1",
        "aws:eu-north-1",
        "aws:us-east-1",
        "aws:us-east-2",
        "aws:us-west-1",
        "aws:us-west-2"
      ]

      message          = local.datadog_synthetic_message_critical
      message_priority = "1"
      name             = "Istio MCI ${local.datadog_synthetic_name}"
      region           = "global"
      service          = local.datadog_synthetic_service
      status           = "paused"
      url              = local.datadog_mci_synthetic_url
    }

    "mci-us-east1" = {
      locations = [
        "aws:us-east-1",
        "aws:us-east-2",
        "aws:us-west-1",
        "aws:us-west-2"
      ]

      message          = local.datadog_synthetic_message_medium
      message_priority = "3"
      name             = "Istio MCI ${local.datadog_synthetic_name}"
      region           = "us-east1"
      service          = local.datadog_synthetic_service

      status = "paused"
      url    = local.datadog_mci_synthetic_url
    }

    "us-east1" = {
      locations = [
        "aws:us-east-1",
        "aws:us-east-2",
        "aws:us-west-1",
        "aws:us-west-2"
      ]

      message          = local.datadog_synthetic_message_medium
      message_priority = "3"
      name             = "Istio Ingress ${local.datadog_synthetic_name}"
      region           = "us-east1"
      service          = local.datadog_synthetic_service
      status           = "paused"
      url              = var.environment == "production" ? "https://us-east1.gcp.osinfra.io/${local.datadog_synthetic_service}" : "https://us-east1.${local.env}.gcp.osinfra.io/${local.datadog_synthetic_service}"
    }
  } : {}

  env_map = {
    "sandbox"        = "sb"
    "non-production" = "nonprod"
    "production"     = "prod"
  }

  env = lookup(local.env_map, var.environment, "none")

  registry           = var.environment == "sandbox" ? "us-docker.pkg.dev/plt-lz-services-tf7f-sb/plt-docker-virtual" : "us-docker.pkg.dev/plt-lz-services-tf79-prod/plt-docker-virtual"
  kubernetes_project = var.environment == "sandbox" ? "plt-k8s-tf39-sb" : var.environment == "production" ? "plt-k8s-tf10-prod" : "plt-k8s-tf33-nonprod"
}
