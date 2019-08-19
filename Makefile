SHELL=/bin/bash -o pipefail
WORKSPACE=$(shell pwd)

pkg_path=github.com/graymeta/gmkit

all: help
	@true

test: fmt vet lint importorder staticcheck ## Runs the test suite
	go test -v -tags=int -race ./...

vet: ## Verifies all code passes a 'go vet'
	go vet -tags=int ./...

importorder: ## Verifies all code has correct import orders (stdlib, internal, 3rd party)
	impi --local $(pkg_path) --scheme stdLocalThirdParty `go list ./...`

# once the codebase is all lintable, we can replace the for loop below with this command:
lint: ## Runs golint on all the code
	golint -set_exit_status `go list $(pkg_path)/...`

fmt: ## Verifies all code is gofmt'ed
	@STATUS=0 ; \
	for f in `find . -type f -name "*.go" | grep -v /vendor/` ; do \
		file=$$(gofmt -l $$f) ; \
		if [[ $$file ]] ; then \
			echo "file not gofmt'ed: $$f" ; \
			STATUS=1 ; \
		fi ; \
	done ; \
	if [ $$STATUS -ne 0 ] ; then \
		exit 1 ; \
	fi

staticcheck: ## runs staticcheck on our packages
	staticcheck -tags=int $(pkg_path)/...

containertest:  ## The job run by Jenkins on each pull request
	docker run \
		-v $(WORKSPACE):/mnt/src/$(pkg_path) \
		-v $(WORKSPACE)/build/run.sh:/run.sh \
		--cap-add SYS_ADMIN \
		builder-metafarm \
	/bin/bash -c "/run.sh; cd /mnt/src/$(pkg_path); PATH=/usr/local/go/bin:$$PATH GOPATH=/mnt make test"

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
