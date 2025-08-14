start-deps:
	docker compose up

stop-deps:
	docker compose down

setup-ssh-tunnel:
	./tools/setup_ssh_tunnel.sh

run-local:
	./tools/run_local.sh

build-docker:
	docker build -t go-template .