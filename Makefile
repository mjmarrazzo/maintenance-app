.PHONY: dev dev/tailwind dev/templ dev/sync_assets

dev/templ:
	templ generate --watch --proxy="http://localhost:1323" --proxybind="0.0.0.0" --cmd="go run ./cmd/api.go" --open-browser=false -v

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