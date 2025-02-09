COMPOSE=docker compose --env-file ./deployment/docker.env  -f ./deployment/docker-compose.yaml
COMPOSE_MONITORING=${COMPOSE} -p container-monitoring

.PHONY: up
up:
	${COMPOSE_MONITORING} up -d

.PHONY: down
down:
	${COMPOSE_MONITORING} down