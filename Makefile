GO=go

.PHONY: all
all:
	make apiserver
	make agent
	make humcli

apiserver:
	$(GO) build -o bin/apiserver cmd/apiserver/main.go 

agent:
	$(GO) build -o bin/agent cmd/agent/main.go 

humcli:
	$(GO) build -o bin/humcli cmd/humcli/main.go

run-apiserver:
	$(GO) run cmd/apiserver/main.go --listen-address 0.0.0.0

run-agent:
	sudo $(GO) run cmd/agent/main.go --config cmd/agent/config.yaml.sample
