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

package-test:
	@goreleaser --snapshot --skip-publish --rm-dist

package:
	@goreleaser --rm-dist
