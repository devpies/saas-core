include .env

.DEFAULT_GOAL := help

# =============================================================
# ADMIN SERVICE
# =============================================================
admin-test: admin-mock	;@ ## Run admin tests. Add -- -v for verbosity.
	go test $(val) -cover ./internal/admin/...
.PHONY: admin-test

admin-mock: ;@ ## Generate admin mocks.
	go generate ./internal/admin/...
.PHONY: admin-mock

admin-db: ;@ ## Enter admin database.
	@pgcli postgres://postgres:postgres@localhost:$(ADMIN_DB_PORT)/admin
.PHONY: admin-db

admin-db-gen: ;@ ## Generate migration files. Required <name> argument.
	@migrate create -ext sql -dir ./internal/admin/res/migrations -seq $(val)
.PHONY: admin-db-gen

admin-db-migrate: ;@ ## Migrate admin database. Optional <num> argument.
	@migrate -path ./internal/admin/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(ADMIN_DB_PORT)/admin?sslmode=disable up $(val)
.PHONY: admin-db-migrate

admin-db-version: ;@ ## Print migration version for admin database.
	@migrate -path ./internal/admin/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(ADMIN_DB_PORT)/admin?sslmode=disable version
.PHONY: admin-db-version

admin-db-rollback: ;@ ## Rollback admin database. Optional <num> argument.
	@migrate -path ./internal/admin/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(ADMIN_DB_PORT)/admin?sslmode=disable down $(val)
.PHONY: admin-db-rollback

admin-db-force: ;@ ## Force version on admin database. Optional <num> argument.
	@migrate -path ./internal/admin/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(ADMIN_DB_PORT)/admin?sslmode=disable force $(val)
.PHONY: admin-db-force

# =============================================================
# REGISTRATION SERVICE
# =============================================================
registration-test: registration-mock	;@ ## Run registration tests. Add -- -v for verbosity.
	go test $(val) -cover ./internal/registration/...
.PHONY: registration-test

registration-mock: ;@ ## Generate registration mocks.
	go generate ./internal/registration/...
.PHONY: registration-mock

# =============================================================
# TENANT SERVICE
# =============================================================
tenant-test: tenant-mock	;@ ## Run tenant tests. Add -- -v for verbosity.
	go test $(val) -cover ./internal/tenant/...
.PHONY: tenant-test

tenant-mock: ;@ ## Generate tenant mocks.
	go generate ./internal/tenant/...
.PHONY: tenant-mock

# =============================================================
# PROJECT SERVICE
# =============================================================
project-test: project-mock	;@ ## Run project tests. Add -- -v for verbosity.
	go test $(val) -cover ./internal/project/...
.PHONY: project-test

project-mock: ;@ ## Generate project mocks.
	go generate ./internal/project/...
.PHONY: project-mock

project-db: ;@ ## Enter project database.
	@pgcli postgres://postgres:postgres@localhost:$(PROJECT_DB_PORT)/project
.PHONY: project-db

project-db-gen: ;@ ## Generate migration files. Required <name> argument.
	@migrate create -ext sql -dir ./internal/project/res/migrations -seq $(val)
.PHONY: project-db-gen

project-db-migrate: ;@ ## Migrate project database. Optional <num> argument.
	@migrate -path ./internal/project/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(PROJECT_DB_PORT)/project?sslmode=disable up $(val)
.PHONY: project-db-migrate

project-db-version: ;@ ## Print migration version for project database.
	@migrate -path ./internal/project/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(PROJECT_DB_PORT)/project?sslmode=disable version
.PHONY: project-db-version

project-db-rollback: ;@ ## Rollback project database. Optional <num> argument.
	@migrate -path ./internal/project/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(PROJECT_DB_PORT)/project?sslmode=disable down $(val)
.PHONY: project-db-rollback

project-db-force: ;@ ## Force version on project database. Optional <num> argument.
	@migrate -path ./internal/project/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(PROJECT_DB_PORT)/project?sslmode=disable force $(val)
.PHONY: project-db-force

# =============================================================
# USER SERVICE
# =============================================================
user-db: ;@ ## Enter user database.
	@pgcli postgres://postgres:postgres@localhost:$(USER_DB_PORT)/user
.PHONY: user-db

user-db-gen: ;@ ## Generate migration files. Required <name> argument.
	@migrate create -ext sql -dir ./internal/user/res/migrations -seq $(val)
.PHONY: user-db-gen

user-db-migrate: ;@ ## Migrate user database. Optional <num> argument.
	@migrate -path ./internal/user/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(USER_DB_PORT)/user?sslmode=disable up $(val)
.PHONY: user-db-migrate

user-db-version: ;@ ## Print migration version for user database.
	@migrate -path ./internal/user/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(USER_DB_PORT)/user?sslmode=disable version
.PHONY: user-db-version

user-db-rollback: ;@ ## Rollback user database. Optional <num> argument.
	@migrate -path ./internal/user/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(USER_DB_PORT)/user?sslmode=disable down $(val)
.PHONY: user-db-rollback

user-db-force: ;@ ## Force version on user database. Optional <num> argument.
	@migrate -path ./internal/user/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(USER_DB_PORT)/user?sslmode=disable force $(val)
.PHONY: user-db-force

# =============================================================
# SUBSCRIPTION SERVICE
# =============================================================
subscription-test: subscription-test	;@ ## Run subscription tests. Add -- -v for verbosity.
	go test $(val) -cover ./internal/subscription/...
.PHONY: subscription-test

subscription-mock: ;@ ## Generate subscription mocks.
	go generate ./internal/subscription/...
.PHONY: subscription-mock

subscription-db: ;@ ## Enter subscription database.
	@pgcli postgres://postgres:postgres@localhost:$(BILLING_DB_PORT)/subscription
.PHONY: subscription-db

subscription-db-gen: ;@ ## Generate migration files. Required <name> argument.
	@migrate create -ext sql -dir ./internal/subscription/res/migrations -seq $(val)
.PHONY: subscription-db-gen

subscription-db-migrate: ;@ ## Migrate subscription database. Optional <num> argument.
	@migrate -path ./internal/subscription/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(BILLING_DB_PORT)/subscription?sslmode=disable up $(val)
.PHONY: subscription-db-migrate

subscription-db-version: ;@ ## Print migration version for subscription database.
	@migrate -path ./internal/subscription/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(BILLING_DB_PORT)/subscription?sslmode=disable version
.PHONY: subscription-db-version

subscription-db-rollback: ;@ ## Rollback subscription database. Optional <num> argument.
	@migrate -path ./internal/subscription/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(BILLING_DB_PORT)/subscription?sslmode=disable down $(val)
.PHONY: subscription-db-rollback

subscription-db-force: ;@ ## Force version on subscription database. Optional <num> argument.
	@migrate -path ./internal/subscription/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(BILLING_DB_PORT)/subscription?sslmode=disable force $(val)
.PHONY: subscription-db-force

init: ;@ ## Initialize project. Do once.
	@./scripts/setup-data-folder.sh
.PHONY: init

ports: ;@ ## Port forward Traefik ports.
	kubectl port-forward --address 0.0.0.0 service/traefik 8000:8000 8080:8080 443:4443 -n default
.PHONY: ports

routes: ;@ ## Apply ingress routes.
	kubectl apply -f ./manifests/traefik-routes.yaml
.PHONY: routes

nats: ;## Port forward NATS port.
	kubectl port-forward statefulset.apps/nats 4222
.PHONY: nats

lint: ;@ ## Run linter. Optional <package path> argument.
	@golangci-lint run $(val)
.PHONY: lint

test: ;@ ## Run all tests. Add -- -v for verbosity.
	go test $(val) -cover ./...
.PHONY: test

help:
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'
	@echo
	@echo "- Setup Instructions -"
	@echo
	@echo "1. tilt up"
	@echo "2. make ports"
	@echo "3. make routes"
	@echo
.PHONY: help

# http://bit.ly/37TR1r2
# TODO: Find a better way
ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),test lint admin-test admin-db-gen admin-db-migrate admin-db-rollback admin-db-force project-test project-db-gen project-db-migrate project-db-rollback project-db-force user-test user-db-gen user-db-migrate user-db-rollback user-db-force subscription-test subscription-db-gen subscription-db-migrate subscription-db-rollback subscription-db-force registration-test tenant-test))
  val := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(val):;@:)
endif
