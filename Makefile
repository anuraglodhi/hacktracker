build:
	@go1.22.3 build -o bin/hacktracker

run: build
	@./bin/hacktracker

test:
	@go1.22.3 test ./...