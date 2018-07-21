FROM golang:1.11.1-stretch
LABEL maintainer="Alexander Mazuruk <a.mazuruk@samsung.com>"

ENV PROJECT="github.com/SamsungSLAV/weles"

RUN go get "${PROJECT}"
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR "${GOPATH}/src/${PROJECT}"

# Copy swagger.yml customised by user.
COPY swagger.yml .

# Build swagger tool.
RUN make tools

# Regenerate server to include customized swagger.yml.
RUN make swagger-server-generate

# Build Weles server.
RUN go build -o /weles cmd/weles-server/main.go
