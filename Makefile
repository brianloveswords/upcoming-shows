pkg=github.com/brianloveswords/spotify
secretfile=auth/secret.go
binary=spotify
release=release/${binary}

debug: ${secretfile}
	@go build -o ${binary}

${secretfile}: check-env
	@> ${secretfile} echo package auth
	@>>${secretfile} echo "var clientID = string([]rune"\
		`python -c "print str(list('${SPOTIFY_ID}')).replace('[', '{').replace(']', '}')"` \
	")"
	@>>${secretfile} echo "var clientSecret = string([]rune"\
		`python -c "print str(list('${SPOTIFY_SECRET}')).replace('[', '{').replace(']', '}')"` \
	")"
	@gofmt -w  ${secretfile}

run-debug: debug
	./${binary}

release: check-env
	@go build -ldflags "-s -w" -o ${release}
	upx ${release}

check-env:
ifndef SPOTIFY_ID
	$(error SPOTIFY_ID is undefined, check your environment exports)
endif
ifndef SPOTIFY_SECRET
	$(error SPOTIFY_SECRET is undefined, check your environment exports)
endif

.PHONY: release debug check-env
