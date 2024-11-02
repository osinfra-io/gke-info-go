# Terraform Core Helpers Module (osinfra.io)
# https://github.com/osinfra-io/terraform-core-helpers

module "helpers" {
  source = "github.com/osinfra-io/terraform-core-helpers?ref=remove-email"

  cost_center         = "x001"
  data_classification = "public"
  repository          = "google-cloud-kubernetes"
  team                = "platform-google-cloud-kubernetes"
}
