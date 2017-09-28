GO_ROOT="/usr/local/go"

# run code analysis
define analysis_only
    @echo "# running code analysis"
    @gometalinter --vendor --exclude=vendor --deadline 1m --enable-gc --disable=aligncheck ./...
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

# install dependecies
install-dep:
	$(call install_dependencies)
	$(fix_sirupsen_case_sensitivity_problem)

# update dependencies
update-dep:
	$(call update_dependencies)
	$(fix_sirupsen_case_sensitivity_problem)

# run all targets
all:
	$(call analysis_only)

.PHONY:analysis install-tools install-dep update-dep
