.PHONY: build run clean install

BINARY_NAME=bibChecker
GO_CMD=go

build:
	$(GO_CMD) build -o $(BINARY_NAME) ./cmd/bibChecker

run:
	$(GO_CMD) run ./cmd/bibChecker

clean:
	rm -f $(BINARY_NAME)

install:
	$(GO_CMD) install ./cmd/bibChecker
