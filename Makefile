package := oauth2-cli

.PHONY: $(package)

.PHONY: build
build: $(package)
	go build -ldflags "-X $(package)/cmd.Version=`git describe --tags --always` -X $(package)/cmd.Date=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` -X $(package)/cmd.Commit=`git rev-parse HEAD`" \
	--o bin/$(package) .

