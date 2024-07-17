# GKE Info Go

An example Go application that shows information about the Google Kubernetes Engine (GKE) cluster.

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

  template:
    metadata:
      labels:
        app: gke-info-go

    spec:
      containers:
        - image: ghcr.io/osinfra-io/gke-info-go:latest
          imagePullPolicy: Always
          name: gke-info-go

          ports:
            - containerPort: 8080

          resources:
            limits:
              cpu: "50m"
              memory: "128Mi"
            requests:
              cpu: "25m"
              memory: "64Mi"

---
apiVersion: v1
kind: Service

metadata:
  name: gke-info-go
  namespace: gke-info

  labels:
    app: gke-info-go

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
