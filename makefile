# Default target
.PHONY: help
help:
	@echo "Available make commands:"
	@echo "  build         - Build the Docker image kadlab"
	@echo "  run           - Build and run Kademlia network"
	@echo "  clean         - Stop and remove the Docker containers"
	@echo "  test          - Run all tests and the test coverage"

.PHONY: build
build:
	@echo "Building Kademlia image..."
	-docker rmi "kadlab"
	docker build -t kadlab .

.PHONY: run
run: build
	@echo "Running Docker containers..."
	docker compose up --build -d

.PHONY: clean
clean:
	@echo "Cleaning up..."
	docker compose down

.PHONY: test
test:
	@echo "Running tests and test coverage..."
	go install golang.org/x/tools/cover
	go test -v ./kademlia -cover