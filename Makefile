swagger:
	swagger serve -F swagger --no-open ./swagger.json

swagger-generate:
	swagger generate spec -o ./swagger.json

fmt:
	go fmt ./...

run: fmt
	MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" \
	MONGO_DATABASE=demo \
	go run main.go

docker-start:
	docker start vigorous_roentgen # name of the mongo container on my local machine.

docker-mongo-install:
	docker run --hostname=74fec6376b1f \
	--env=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin \
	--env=GOSU_VERSION=1.17 \
	--env=JSYAML_VERSION=3.13.1 \
	--env=JSYAML_CHECKSUM=662e32319bdd378e91f67578e56a34954b0a2e33aca11d70ab9f4826af24b941 \
	--env=MONGO_PACKAGE=mongodb-org \
	--env=MONGO_REPO=repo.mongodb.org \
	--env=MONGO_MAJOR=8.0 \
	--env=MONGO_VERSION=8.0.3 \
	--env=HOME=/data/db \
	--volume=/data/configdb \
	--volume=/data/db \
	--network=bridge \
	--restart=no \
	--label='desktop.docker.io/wsl-distro=Ubuntu-20.04' \
	--label='org.opencontainers.image.ref.name=ubuntu' \
	--label='org.opencontainers.image.version=24.04' \
	--runtime=runc \
	--name vigorous_roentgen \
	-e MONGO_INITDB_ROOT_USERNAME=admin \
	-e MONGO_INITDB_ROOT_PASSWORD=password \
	-p 27017:27017 -d mongo:latest