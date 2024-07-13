service_name="domain"
protoc --proto_path proto/ --go_out ops-backend/${service_name}/api/gen/v1 --go_opt paths=source_relative --go-grpc_out ops-backend/${service_name}/api/gen/v1 --go-grpc_opt paths=source_relative --proto_path proto/${service_name} ${service_name}.proto

protoc --proto_path proto/ --grpc-gateway_out ops-backend/${service_name}/api/gen/v1 --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --proto_path proto/${service_name} ${service_name}.proto

service_name="whereip"
protoc --proto_path proto/ --go_out ops-backend/${service_name}/api/gen/v1 --go_opt paths=source_relative --go-grpc_out ops-backend/${service_name}/api/gen/v1 --go-grpc_opt paths=source_relative --proto_path proto/${service_name} ${service_name}.proto