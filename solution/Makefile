GOARCH=amd64
BINARY=minerva
SRC_DIR=src
SRC=main.go config.go state.go websocket.go

all: deps $(shell uname)

Linux:
	cd ${SRC_DIR}; \
	GOOS=linux GOARCH=${GOARCH} go build -o ../${BINARY}-linux-${GOARCH} ${SRC};

Darwin:
	cd ${SRC_DIR}; \
	GOOS=darwin GOARCH=${GOARCH} go build -o ../${BINARY}-darwin-${GOARCH} ${SRC};

deps:
	go get -u github.com/sirupsen/logrus github.com/gorilla/websocket github.com/spf13/viper;

clean:
	rm -f ${BINARY}
