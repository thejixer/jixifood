#!/bin/bash
protoc -I=proto/ --go_out=. --go-grpc_out=. proto/*.proto
echo "Proto files generated successfully!"