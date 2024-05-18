# GKE Info Go

[![Docker Build and Test](https://github.com/osinfra-io/gke-info-go/actions/workflows/build-and-test.yml/badge.svg)](https://github.com/osinfra-io/gke-info-go/actions/workflows/build-and-test.yml)

## Usage

```yaml
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
        imagePullPolicy: IfNotPresent
        name: gke-info-go
        ports:
        - containerPort: 8080
      imagePullSecrets:
       - name: github-container-registry-key

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

After deploying the above, you can check the status of the cluster by running:

```bash
kubectl port-forward --namespace gke-info $(kubectl get pod --namespace gke-info --selector="app=gke-info-go" --output jsonpath='{.items[0].metadata.name}') 8080:8080
```

Open your browser to <http://localhost:8080/cluster-name> and you should see the name of the cluster.
