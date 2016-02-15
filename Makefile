TARGET=meshbird

all: clean build

clean:
	rm -rf $(TARGET)

depends:
	go get -u -v

build:
	go build -v -ldflags="-X main.Version=`cat VERSION`" -o $(TARGET) *.go

fmt:
	go fmt *.go

xc:
	go get github.com/laher/goxc
	goxc -d dist -os="linux,darwin,freebsd" -include 'LICENSE,VERSION' -pv `cat VERSION` -build-ldflags="-X main.Version=`cat VERSION`" xc copy-resources archive-tar-gz deb downloads-page

sign:
	sudo dpkg-sig --sign builder dist/`cat VERSION`/*deb

deploy:
	sudo mkdir -p /var/www/dist/`cat VERSION`
	mv dist/`cat VERSION`/*deb /tmp/
	sudo reprepro -b /var/www/debian --confdir /root/conf includedeb wheezy /tmp/*deb
	sudo reprepro -b /var/www/debian --confdir /root/conf includedeb jessie /tmp/*deb
	sudo reprepro -b /var/www/debian --confdir /root/conf includedeb squeeze /tmp/*deb
	sudo reprepro -b /var/www/debian --confdir /root/conf includedeb precise /tmp/*deb
	sudo reprepro -b /var/www/debian --confdir /root/conf includedeb trusty /tmp/*deb
	sudo reprepro -b /var/www/debian --confdir /root/conf includedeb vivid /tmp/*deb
	sudo reprepro -b /var/www/debian --confdir /root/conf includedeb wily /tmp/*deb
	rm /tmp/*.deb
	sudo cp -rf dist/`cat VERSION`/{darwin,freebsd}* /var/www/dist/`cat VERSION`/
	sudo cp dist/`cat VERSION`/*.tar.gz /var/www/dist/`cat VERSION`/
	sudo chown -R www-data:www-data /var/www/{debian,dist}

xcupload:
	gsutil -m cp -a public-read -r dist/ gs://meshbird.com/dist
