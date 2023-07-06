.PHONY: ensure_deps build test lint push_init_chat watch

default: test lint

ensure_deps:
	go mod vendor
	go mod tidy

build:
	@ echo "🛠  Start building 🛠"
	@ go build -a \
			 -o "bin/bot" "./internal"
	@ echo "🛠  Build done 🛠"

test:
	go test -v -mod=vendor ./...

lint:
	@ echo "🪛  Start linting 🪛"
	@ golangci-lint run ./...
	@ echo "🪛  Lint done 🪛"

lint_fix:
	@ echo "🪛  Start linting with Fix 🪛"
	@ golangci-lint run --fix  ./...
	@ echo "🪛  Fixed your lint 🪛"

push_init_chat:
	@ echo "Start pushing data in MONGO DB"
	@ go run ./scripts/init-data.go

watch:
	@ echo "👀  Continue working... I'm watching 👀"
	@ gow -c run ./internal