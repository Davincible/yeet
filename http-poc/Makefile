
.PHONY: proto
proto:
	@protoc \
		--proto_path=. \
		--go_out=:. \
                --go-grpc_out=. \
		proto/echo.proto
