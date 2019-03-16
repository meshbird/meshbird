build:
	protoc --go_out=../.. protocol.proto
prepare:
	add-apt-repository ppa:maarten-fonville/protobuf
	apt-get updatego install ./vendor/github.com/golang/protobuf/protoc-gen-go/
	apt-get -y install protobuf-compiler
	go get -u github.com/golang/protobuf/protoc-gen-go