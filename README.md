# go-swn

**SWN** - **S**overeign **W**eb **N**ode. Universal unit to distribute events and actions among the network.

It was developed as core part of [Neonyx Ecosystem](https://neonyx.io)

## Background
```WARN``` this repository and all related are under active development and should be considered as v0.0.1-beta.

[specs](https://github.com/neonyxhub/swn-specs) - technical specification of SWN processes

[docs](https://github.com/neonyxhub/sws-docs/tree/main/swn) - documentation about SWN usage

## Contribution

### Pre-requisites

* docker
* go >= 1.21.1
* protoc
* protoc-gen-go, protoc-gen-go-grpc mockgen
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/golang/mock/mockgen@latest
```
* [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)

### Development
- Change in Makefile `PROJECT=` variable to the current project name
- `make init`
	- this will run `make configure`, `make add-commit-hook` targets
- Do changes in code
- `go run cmd/swn/swn.go`
- `go run cmd/swn/swn.go -debug=true` -- check all flags via `-h`

### Run as Docker container via Makefile
- `make build-dev` -- build dev image
- `make dev` -- run container of dev image
- You can add necessary `docker run` arguments to it as:
    - `make dev "ARGS=-p 50051:50051"`. It will expose internal port
    - `make dev "CMD_ARGS=-debug=true"`

### Git commit
* Git commit messages should follow the pattern `<subject>: <description>`
* Check [deployment/commit-msg-hook](deployment/commit-msg-hook) for details

### Tests

Run unit tests: `make test`

### Linter

Run linter: `make lint`

## Deployment

### release

- Increment `SUBLEVEL=` in Makefile [TBD]
- `make build-release`
