TARGET=meshbird

all: clean depends build

clean:
	rm -rf $(TARGET)

depends:
	go get -d

build:
	go build -o $(TARGET) *.go

fmt:
	go fmt *.go
