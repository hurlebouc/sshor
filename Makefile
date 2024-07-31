build:
	go build -v .

build-all:
	GOOS=linux GOARCH=amd64 go build -v -o sshor-linux-amd64 .
	GOOS=linux GOARCH=386 go build -v -o sshor-linux-386 .
	GOOS=linux GOARCH=arm go build -v -o sshor-linux-arm .
	GOOS=linux GOARCH=arm64 go build -v -o sshor-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build -v -o sshor-windows-amd64.exe .
	GOOS=windows GOARCH=386 go build -v -o sshor-windows-386.exe .
	GOOS=windows GOARCH=arm go build -v -o sshor-windows-arm .
	GOOS=windows GOARCH=arm64 go build -v -o sshor-windows-arm64 .
	GOOS=darwin GOARCH=amd64 go build -v -o sshor-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -v -o sshor-darwin-arm64 .

sign:
	for file in target/bin/*; do if [[ $$file != *.asc ]] then gpg --detach-sign --armor --batch --yes $${file}; fi; done

deploy-latest: sign
	name=$$(gh release view --json name --template '{{ .name }}') && gh release upload "$$name" target/bin/*

docker-build:
	podman-compose build sshor && container=$$(podman create localhost/sshor:latest) && rm -rf target/bin && podman export "$$container" | tar -xvf -