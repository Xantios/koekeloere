main:
	@rm -rf ./bin/ ; mkdir ./bin
	@go build -o ./bin/koekeloere
	@chmod +x ./bin/koekeloere
