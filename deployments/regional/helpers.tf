# Terraform Core Helpers Module (osinfra.io)
# https://github.com/osinfra-io/terraform-core-helpers

module "helpers" {
  source = "github.com/osinfra-io/terraform-core-helpers?ref=v0.1.1"

  cost_center         = "x001"
  data_classification = "public"
  repository          = "google-cloud-kubernetes"
  team                = "platform-google-cloud-kubernetes"
}
