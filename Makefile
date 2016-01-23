TARGET=meshbird

all: clean build

clean:
	rm -rf $(TARGET)

depends:
	go get -u -v

build:
	go build -v -o  $(TARGET) *.go

fmt:
	go fmt *.go
