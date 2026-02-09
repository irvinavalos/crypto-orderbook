TARGET := prog

build:
	@go build -o ./bin/$(TARGET) ./main.go
	@chmod +x ./bin/$(TARGET)

run: build
	@./bin/$(TARGET)
