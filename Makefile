.PHONY: dev backend frontend

# Default target
dev:
	$(MAKE) -j 2 backend frontend

backend:
	go run cmd/main.go

frontend:
	cd app && yarn dev
