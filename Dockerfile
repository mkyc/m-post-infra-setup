# golang builder
FROM golang:1.15.2 as builder

ARG ARG_GO_MODULE_NAME
ENV ARG_GO_MODULE_NAME=${ARG_GO_MODULE_NAME:-github.com/epiphany-platform/m-host-init}
ARG ARG_M_VERSION
ENV ARG_M_VERSION=${ARG_M_VERSION:-dev}

RUN mkdir -p $GOPATH/src/$ARG_GO_MODULE_NAME
COPY . $GOPATH/src/$ARG_GO_MODULE_NAME
WORKDIR $GOPATH/src/$ARG_GO_MODULE_NAME

RUN go get -v &&\
  go get github.com/ahmetb/govvv &&\
  CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w $(govvv -flags -pkg $ARG_GO_MODULE_NAME/cmd -version $ARG_M_VERSION)" -x -o /entrypoint $ARG_GO_MODULE_NAME

# main
FROM quay.io/ansible/ansible-runner:stable-2.10-devel

COPY resources/requirements.yml .
RUN ansible-galaxy install -r requirements.yml

ENV RESOURCES "/resources"
ENV SHARED "/shared"

WORKDIR /workdir
ENTRYPOINT ["/workdir/entrypoint"]

COPY resources/project /resources/project
COPY --from=builder /entrypoint /workdir

# TODO change source image to be able to change user (additional context: https://github.com/ansible/ansible-runner/issues/611)
#ARG ARG_HOST_UID=1000
#ARG ARG_HOST_GID=1000
#RUN chown -R $ARG_HOST_UID:$ARG_HOST_GID /workdir /resources ~/.ansible/
#
#USER $ARG_HOST_UID:$ARG_HOST_GID
