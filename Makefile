TARGET=meshbird

all: clean build

clean:
	rm -rf $(TARGET)

depends:
	go get -v

build:
	go build -v -o  $(TARGET) *.go

fmt:
	go fmt *.go
