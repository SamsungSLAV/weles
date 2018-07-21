FROM golang:1.11.1-stretch
LABEL maintainer="Alexander Mazuruk <a.mazuruk@samsung.com>"

ENV PROJECT="github.com/SamsungSLAV/weles"

RUN go get -d "${PROJECT}"
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR "${GOPATH}/src/${PROJECT}"

# Copy swagger.yml customized by user.
# Only instance-specific values should be changed:
# * info section - contact email and terms of service,
# * host section.
COPY swagger.yml .

# Build swagger tool.
RUN make tools

# Regenerate server to include customized swagger.yml.
RUN make swagger-server-generate

# Build Weles server.
RUN go build -o /weles cmd/weles-server/main.go
