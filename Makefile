.PHONY: test coverage coverage-html

test:
	go test ./... -coverprofile=cover

coverage: test
	go tool cover -func=cover

coverage-html: test
	go tool cover -html=cover
