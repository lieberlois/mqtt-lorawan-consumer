.PHONY: build clean

build:
	@echo "Downloading packages..."
	@go mod download
	@echo "Building..."
	@GOOS=windows GOARCH=amd64 go build -o build/mqtt_lorawan_consumer.exe ./main.go
	@go build -o build/mqtt_lorawan_consumer ./main.go
	@echo "Done!"

clean:
	@echo "Cleanup..."
	@rm -rf build/*
	@rm -rf docs/public