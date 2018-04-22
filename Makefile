PROJECT_PATH 		:= $(basename $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
ifeq ($(PROJECT_PATH),)
$(error Variables DIR_PROJECT not set correctly)
endif
PROJECT_NAME 		:= $(lastword $(subst /, ,$(PROJECT_PATH)))
PROJECT_PACKAGE		:= $(patsubst %/,%,$(subst $(GOPATH)/src/,,$(PROJECT_PATH)))

PACKAGES			:= $(shell go list $(PROJECT_PACKAGE)/... | grep -v vendor/)

DIR_VENDOR			:= $(PROJECT_PATH)vendor/
DIR_BUILD			:= $(PROJECT_PATH)build/
DIR_COVERAGE		:= $(DIR_BUILD)coverage
DIR_RELEASES		:= $(DIR_BUILD)releases

.DEFAULT_GOAL		:= help

.PHONY: dep
dep: ## Update dependencies
	dep ensure -v -update -no-vendor
	dep ensure -v -vendor-only

.PHONY: setup
setup: clean ## Install all the project's tools and dependencies
	go get -v -u github.com/alecthomas/gometalinter
	go get -v -u golang.org/x/tools/cmd/cover
	go get -v -u github.com/golang/dep/cmd/dep
	gometalinter --install --update
	@$(MAKE) dep

.PHONY: test-lint
test-lint: ## Run linters against package and sub-packages
	gometalinter \
		--no-config \
		--deadline=2m \
		--enable-gc \
		--line-length=100 \
		--vendor \
		--disable-all \
		--enable=deadcode \
		--enable=dupl \
		--enable=errcheck \
		--enable=gas \
		--enable=goconst \
		--enable=gocyclo \
		--enable=goimports \
		--enable=golint \
		--enable=gotype \
		--enable=gotypex \
		--enable=ineffassign \
		--enable=interfacer \
		--enable=lll \
		--enable=maligned \
		--enable=megacheck \
		--enable=misspell \
		--enable=nakedret \
		--enable=safesql \
		--enable=structcheck \
		--enable=unconvert \
		--enable=unparam \
		--enable=varcheck \
		--enable=vet \
		./...

.PHONY: test-dep
test-dep: ## Check for useless and/or missing dependencies
	dep status
	@test $(shell dep status 2>&1 | grep -c "MISSING PACKAGES") -eq 0

.PHONY: test
test: ## Run the unit/functionnal tests
	mkdir -p $(DIR_COVERAGE)
	@echo 'mode: atomic' > $(DIR_COVERAGE)/coverage.txt
	go test -covermode=atomic -coverprofile=$(DIR_COVERAGE)/coverage.txt -v -race -timeout=30s $(BUILD_TAGS) $(PACKAGES)

.PHONY: test-all
test-all: test-lint test-dep test ## Run all kind of tests (code quality, missing dependencies, units, functionnal)

.PHONY: cover
cover: test ## Run tests to compute coverage and open the coverage report
	mkdir -p $(DIR_COVERAGE)
	go tool cover -html=$(DIR_COVERAGE)/coverage.txt

.PHONY: clean
clean: ## Remove vendors, build, and temporary files
	$(RM) -r $(DIR_BUILD)
	$(RM) -r $(DIR_VENDOR)
	go clean

.PHONY: ci
ci: setup test-dep test-lint test ## Useful alias for CI

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@echo "Available targets descriptions:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
