protoc --proto_path proto/ --go_out ops-backend/domain/api/gen/v1 --go_opt paths=source_relative --go-grpc_out ops-backend/domain/api/gen/v1 --go-grpc_opt paths=source_relative --proto_path proto/domain domain.proto

protoc --proto_path proto/ --grpc-gateway_out ops-backend/domain/api/gen/v1 --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --proto_path proto/domain domain.proto

protoc --proto_path proto/ --go_out ops-backend/whereip/api/gen/v1 --go_opt paths=source_relative --go-grpc_out ops-backend/whereip/api/gen/v1 --go-grpc_opt paths=source_relative --proto_path proto/whereip whereip.proto