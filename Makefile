IMAGE_PREFIX=u03013112
IMAGE_NAME=$(IMAGE_PREFIX)/ss-user
all:
	GOFLAGS=-mod=vendor CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/ss-user -a -installsuffix cgo -ldflags '-w'
	docker build -t $(IMAGE_NAME) .
push:
	docker push $(IMAGE_NAME)
clean:
	docker rmi  $(IMAGE_NAME)
	rm -rf build/*