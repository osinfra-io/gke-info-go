# GKE Info Go

[![Docker Build and Test](https://github.com/osinfra-io/gke-info-go/actions/workflows/build-and-test.yml/badge.svg)](https://github.com/osinfra-io/gke-info-go/actions/workflows/build-and-test.yml) [![Docker Build and Push](https://github.com/osinfra-io/gke-info-go/actions/workflows/build-and-push.yml/badge.svg)](https://github.com/osinfra-io/gke-info-go/actions/workflows/build-and-push.yml)

## Usage

```yaml
---
apiVersion: v1
kind: Namespace

metadata:
  name: gke-info

---
apiVersion: apps/v1
kind: Deployment

metadata:
  name: gke-info-go
  namespace: gke-info

spec:
  replicas: 1
  selector:
    matchLabels:
      app: gke-info-go
      version: v1

  template:
    metadata:
      labels:
        app: gke-info-go
        version: v1

    spec:
      containers:
        - image: ghcr.io/osinfra-io/gke-info-go:latest
          imagePullPolicy: Always
          name: gke-info-go

          ports:
            - containerPort: 8080

          resources:
            limits:
              cpu: "100m"
              memory: "256Mi"
            requests:
              cpu: "50m"
              memory: "128Mi"

---
apiVersion: v1
kind: Service

metadata:
  name: gke-info-go
  namespace: gke-info

  labels:
    app: gke-info-go
    version: v1

spec:
  ports:
    - name: http
      port: 8080
      targetPort: 8080

  selector:
    app: gke-info-go

```

After deploying, you can get the information about the GKE cluster by running the following command:

```bash
kubectl port-forward --namespace gke-info $(kubectl get pod --namespace gke-info --selector="app=gke-info-go" --output jsonpath='{.items[0].metadata.name}') 8080:8080
```

Curl the endpoint:

```bash
curl http://localhost:8080/gke-info
```
