.PHONY: clean

clean:
	rm -f proto/hello/v1/*.go proto/org/project/hello/v1/*.go

proto: clean
	protoc \
		--go_out=plugins=grpc:. \
		--go_opt=paths=source_relative \
		proto/hello/v1/hello.proto \
		proto/org/project/hello/v1/hello.proto
