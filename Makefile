binary=upcoming-shows
release=release/${binary}
ldxflags=\
	-X main.clientID=${SPOTIFY_ID} \
	-X main.clientSecret=${SPOTIFY_SECRET}

debug: check-env
	@go build -ldflags "${ldxflags}"

release: check-env
	@go build -ldflags "-s -w ${ldxflags}" -o ${release}
	upx ${release}

check-env:
ifndef SPOTIFY_ID
	$(error SPOTIFY_ID is undefined, check your environment exports)
endif
ifndef SPOTIFY_SECRET
	$(error SPOTIFY_SECRET is undefined, check your environment exports)
endif

.PHONY: release debug check-env
