run:
	go run ./cmd/pmdb

tw:
	tailwindcss -i ./static/input.css -o ./static/styles.css -w -m
