.PHONY: dev dev/tailwind dev/templ dev/sync_assets

dev/templ:
	templ generate --watch --proxy="http://localhost:1323" --proxybind="0.0.0.0" --cmd="go run ./main.go" --open-browser=false -v

dev/tailwind:
	npx --yes @tailwindcss/cli -i ./input.css -o ./public/tailwind.css --watch

dev/sync_assets:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "templ generate" \
	--build.bin "true" \
	--build.delay "500" \
	--build.exclude_dir "" \
	--build.include_dir "assets" \
	--build.include_ext "js,css"

dev:
	make -j3 dev/tailwind dev/templ dev/sync_assets

IMAGE ?= 192.168.1.175:6000/templ
TAG   ?= $(shell git rev-parse --short HEAD)

docker/build:
	docker build -t $(IMAGE):$(TAG) .

docker/push: docker/build
	docker push $(IMAGE):$(TAG)

docker/run:
	docker run --rm -p 1323:1323 $(IMAGE):$(TAG) --insecure-registry