# Terraform Documentation

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_datadog"></a> [datadog](#provider\_datadog) | 3.43.1 |
| <a name="provider_google"></a> [google](#provider\_google) | 6.0.1 |
| <a name="provider_kubernetes"></a> [kubernetes](#provider\_kubernetes) | 2.32.0 |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [datadog_synthetics_test.this](https://registry.terraform.io/providers/datadog/datadog/latest/docs/resources/synthetics_test) | resource |
| [kubernetes_deployment_v1.gke_info_go](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/deployment_v1) | resource |
| [kubernetes_manifest.gke_info_go](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/manifest) | resource |
| [kubernetes_service_v1.gke_info_go](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/service_v1) | resource |
| [kubernetes_service_v1.gke_info_go_regional](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/service_v1) | resource |
| [google_client_config.current](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/client_config) | data source |
| [google_container_cluster.this](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/container_cluster) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_datadog_api_key"></a> [datadog\_api\_key](#input\_datadog\_api\_key) | Datadog API key | `string` | n/a | yes |
| <a name="input_datadog_app_key"></a> [datadog\_app\_key](#input\_datadog\_app\_key) | Datadog APP key | `string` | n/a | yes |
| <a name="input_environment"></a> [environment](#input\_environment) | The full environment name for example: `sandbox`, `non-production`, `production` | `string` | `"sandbox"` | no |
| <a name="input_gke_info_go_replicas"></a> [gke\_info\_go\_replicas](#input\_gke\_info\_go\_replicas) | The number of replicas for the gke-info deployment | `number` | `1` | no |
| <a name="input_gke_info_go_version"></a> [gke\_info\_go\_version](#input\_gke\_info\_go\_version) | The version of the gke-info deployment | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | The region to deploy the resources into | `string` | n/a | yes |
| <a name="input_zone"></a> [zone](#input\_zone) | The zone to deploy the resources to | `string` | n/a | yes |

## Outputs

No outputs.
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
