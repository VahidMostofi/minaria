check_install:
	which swagger || go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	swagger generate spec -o ./swagger.yml --scan-models

swagger_gen_client: swagger
	swagger generate client -f ./swagger.yml -t sdk