FROM openshift/origin-release:golang-1.14 AS builder

ENV GIT_COMMITTER_NAME devtools
ENV GIT_COMMITTER_EMAIL devtools@redhat.com


WORKDIR /go/src/github.com/redhat-developer/service-binding-admission-controller

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/webhook-server ./cmd/service-binding-admission-controller

FROM registry.access.redhat.com/ubi8/ubi-minimal
COPY --from=builder /go/src/github.com/redhat-developer/service-binding-admission-controller/bin/webhook-server /usr/local/bin/webhook-server

USER 10001
ENTRYPOINT ["/usr/local/bin/webhook-server"]
