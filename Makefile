GO_LDFLAGS=-ldflags " -w"

TAG=v0323
PREFIX=barnettzqg/blog
build:
	@echo "🐳 $@"
	@go build -a -installsuffix cgo ${GO_LDFLAGS} .
image: clean
	@echo "🐳 $@"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ${GO_LDFLAGS} .
	@docker build -t $(PREFIX):$(TAG) .
	@docker push $(PREFIX):$(TAG)
	@rm -f journey
clean:
	@rm -f journey	