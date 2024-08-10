variable "env" {
  description = "This is the environment suffix for example: `sb` (Sandbox), `non-prod` (Non-Production), `prod` (Production)"
  type        = string
  default     = "sb"
}

variable "environment" {
  description = "The full environment name for example: `sandbox`, `non-production`, `production`"
  type        = string
  default     = "sandbox"
}

variable "datadog_api_key" {
  description = "Datadog API key"
  type        = string
}

variable "datadog_app_key" {
  description = "Datadog APP key"
  type        = string
}

variable "gke_info_go_replicas" {
  description = "The number of replicas for the gke-info deployment"
  type        = number
  default     = 1
}

variable "gke_info_go_version" {
  description = "The version of the gke-info deployment"
  type        = string
}

variable "region" {
  description = "The region to deploy the resources into"
  type        = string
}

variable "zone" {
  description = "The zone to deploy the resources to"
  type        = string
}
