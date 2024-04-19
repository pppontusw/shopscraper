# Go package and binary names
API_PKG = ./cmd/api
MAILER_PKG = ./cmd/mailer
SCRAPER_PKG = ./cmd/scraper
API_BIN = api
MAILER_BIN = mailer
SCRAPER_BIN = scraper

all: build
	
.PHONY: start-dependencies
start-dependencies:
	docker-compose up -d postgres

.PHONY: test
test: start-dependencies
	go test ./...

.PHONY: test-ci
test-ci:
	go test ./...

.PHONY: test-frontend
test-frontend:
	cd frontend && npm test -- --watchAll=false

.PHONY: watch-frontend
watch-frontend:
	cd frontend && npm test

.PHONY: coverage
coverage: start-dependencies
	go test -cover ./...

.PHONY: build
build:
	go build -o $(API_BIN) $(API_PKG)
	go build -o $(MAILER_BIN) $(MAILER_PKG)
	go build -o $(SCRAPER_BIN) $(SCRAPER_PKG)

.PHONY: build-frontend
build-frontend:
	cd frontend && npm install && npm run build

.PHONY: deploy-frontend
deploy-frontend: build-frontend
	@if [ -z "$$S3_BUCKET_NAME" ]; then \
		echo "S3_BUCKET_NAME is not set. Aborting."; \
		exit 1; \
	fi
	@if [ -z "$$DISTRIBUTION_ID" ]; then \
		echo "DISTRIBUTION_ID is not set. Aborting."; \
		exit 1; \
	fi
	@if [ -z "$$REACT_APP_API_URL" ]; then \
		echo "REACT_APP_API_URL is not set. Aborting."; \
		exit 1; \
	fi
	aws s3 sync frontend/build s3://$$S3_BUCKET_NAME --delete
	aws cloudfront create-invalidation --distribution-id $$DISTRIBUTION_ID --paths "/*" --no-cli-pager
	
.PHONY: build-docker
build-docker:
	docker build . -t shopscraper:local

.PHONY: run-api
run-api: start-dependencies
	go run ./cmd/api/main.go

.PHONY: run-mailer
run-mailer: start-dependencies
	go run ./cmd/mailer/main.go

.PHONY: run-scraper
run-scraper: start-dependencies
	go run ./cmd/scraper/main.go
	
.PHONY: run-frontend
run-frontend:
	cd frontend && npm start

.PHONY: clean
clean:
	docker-compose down
	go clean
	rm -f $(API_BIN) $(MAILER_BIN) $(SCRAPER_BIN)
	rm -rf frontend/build

.PHONY: clean-all
clean-all:
	docker-compose down
	go clean
	rm -f $(API_BIN) $(MAILER_BIN) $(SCRAPER_BIN)
	rm -rf frontend/build
	rm -rf frontend/node_modules
	
.PHONY: release
release:
	$(eval CURRENT_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`))
	$(eval MAJOR=$(shell echo $(CURRENT_VERSION) | sed 's/v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\)/\1/'))
	$(eval MINOR=$(shell echo $(CURRENT_VERSION) | sed 's/v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\)/\2/'))
	$(eval PATCH=$(shell echo $(CURRENT_VERSION) | sed 's/v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\)/\3/'))
	$(eval NEW_PATCH=$(shell echo $$(( $(PATCH) + 1 ))))
	$(eval NEW_VERSION=v$(MAJOR).$(MINOR).$(NEW_PATCH))
	@echo "Current version: $(CURRENT_VERSION)"
	@echo "New version: $(NEW_VERSION)"
	git tag $(NEW_VERSION)
	git push origin $(NEW_VERSION)

.PHONY: release-minor
release-minor:
	$(eval CURRENT_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`))
	$(eval MAJOR=$(shell echo $(CURRENT_VERSION) | sed 's/v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\)/\1/'))
	$(eval MINOR=$(shell echo $(CURRENT_VERSION) | sed 's/v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\)/\2/'))
	$(eval NEW_MINOR=$(shell echo $$(( $(MINOR) + 1 ))))
	$(eval NEW_VERSION=v$(MAJOR).$(NEW_MINOR).0)
	@echo "Current version: $(CURRENT_VERSION)"
	@echo "New version: $(NEW_VERSION)"
	git tag $(NEW_VERSION)
	git push origin $(NEW_VERSION)

.PHONY: release-major
release-major:
	$(eval CURRENT_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`))
	$(eval MAJOR=$(shell echo $(CURRENT_VERSION) | sed 's/v\([0-9]*\)\.\([0-9]*\)\.\([0-9]*\)/\1/'))
	$(eval NEW_MAJOR=$(shell echo $$(( $(MAJOR) + 1 ))))
	$(eval NEW_VERSION=v$(NEW_MAJOR).0.0)
	@echo "Current version: $(CURRENT_VERSION)"
	@echo "New version: $(NEW_VERSION)"
	git tag $(NEW_VERSION)
	git push origin $(NEW_VERSION)