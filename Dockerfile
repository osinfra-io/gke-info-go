FROM golang:1.22.3

ARG DD_GIT_REPOSITORY_URL
ARG DD_GIT_COMMIT_SHA
ENV DD_GIT_REPOSITORY_URL=${DD_GIT_REPOSITORY_URL}
ENV DD_GIT_COMMIT_SHA=${DD_GIT_COMMIT_SHA}

WORKDIR /app

COPY main.go .

RUN go mod init gke-info
RUN go mod tidy
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
