all: clean build

clean:
	rm -rf *.pb.go

build:
	protoc  --go_out=. *.proto
