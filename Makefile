SHELL=/bin/bash -o pipefail
WORKSPACE=$(shell pwd)

all: help
	@true

deps: ## Installs deps
	go get -t ./...

test: deps fmt vet lint importorder staticcheck ## Runs the test suite
	go test -v -tags=int -race ./...

vet: ## Verifies all code passes a 'go vet'
	go vet -tags=int ./...

importorder: ## Verifies all code has correct import orders (stdlib, internal, 3rd party)
	impi --local github.com/graymeta/gmkit --scheme stdLocalThirdParty `go list ./...`

# once the codebase is all lintable, we can replace the for loop below with this command:
lint: ## Runs golint on all the code
	golint -set_exit_status `go list github.com/graymeta/gokit/... | grep -v vendor`

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
	staticcheck -tags=int github.com/graymeta/mf2/...

containertest:  ## The job run by Jenkins on each pull request
	docker run \
		-v $(WORKSPACE):/mnt/src/github.com/graymeta/gmkit \
		--cap-add SYS_ADMIN \
		builder-metafarm \
	/bin/bash -c "cd /mnt/src/github.com/graymeta/gmkit; PATH=/usr/local/go/bin:$$PATH GOPATH=/mnt make test"

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
