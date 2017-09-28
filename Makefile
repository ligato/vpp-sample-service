GO_ROOT="/usr/local/go"

# run code analysis
define analysis_only
    @echo "# running code analysis"
    @gometalinter --vendor --exclude=vendor --deadline 1m --enable-gc --disable=aligncheck ./...
    @echo "# done"
endef

# build vpp-agent plugin (bgp-vpp-agent)
define build_bgpplugin
    @echo "# building bgpplugin"
    @cd agent && ${GO_ROOT}/bin/go build -v ${LDFLAGS}
    @echo "# done"
endef

# build examples
define build_example
    @echo "# building examples"
    @cd examples && ${GO_ROOT}/bin/go build -v ${LDFLAGS} end_to_end_example.go
    @echo "# done"
endef

# clean bgpplugin
define clean_bgpplugin
    @echo "# cleaning bgpplugin"
    @rm -f bgpplugin/plugin
    @echo "# done"
endef

# install dependencies according to glide.yaml & glide.lock (in case vendor dir was deleted)
define install_dependencies
	$(if $(shell command -v glide install 2> /dev/null),$(info glide dependency manager is ready),$(error glide dependency manager missing, info about installation can be found here https://github.com/Masterminds/glide))
	@echo "# installing dependencies, please wait ..."
	@glide install --strip-vendor
endef

# clean update dependencies according to glide.yaml (re-downloads all of them)
define update_dependencies
	$(if $(shell command -v glide install 2> /dev/null),$(info glide dependency manager is ready),$(error glide dependency manager missing, info about installation can be found here https://github.com/Masterminds/glide))
	@echo "# updating dependencies, please wait ..."
	@-cd vendor && rm -rf *
	@echo "# vendor dir cleared"
	@-rm -rf glide.lock
	@glide cc
	@echo "# glide cache cleared"
	@glide install --strip-vendor
endef

# fix sirupsen/Sirupsen problem
define fix_sirupsen_case_sensitivity_problem
    @echo "# fixing sirupsen case sensitivity problem, please wait ..."
    @-rm -rf vendor/github.com/Sirupsen
    @-find ./ -type f -name "*.go" -exec sed -i -e 's/github.com\/Sirupsen\/logrus/github.com\/sirupsen\/logrus/g' {} \;
endef


# build all binaries
build:
	$(call build_bgpplugin)
	$(call build_example)

# get tools (analysis,mocking,...)
install-tools:
	@${GO_ROOT}/bin/go get -u -f "github.com/alecthomas/gometalinter"
	@gometalinter --install
	@${GO_ROOT}/bin/go get -u -f "github.com/golang/mock/gomock"
	@${GO_ROOT}/bin/go get -u -f "github.com/golang/mock/mockgen"
	@${GO_ROOT}/bin/go install "github.com/golang/mock/mockgen"

# run & print code analysis
analysis:
	$(call analysis_only)

# clean
clean:
	$(call clean_bgpplugin)

# install dependecies
install-dep:
	$(call install_dependencies)
	$(fix_sirupsen_case_sensitivity_problem)

# update dependencies
update-dep:
	$(call update_dependencies)
	$(fix_sirupsen_case_sensitivity_problem)

# generate plugin mock for tests
generate-test-mocks:
	    @mockgen -source=vendor/github.com/ligato/vpp-agent/clientv1/defaultplugins/data_change_api.go -destination=mocks/data_change_api.go -package=mocks -imports .=github.com/ligato/vpp-agent/clientv1/defaultplugins
	    @mockgen -source=vendor/github.com/ligato/vpp-agent/clientv1/defaultplugins/data_resync_api.go -destination=mocks/data_resync_api.go -package=mocks -imports .=github.com/ligato/vpp-agent/clientv1/defaultplugins

# run all tests
test:
	@echo "# running unit tests"
	@go test $$(go list ./... | grep -v /vendor/)

# get coverage percentage
coverage:
	@echo "# getting test coverage"
	@go test -cover $$(go list ./... | grep -v /vendor/)

# run all targets
all:
	$(call analysis_only)
	$(call build_bgpplugin)
	$(call build_example)

.PHONY: build analysis clean install-tools install-dep update-dep test coverage
