include Makeroutines.mk

# run code analysis
define analysis_only
	@echo "# running code analysis"
	@gometalinter --vendor --exclude=vendor --deadline 1m --enable-gc --disable=aligncheck --disable=gotype --disable=gotypex --exclude=mock ./...
	@echo "# done"
endef

# build vpp-agent plugin (bgp-vpp-agent)
define build_bgpplugin
    @echo "# building bgpplugin"
    @cd plugins && go build -a -v ${LDFLAGS}
    @echo "# done"
endef

# build examples
define build_example
    @echo "# building examples"
    @cd examples && go build -v ${LDFLAGS} end_to_end_example.go
    @echo "# done"
endef

# verify that links in markdown files are valid
# requires npm install -g markdown-link-check
define check_links_only
    @echo "# checking links"
    @./scripts/check_links.sh
    @echo "# done"
endef

# run all tests with coverage
define test_cover_only
	@echo "# running unit tests with coverage analysis"
	@go test -covermode=count -coverprofile=${COVER_DIR}coverage.out ./plugins/vppl3bgp
    @echo "# coverage data generated into ${COVER_DIR}coverage.out"
    @echo "# done"
endef

# run all targets
all:
	@echo "# running all"
	@make install-tools
	@make install-dep
	@make update-dep
	@make analysis
	@make build
	@make test-cover
	@make run-examples
	@make clean-examples

# build all binaries
build:
	@echo "# building"
	@go build -a ./plugins/...
	@echo "# done"

# run & print code analysis
analysis:
	$(call analysis_only)

# get tools (analysis,mocking,...)
install-tools:
	@go get -u -f "github.com/alecthomas/gometalinter"
	@gometalinter --install
	@go get -u -f "github.com/golang/mock/gomock"
	@go get -u -f "github.com/golang/mock/mockgen"
	@go install "github.com/golang/mock/mockgen"

# install dependecies
install-dep:
	$(call install_dependencies)
	$(fix_sirupsen_case_sensitivity_problem)

# update dependencies
update-dep:
	$(call update_dependencies)
	$(fix_sirupsen_case_sensitivity_problem)

# get coverage percentage
coverage:
	@echo "# getting test coverage"
	@go test -cover $$(go list ./plugins/... | grep -v /vendor/)

# run all tests
test:
	@echo "# running unit tests"
	@go test $$(go list ./... | grep -v /vendor/)

# run tests with coverage report
test-cover:
	$(call test_cover_only)

# run examples
run-examples:
	@echo "# running examples"
	@./scripts/run_vpp_l3_bgp_routes.sh
	@echo "# done"

# run clean examples
clean-examples:
	@rm -f examples/gobgp_watch_plugin/gobgp_watch_plugin
	@rm -f docker/gobgp_route_reflector/gobgp-client-in-docker/log
	@rm -f log
	@rm -f fib

# validate links in markdown files
check_links:
	$(call check_links_only)

.PHONY: build analysis install-tools install-dep update-dep run-examples all test test-cover clean-examples coverage check_links
