
VERSION=$(shell awk '/([0-9]{1}.?){3}/ {print $$4;}' main.go)

build/gotop:
	@go build

build/nfpm.rpm:
	@docker run --rm \
	-v "$(PWD)/build:/tmp/pkg" \
	-e "VERSION=$(VERSION)" \
	goreleaser/nfpm pkg \
		--config /tmp/pkg/nfpm.yaml \
		--target /tmp/pkg/nfpm.rpm
