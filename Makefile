
all: clean build 

build: 
	glide create --non-interactive 
	glide install 
	go mod init github.com/wfchiang/davic-micro-core 

clean: 
	-rm -rf vendor
	-rm rm go.mod 
	-rm glide.lock 
	-rm glide.yaml  