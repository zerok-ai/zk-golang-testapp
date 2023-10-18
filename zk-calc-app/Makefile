generate-proto:
	protoc --go_out=. --go-grpc_out=. calculator.proto

grpcurl-list:
	grpcurl -plaintext localhost:50051 list

grpcurl-list-methods:
	grpcurl -plaintext localhost:50051 list calculator.Calculator

grpcurl-add:
	grpcurl -plaintext -d '{"a": 5, "b": 3}' -H "traceparent: 00-0018317d2edc7d4cfa75b39e88c2d570-176069a07e84acaf-01" localhost:50051 calculator.Calculator.Add

grpcurl-other-add:
	grpcurl -plaintext -d '{"a": 5, "b": 3, "forward": "50052"}' -H "traceparent: 00-0018317d2edc7d4cfa75b39e88c2d570-176069a07e84acaf-01" localhost:50051 calculator.Calculator.CallOtherAdd

grpcurl-subtract:
	grpcurl -plaintext -d '{"a": 5, "b": 3}' -H "traceparent: 00-0018317d2edc7d4cfa75b39e88c2d570-176069a07e84acaf-01" localhost:50051 calculator.Calculator.Subtract

grpcurl-other-subtract:
	grpcurl -plaintext -d '{"a": 5, "b": 3, "forward": "50052"}' -H "traceparent: 00-0018317d2edc7d4cfa75b39e88c2d570-176069a07e84acaf-01" localhost:50051 calculator.Calculator.CallOtherSubtract
