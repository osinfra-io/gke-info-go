variable "datadog_api_key" {
  description = "Datadog API key"
  type        = string
}

variable "datadog_app_key" {
  description = "Datadog APP key"
  type        = string
}

variable "istio_test_replicas" {
  description = "The number of replicas for the istio-test deployment"
  type        = number
  default     = 1
}

variable "istio_test_version" {
  description = "The version of the istio-test deployment"
  type        = string
}
