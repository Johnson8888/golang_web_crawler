
GOCMD=GO111MODULE=on go


build_mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o target/web_crawler_mac main.go


build_windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o target/web_crawler_win.exe main.go


build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o target/web_crawler_linux main.go






clean:
	rm -rf target/
