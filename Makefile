OUTPUT_DIR := binary
BINARY_NAME := termodoro

PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

build:
	go build -o $(OUTPUT_DIR)/$(BINARY_NAME)

run: build
	./$(OUTPUT_DIR)/$(BINARY_NAME)

test:
	go test ./... -v

cross-build: $(PLATFORMS)

$(PLATFORMS):
	@mkdir -p $(OUTPUT_DIR)
	GOOS=$(word 1,$(subst /, ,$@)) GOARCH=$(word 2,$(subst /, ,$@)) go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-$(word 1,$(subst /, ,$@))-$(word 2,$(subst /, ,$@))$(if $(findstring windows,$@),.exe,) .

clean:
	rm -rf $(OUTPUT_DIR)

.PHONY: build run test cross-build clean $(PLATFORMS)
