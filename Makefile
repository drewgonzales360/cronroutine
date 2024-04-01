# Builds target for whatever OS this is called from.
build:
	go build

# Runs a script to test basic, happy-path functionality inside the container
test:
	go test -v ./...

# Turns on some hooks to check format and build status before commiting/pushing. Optional, but helpful.
githooks:
	git config --local core.hooksPath .githooks/
