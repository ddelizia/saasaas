test:
	go test ./... -race

analysis:
	go vet ./...
	go vet -vettool=$(which bodyclose) ./...
	staticcheck ./...

crawl:
	go run pkg/crawler/main.go

build-account:
	docker build --build-arg APP="./cmd/account/http" --target release -t saasaas/account .
	