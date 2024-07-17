run:
	go run ./cmd/api

generate_proto:
	@if "$(name)" == "" ( \
		echo Error: name parameter is required & \
		exit 1 \
	)

	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./$(name)/$(name).proto 