VERSION:=0.1

all: windows linux_amd64 release

linux_amd64:
	mkdir -p linux
	GOOS="linux" GOARCH="amd64" go build
	mv xcup linux/

windows:
	mkdir windows
	GOOS="windows" GOARCH="amd64" go build
	mv xcup.exe windows/

release:
	zip -r xcup-${VERSION}.zip linux/ windows/

clean:
	@rm -rf linux
	@rm -rf windows
	@rm -f xcup*.zip
