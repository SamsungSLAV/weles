BIN_DIR = bin

DEV_TOOLS_DIR = $(BIN_DIR)/dev-tools
DEV_TOOLS = ./vendor/github.com/golang/mock/mockgen ./vendor/github.com/go-swagger/go-swagger/cmd/swagger
MOCKGEN_BIN = $(DEV_TOOLS_DIR)/mockgen
SWAGGER_BIN = $(DEV_TOOLS_DIR)/swagger

DEV_TOOLS_BIN = $(MOCKGEN_BIN) $(SWAGGER_BIN)

rwildcard=$(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) $(filter $(subst *,%,$2),$d))

WELES_FILES = $(filter-out *_test.go, $(call rwildcard, , *.go))

SERVER_MAIN = cmd/weles-server/main.go
SERVER_BIN = $(BIN_DIR)/weles

BUILD_DOCKER_IMAGE = weles-build-img
BUILD_DOCKER_CONTAINER = weles-build

.PHONY: all
all: docker-build

.PHONY: clean
clean: clean-docker-build clean-vendor

.PHONY: docker-build
docker-build: $(SERVER_BIN)
	docker rm $(BUILD_DOCKER_CONTAINER)

$(SERVER_BIN): docker-container | $(BIN_DIR)
	docker cp "$(BUILD_DOCKER_CONTAINER):/$(@F)" $(@)

.PHONY: docker-container
docker-container: docker-build-img
	docker create --name "$(BUILD_DOCKER_CONTAINER)" "$(BUILD_DOCKER_IMAGE)"

.PHONY: docker-build-img
docker-build-img:
	docker build --tag "$(BUILD_DOCKER_IMAGE)" .

$(BIN_DIR):
	mkdir -p "$(BIN_DIR)"

.PHONY: clean-docker-build
clean-docker-build:
	-docker rm $(BUILD_DOCKER_CONTAINER)
	-docker rmi $(BUILD_DOCKER_IMAGE)
	-rm -f "$(SERVER_BIN)"
	-rmdir "$(BIN_DIR)"

.PHONY: server
server: vendor $(WELES_FILES)
	go build -o $(SERVER_BIN) $(SERVER_MAIN)

# dep ensure is run after swagger generation to update Gopkg.lock with packages needed to build server
.PHONY: swagger-server-generate
swagger-server-generate:  swagger.yml COPYING
	./$(DEV_TOOLS_DIR)/swagger generate server \
		-A weles \
		-f ./swagger.yml \
		-m ../weles \
		-s ./server \
		-r ./COPYING \
		--flag-strategy pflag \
		--exclude-main \
		--skip-models \
		--compatibility-mode=modern
	dep ensure

.PHONY: swagger-models-generate
swagger-models-generate:  swagger.yml COPYING
	./$(DEV_TOOLS_DIR)/swagger generate model \
		-f ./swagger.yml \
		-m ../weles \
		-s ./server \
		-r ./COPYING \

.PHONY: swagger-docs-html
swagger-docs-html: swagger.yml
	mkdir -p doc/build/swagger
	docker run \
		--rm \
		-v $$PWD:/local \
		--user `id -u $$USER`:`id -g $$USER` \
		swaggerapi/swagger-codegen-cli \
		generate \
			-i local/swagger.yml \
			-l html \
			-o /local/doc/build/swagger/

vendor: Gopkg.lock
	dep ensure -v -vendor-only

Gopkg.lock: Gopkg.toml
	dep ensure -v -no-vendor

.PHONY: dep-update
dep-update: clean-vendor
	dep ensure -update

# clean-vendor has not been added to vendor dependencies as dep is able to check out appropriate
# packages on versions set in the Gopkg.lock file. Removing vendor would force dep to re-download
# all the packages instead of only the missing ones. Global prune is turned off (see Gopkg.toml
# for explanation) thus vendor recipe will leave unused packages in the vendor/ directory. If that
# bothers you, run sequentially clean-vendor and vendor recipes.
.PHONY: clean-vendor
clean-vendor:
	rm -rf vendor

# Due to lack of standard approach to naming and separation of both interfaces and generated mock files
# below recipe does not have any file dependencies and is PHONY. Interface changes should be rare thus
# it is up to the developer to regenerate mocks after interface changes.
.PHONY: mocks
mocks: tools
	go generate ./mock
	go generate ./manager
	go generate ./controller/mock

.PHONY: tools
tools: vendor $(DEV_TOOLS_BIN)

# This recipe will rebuild all tools on vendor directory change.
# Due to short build time it is not treated as issue.
$(DEV_TOOLS_DIR)/%: $(DEV_TOOLS)
	go build -o $@ $(filter %$(@F),$(DEV_TOOLS))
