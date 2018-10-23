VERSION=$(shell awk '/([0-9]{1}.?){3}/ {print $$4;}' main.go)

.PHONY: default
default: dist/gotop.rpm dist/gotop.deb

dist/gotop:
	@GOOS=linux GOARCH=amd64 go build -o $@

dist:
	@mkdir $@

dist/gotop.rpm: dist dist/gotop
	@docker run --rm \
	-v "$(PWD)/build:/tmp/build" \
	-v "$(PWD)/dist:/tmp/dist" \
	-e "VERSION=$(VERSION)" \
	goreleaser/nfpm pkg \
		--config /tmp/build/gotop-nfpm.yml \
		--target /tmp/dist/gotop.rpm

dist/gotop.deb: dist dist/gotop
	@docker run --rm \
	-v "$(PWD)/build:/tmp/build" \
	-v "$(PWD)/dist:/tmp/dist" \
	-e "VERSION=$(VERSION)" \
	goreleaser/nfpm pkg \
		--config /tmp/build/gotop-nfpm.yml \
		--target /tmp/dist/gotop.deb

.PHONY: clean
clean:
	@-rm -rf dist
