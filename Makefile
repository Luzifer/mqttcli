publish:
	bash ./ci/build.sh

build:
	go build \
		-buildvcs=false \
		-ldflags "-s -w -buildid= -X main.version=$(PRODUCT_VERSION)" \
		-mod=readonly \
		-trimpath

# -- Vulnerability scanning --

trivy:
	trivy fs . \
		--dependency-tree \
		--exit-code 1 \
		--format table \
		--ignore-unfixed \
		--quiet \
		--scanners config,license,secret,vuln \
		--severity HIGH,CRITICAL \
		--skip-dirs docs
