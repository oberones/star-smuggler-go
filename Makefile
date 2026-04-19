SHELL := /bin/bash

ROOT_DIR := $(CURDIR)
GO ?= go
DOTNET ?= dotnet
GODOT_BIN ?= /Applications/Godot\ .NET.app/Contents/MacOS/Godot

.PHONY: help test test-go test-integration test-golden run-go build-dotnet smoke smoke-go smoke-travel godot-headless godot-open fmt fmt-go

help:
	@printf "Available targets:\n"
	@printf "  %-18s %s\n" "help" "Show this help message"
	@printf "  %-18s %s\n" "test" "Run the main Go test suite"
	@printf "  %-18s %s\n" "test-go" "Alias for test"
	@printf "  %-18s %s\n" "test-integration" "Run integration tests only"
	@printf "  %-18s %s\n" "test-golden" "Run golden tests only"
	@printf "  %-18s %s\n" "run-go" "Run the Go bootstrap entrypoint"
	@printf "  %-18s %s\n" "build-dotnet" "Build the Godot C# project"
	@printf "  %-18s %s\n" "smoke" "Run the default headless Godot smoke test"
	@printf "  %-18s %s\n" "smoke-go" "Run the headless smoke test with STARSMUGGLER_GO_RUNTIME=1"
	@printf "  %-18s %s\n" "smoke-travel" "Run the travel-flow smoke check"
	@printf "  %-18s %s\n" "godot-headless" "Launch Godot headless directly"
	@printf "  %-18s %s\n" "godot-open" "Open the project in the .NET-enabled Godot editor"
	@printf "  %-18s %s\n" "fmt" "Format Go source files"
	@printf "  %-18s %s\n" "fmt-go" "Alias for fmt"

test: test-go

test-go:
	$(GO) test ./...

test-integration:
	$(GO) test ./tests/integration/...

test-golden:
	$(GO) test ./tests/golden/...

run-go:
	$(GO) run ./cmd/starsmuggler

build-dotnet:
	$(DOTNET) build StarSmugglerGo.sln

smoke:
	bash tests/smoke/run_headless_smoke.sh

smoke-go:
	STARSMUGGLER_GO_RUNTIME=1 bash tests/smoke/run_headless_smoke.sh

smoke-travel:
	bash tests/smoke/travel_flow_check.sh

godot-headless:
	"$(GODOT_BIN)" --headless --path "$(ROOT_DIR)" --quit-after 1

godot-open:
	open -a "/Applications/Godot .NET.app" "$(ROOT_DIR)"

fmt: fmt-go

fmt-go:
	@files="$$(find cmd internal tests -name '*.go' -type f | sort)"; \
	if [ -z "$$files" ]; then \
		echo "No Go files found."; \
	else \
		$(GO) fmt $$files; \
	fi
