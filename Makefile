
run-disk-agent:
	go run cmd/disk-agent/main.go --config cmd/disk-agent/config.yaml

run-api:
	go run cmd/api/main.go

.PHONY: run-disk-agent run-api