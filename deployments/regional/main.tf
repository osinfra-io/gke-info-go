# Datadog Synthetics Test Resource
# https://registry.terraform.io/providers/DataDog/datadog/latest/docs/resources/synthetics_test

resource "datadog_synthetics_test" "this" {
  for_each = local.datadog_synthetic_tests

  assertion {
    type     = "statusCode"
    operator = "is"
    target   = "200"
  }

  assertion {
    type     = "responseTime"
    operator = "lessThan"
    target   = 1000
  }

  locations = each.value.locations
  message   = each.value.message
  name      = "${each.value.name} ${each.value.region} ${module.helpers.environment}"

  options_list {
    tick_every = 300

    retry {
      count    = 2
      interval = 120
    }

    monitor_priority = each.value.message_priority
  }

  request_definition {
    method = "GET"
    url    = each.value.url
  }

  request_headers = each.value.region == "global" ? {} : {
    Body = "services-${each.value.region}"
  }

  status  = each.value.status
  subtype = "http"

  tags = [
    "env:${module.helpers.environment}",
    "service:${each.value.service}",
    "region:${each.value.region}",
    "team:${module.helpers.team}"
  ]

  type = "api"
}



# Kubernetes Deployment Resource
# https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/deployment_v1

# This is a simple deployment that is used to get cluster status information and test end-to-end connectivity to the cluster
# through Datadog synthetic tests.

resource "kubernetes_deployment_v1" "gke_info_go" {
  # Minimize the admission of containers with the NET_RAW capability
  # checkov:skip=CKV_K8S_28: This needs some additional investigation

  # Apply security context to your pods, deployments and daemon_set
  # checkov:skip=CKV_K8S_29: A user is set in the container and we do not need to override it

  # Apply security context to your pods and containers
  # checkov:skip=CKV_K8S_30: A user is set in the container and we do not need to override it

  # Image should use digest
  # checkov:skip=CKV_K8S_43: We are using the image tag for the deployment

  metadata {
    labels = {
      "tags.datadoghq.com/env"     = module.helpers.environment
      "tags.datadoghq.com/service" = "gke-info-go"
      "tags.datadoghq.com/version" = var.gke_info_go_version
    }

    name      = "gke-info-go"
    namespace = "gke-info"
  }

  spec {
    replicas = var.gke_info_go_replicas

    selector {
      match_labels = {
        "app" = "gke-info-go"
      }
    }

    template {
      metadata {
        annotations = {
          "apm.datadoghq.com/env" = jsonencode({
            "DD_ENV"     = module.helpers.environment
            "DD_SERVICE" = "gke-info-go"
            "DD_VERSION" = var.gke_info_go_version
          })
        }

        labels = {
          # Enable Admission Controller to mutate new pods part of this deployment
          "admission.datadoghq.com/enabled" = "true"
          "app"                             = "gke-info-go"
          "tags.datadoghq.com/env"          = module.helpers.environment
          "tags.datadoghq.com/service"      = "gke-info-go"
          "tags.datadoghq.com/version"      = var.gke_info_go_version
        }
      }

      spec {
        container {
          env {
            name  = "DD_APPSEC_ENABLED"
            value = "true"
          }
          env {
            name = "DD_ENV"
            value_from {
              field_ref {
                field_path = "metadata.labels['tags.datadoghq.com/env']"
              }
            }
          }
          env {
            name = "DD_SERVICE"
            value_from {
              field_ref {
                field_path = "metadata.labels['tags.datadoghq.com/service']"
              }
            }
          }
          env {
            name = "DD_VERSION"
            value_from {
              field_ref {
                field_path = "metadata.labels['tags.datadoghq.com/version']"
              }
            }
          }

          name              = "gke-info-go"
          image             = "${local.registry}/gke-info-go:${var.gke_info_go_version}"
          image_pull_policy = "Always"

          resources {
            requests = {
              cpu    = "10m"
              memory = "32Mi"
            }
            limits = {
              cpu    = "20m"
              memory = "64Mi"
            }
          }

          port {
            container_port = 8080
          }

          liveness_probe {
            http_get {
              path = "/gke-info-go/health"
              port = "8080"
            }

            initial_delay_seconds = 10
            timeout_seconds       = 5
            period_seconds        = 10
            failure_threshold     = 5
          }

          readiness_probe {
            http_get {
              path = "/gke-info-go/health"
              port = "8080"
            }

            initial_delay_seconds = 10
            timeout_seconds       = 5
            period_seconds        = 10
            failure_threshold     = 5
          }
        }

        topology_spread_constraint {
          label_selector {
            match_labels = {
              "app" = "gke-info-go"
            }
          }

          max_skew           = 1
          topology_key       = "kubernetes.io/zone"
          when_unsatisfiable = "ScheduleAnyway"
        }
      }
    }
  }
}

# Kubernetes Manifest Resource
# https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/manifest

resource "kubernetes_manifest" "gke_info_go" {
  manifest = {
    apiVersion = "security.istio.io/v1"
    kind       = "AuthorizationPolicy"

    metadata = {
      name      = "gke-info-go"
      namespace = "gke-info"
    }

    spec = {
      action = "ALLOW"
      rules = [
        {
          from = [
            {
              source = {
                namespaces = ["istio-ingress"]
              }
            }
          ]

          to = [
            {
              operation = {
                methods = ["GET"]

                # The authorization policy below uses the ALLOW-with-positive-matching pattern to allow requests to specific paths.

                paths = [
                  "/gke-info-go/health",
                  "/gke-info-go/metadata/cluster-location",
                  "/gke-info-go/metadata/cluster-name",
                  "/gke-info-go/metadata/instance-zone"
                ]
              }
            }
          ]
        }
      ]

      selector = {
        matchLabels = {
          app = "gke-info-go"
        }
      }
    }
  }
}

# Kubernetes Service Resource
# https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/service_v1

resource "kubernetes_service_v1" "gke_info_go" {
  metadata {
    name      = "gke-info-go"
    namespace = "gke-info"
  }

  spec {
    type = "ClusterIP"
    selector = {
      app = "gke-info-go"
    }

    port {
      name        = "http"
      port        = 8080
      target_port = 8080
    }
  }
}

resource "kubernetes_service_v1" "gke_info_go_regional" {
  metadata {
    name      = "gke-info-go-${module.helpers.region}-${module.helpers.zone}"
    namespace = "gke-info"
  }

  spec {
    type = "ClusterIP"
    selector = {
      app = "gke-info-go"
    }

    port {
      name        = "http"
      port        = 8080
      target_port = 8080
    }
  }
}
