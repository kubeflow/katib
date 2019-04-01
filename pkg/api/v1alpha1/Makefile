api.pb.go: api.proto
	protoc -I. api.proto --go_out=plugins=grpc:.
health.pb.go: health/health.proto
	protoc -I. health/health.proto --go_out=plugins=grpc:.
