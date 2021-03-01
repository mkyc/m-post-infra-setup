# golang builder
FROM golang:1.15.2 as builder
ARG ARG_GO_MODULE_NAME="github.com/mkyc/m-host-init"
ENV GO_MODULE_NAME=$ARG_GO_MODULE_NAME
ARG ARG_M_VERSION="dev"
ENV M_VERSION=$ARG_M_VERSION
RUN mkdir -p $GOPATH/src/$GO_MODULE_NAME
COPY . $GOPATH/src/$GO_MODULE_NAME
WORKDIR $GOPATH/src/$GO_MODULE_NAME
RUN go get -v
RUN go get github.com/ahmetb/govvv
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w $(govvv -flags -pkg $GO_MODULE_NAME/cmd -version $M_VERSION)" -x -o /entrypoint $GO_MODULE_NAME

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

ARG ARG_HOST_UID=1000
ARG ARG_HOST_GID=1000
RUN chown -R $ARG_HOST_UID:$ARG_HOST_GID \
    /workdir \
    /resources

USER $ARG_HOST_UID:$ARG_HOST_GID
