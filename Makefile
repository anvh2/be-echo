genpb:
	protoc -I/usr/local/include -Igrpc-gen \
		-I$$GOPATH/src \
		-I$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		-I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--plugin=protoc-gen-gogo=$$GOPATH/bin/protoc-gen-gogo \
		--gogo_out=plugins=grpc:grpc-gen \
		idl/echo.proto

