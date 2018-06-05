#python gen
protoc --python_out=plugins=grpc:. *.proto
#golang gen
protoc --go_out=plugins=grpc:. *.proto;
