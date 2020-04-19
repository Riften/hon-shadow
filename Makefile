.PHONY:shadow
.PHONY:protos
shadow:
	go build -o bin/shadow github.com/Riften/hon-shadow
protos:
	cd pb/protos; protoc --go_out=../. *.proto
