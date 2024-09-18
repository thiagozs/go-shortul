# Set default values
PORT ?= 9090
SUPERSCRT ?= test
DOMAIN ?= localhost
HTTPS ?= false
LOCAL ?= true

build:
	@echo "Building the Go application..."
	rm -fr url-shortener
	go build -o url-shortener .

run: build
	@echo "Running the Go application on PORT $(PORT), SUPERSCRT $(SUPERSCRT), DOMAIN $(DOMAIN)..."
	./url-shortener -port=$(PORT) \
		-token=$(SUPERSCRT) \
		-domain=$(DOMAIN) \
		-https=$(HTTPS) \
		-local=$(LOCAL)

docker-build:
	@echo "Building Docker image with PORT=$(PORT), SUPERSCRT=$(SUPERSCRT) and DOMAIN $(DOMAIN)..."
	docker build -t thiagozs/url-shortener:latest --build-arg PORT=$(PORT) --build-arg SUPERSCRT=$(SUPERSCRT) --build-arg DOMAIN=$(DOMAIN) --build-arg HTTPS=$(HTTPS) --build-arg LOCAL=$(LOCAL) .

docker-run:
	@echo "Running Docker container on PORT $(PORT), SUPERSCRT $(SUPERSCRT) and DOMAIN $(DOMAIN)..."
	# Remove existing container if it exists
	@if [ $(shell docker ps -a -q -f name=url-shortener) ]; then \
		echo "Removing existing container..."; \
		docker rm -f url-shortener; \
	fi
	# Run new container instance
	docker run -p $(PORT):$(PORT) -e PORT=$(PORT) -e SUPERSCRT=$(SUPERSCRT) -e DOMAIN=$(DOMAIN) -e LOCAL=$(LOCAL) -e HTTPS=$(HTTPS) --name url-shortener thiagozs/url-shortener:latest

clean:
	@echo "Cleaning up..."
	rm -f url-shortener
	docker rm -f url-shortener || true
	docker rmi thiagozs/url-shortener:latest || true
