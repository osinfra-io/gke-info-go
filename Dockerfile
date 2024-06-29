FROM golang:1.22.4

# Ensure that HEALTHCHECK instructions have been added to container images
# checkov:skip=CKV_DOCKER_2: Since Kubernetes 1.8, the Docker HEALTHCHECK has been disabled explicitly

# Create a non-root user for security purposes

RUN useradd -m gke-info
USER gke-info

ARG DD_GIT_REPOSITORY_URL
ARG DD_GIT_COMMIT_SHA
ENV DD_GIT_REPOSITORY_URL=${DD_GIT_REPOSITORY_URL}
ENV DD_GIT_COMMIT_SHA=${DD_GIT_COMMIT_SHA}

# Set the working directory inside the container

WORKDIR /app

# Pre-copy/cache go.mod for pre-downloading dependencies and only re-downloading them in subsequent builds if they change

COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy

COPY cmd/ /app/cmd/
COPY internal/ /app/internal/

# Build the application

RUN GOOS=linux CGO_ENABLED=0 go build -o main cmd/http/main.go

# Expose the port your application listens on

EXPOSE 8080

# Define the command to run your application

CMD ["./main"]
