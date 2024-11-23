CURDIR = $(shell pwd)
GENERATEDDIR = ${CURDIR}/generated
PROTODIR = ${CURDIR}/proto

auth-generate-proto:
	rm -rf ${GENERATEDDIR}/auth
	mkdir -p ${GENERATEDDIR}
	protoc --proto_path=${PROTODIR} \
	--go_out=${GENERATEDDIR} \
	--go-grpc_out=${GENERATEDDIR} \
	${PROTODIR}/auth/auth.proto 
	go mod tidy