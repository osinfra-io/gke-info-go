# Local Values
# https://www.terraform.io/docs/language/values/locals.html

locals {
  datadog_mci_synthetic_url          = var.env == "prod" ? "https://gcp.osinfra.io/${local.datadog_synthetic_service}" : "https://${var.env}.gcp.osinfra.io/${local.datadog_synthetic_service}"
  datadog_synthetic_message_critical = var.env == "prod" ? "@hangouts-Platform-CriticalHighPriority" : ""
  datadog_synthetic_message_medium   = var.env == "prod" ? "@hangouts-Platform-MediumLowInfoPriority" : ""
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
      status           = "live"
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

      status = var.env == "sb" ? "live" : "paused"
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
      status           = "live"
      url              = var.env == "prod" ? "https://us-east1.gcp.osinfra.io/${local.datadog_synthetic_service}" : "https://us-east1.${var.env}.gcp.osinfra.io/${local.datadog_synthetic_service}"
    }
  } : {}

  registry           = var.env == "sb" ? "us-docker.pkg.dev/plt-lz-services-tf7f-sb/platform-docker-virtual" : "us-docker.pkg.dev/plt-lz-services-tf79-prod/platform-docker-virtual"
  kubernetes_project = var.env == "sb" ? "plt-k8s-tf39-sb" : var.env == "prod" ? "plt-k8s-tf10-prod" : "plt-k8s-tf33-nonprod"
}
