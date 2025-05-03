.PHONY: all build clean install run deps check-deps

# Go parameters
GOCMD=go
GOINSTALL=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
BINARY_NAME=server-ui
SRC_DIR=src
DIST_DIR=dist

# Check for required system packages (Fedora)
check-deps:
	@echo "Checking for required dependencies..."
	@which pkg-config >/dev/null || (echo "pkg-config is required" && exit 1)
	@pkg-config --exists gtk+-3.0 || (echo "gtk3-devel is required" && sudo dnf install -y gtk3-devel)

all: check-deps deps build

build:
	mkdir -p $(DIST_DIR)
	cd $(SRC_DIR) && $(GOINSTALL) -o ../$(DIST_DIR)/$(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -rf $(DIST_DIR)

deps:
	$(GOGET) github.com/andlabs/ui
	$(GOGET) -u -v .

run: build
	./$(DIST_DIR)/$(BINARY_NAME)

# Help target
help:
	@echo "Available targets:"
	@echo "  all        - Check dependencies, install them, and build the application (default)"
	@echo "  build      - Build and install the server UI application"
	@echo "  clean      - Clean Go cache and remove binary"
	@echo "  check-deps - Check and install required system dependencies"
	@echo "  deps       - Install Go dependencies"
	@echo "  run        - Build and run the application"