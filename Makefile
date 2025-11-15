.PHONY: test lint fmt vet ci setup

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run --timeout=5m ./...

fmt:
	go fmt ./...

fmt-check:
	@if [ "$$(gofmt -s -l . | wc -l)" -gt 0 ]; then \
		echo "Code is not formatted. Run 'make fmt'"; \
		gofmt -s -d .; \
		exit 1; \
	fi

vet:
	go vet ./...

ci: lint test fmt-check vet

setup:
	@echo "Installing golangci-lint..."
	@if [ "$$(uname)" = "Darwin" ]; then \
		brew install golangci-lint; \
	else \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.60.0; \
	fi
	@echo "Installation completed!"