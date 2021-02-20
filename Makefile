export CGO_ENABLED=0
default: release

debug:
	go build ${SERVER}
	go build ${CLIENT}

release:
	go build -ldflags "-s -w"

clean:
	@rm -f update-prometheus
	@rm -rf dist

testpkg: release
	@goreleaser --snapshot --skip-publish --rm-dist

package: release
	@goreleaser --rm-dist
