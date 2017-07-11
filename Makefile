all: assets/in assets/check assets/out

assets:
	mkdir assets

assets/in: assets in/in.go
	GOARCH=amd64 GOOS=linux go build -o assets/in in/in.go

assets/out: assets out/out.sh
	cp out/out.sh assets/out

assets/check: assets check/check.go
	GOARCH=amd64 GOOS=linux go build -o assets/check check/check.go
