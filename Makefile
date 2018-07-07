binary=upcoming-shows
release=release/${binary}
ldxflags=\
	-X main.clientID=${SPOTIFY_ID} \
	-X main.clientSecret=${SPOTIFY_SECRET}

debug:
	@go build -ldflags "${ldxflags}"

release:
	@go build -ldflags "-s -w ${ldxflags}" -o ${release}
	upx ${release}

.PHONY: release debug
