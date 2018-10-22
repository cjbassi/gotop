
VERSION=$(shell awk '/([0-9]{1}.?){3}/ {print $$4;}' main.go)

.PHONY: all
all: pkg/gotop.rpm pkg/gotop.deb

build/gotop:
	@GOOS=linux GOARCH=amd64 go build -o $@

pkg:
	@mkdir $@

pkg/gotop.rpm: pkg build/gotop
	@docker run --rm \
	-v "$(PWD)/build:/tmp/pkg" \
	-e "VERSION=$(VERSION)" \
	goreleaser/nfpm pkg \
		--config /tmp/pkg/gotop-nfpm.yaml \
		--target /tmp/pkg/gotop.rpm \
	&& mv ./build/gotop.rpm $@ 

pkg/gotop.deb: pkg build/gotop
	@docker run --rm \
	-v "$(PWD)/build:/tmp/pkg" \
	-e "VERSION=$(VERSION)" \
	goreleaser/nfpm pkg \
		--config /tmp/pkg/gotop-nfpm.yaml \
		--target /tmp/pkg/gotop.deb \
	&& mv ./build/gotop.deb $@ 

.PHONY: clean
clean:
	@-rm -f build/gotop
	@-rm -rf pkg