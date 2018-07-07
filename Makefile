binary=upcoming-shows

build:
	@go build -ldflags "-s -w -X main.clientID=${SPOTIFY_ID} -X main.clientSecret=${SPOTIFY_SECRET}"
	upx ${binary}
