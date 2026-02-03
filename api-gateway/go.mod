module github.com/adityakw90/microservice-sample-app/api-gateway

go 1.25.5

require (
	github.com/adityakw90/service-user-proto v0.0.0
	github.com/gorilla/mux v1.8.1
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11
)

require (
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251029180050-ab9386a59fda // indirect
)

replace github.com/adityakw90/service-user-proto => /media/adit/SSD/project/MTAmedia/repo/service-user/service-user-proto
