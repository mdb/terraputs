VERSION = 0.1.7
SOURCE = ./...

.PHONY: help \
	build \
	vet \
	test-fmt \
	test \
	testdata \
	clean

.DEFAULT_GOAL := build

help:
	# build:     build terraputs (default make target)
	# tools: 		 install build dependencies cited in tools.go
	# vet:       run 'go vet' against source code
	# test-fmt:  validate that source code is formatted correctly
	# test:      run automated tests
	# testdata:  generate Terraform state JSON for use in tests
	# check-tag: check if a $(VERSION) git tag already exists
	# tag:       create a $(VERSION) git tag
	# release:   build and publish a terraputs GitHub release
	# clean:     remove testdata fixtures and compiled artifacts

tools:
	echo "Installing tools from tools.go"
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

build: tools
	goreleaser release \
		--snapshot \
		--skip-publish \
		--rm-dist

vet:
	go vet $(SOURCE)

test-fmt:
	test -z $(shell go fmt $(SOURCE))

test: vet test-fmt
	go test -cover $(SOURCE) -count=1

define generate-testdata
	docker run \
		--interactive \
		--tty \
		--volume $(shell pwd):/src \
		--workdir /src/testdata/$(1) \
		--entrypoint /bin/sh \
		hashicorp/terraform:$(2) \
			-c \
				"terraform init && \
				terraform apply -auto-approve && \
				terraform show -json > show.json"
endef

testdata:
	$(call generate-testdata,basic,1.0.5)
	$(call generate-testdata,nooutputs,1.0.5)
	$(call generate-testdata,emptyconfig,1.0.5)
	$(call generate-testdata,emptyconfig-1.1.5,1.1.5)
	$(call generate-testdata,basic-1.1.5,1.1.5)

check-tag:
	./scripts/ensure_unique_version.sh "$(VERSION)"

tag: check-tag
	@echo "creating git tag $(VERSION)"
	@git tag $(VERSION)
	@git push origin $(VERSION)

release: tools
	goreleaser release \
		--rm-dist

demo:
	svg-term \
		--cast 423523 \
		--out demo.svg \
		--window \
		--no-cursor

define clean-testdata
	rm -rf testdata/$(1)/.terraform*
	rm -rf testdata/$(1)/terraform.tfstate.backup
	rm -rf testdata/$(1)/greeting.txt
endef

clean:
	$(call clean-testdata,basic)
	$(call clean-testdata,nooutputs)
	$(call clean-testdata,emptyconfig)
	rm -rf dist
