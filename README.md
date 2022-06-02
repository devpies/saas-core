# SaaS-Core

This project is a part of "AWS SaaS app in 30 days" - _Proof of Concept_

### Requirements
- go v1.18 or higher
- [tilt](https://tilt.dev/)
- docker-desktop - enable kubernetes
- aws account
- mac/linux environment
- [pgcli](https://www.pgcli.com/)
- [golangci-lint](https://github.com/golangci/golangci-lint)
- [go-migrate](https://github.com/golang-migrate/migrate)
- [saas-infra resources](https://github.com/devpies/saas-infra/tree/main/local/saas)

## Getting Started
Print a description of each supported makefile command.

```bash
> make
admin             Run admin app with live reload.
db-admin-gen      Generate migration files. Required <name> argument.
db-admin          Enter admin database.
db-admin-migrate  Migrate admin database. Optional <num> argument.
db-admin-version  Print migration version for admin database.
db-admin-rollback Rollback admin database. Optional <num> argument.
lint              Run linter.
...
```


## Environment Variables

The `.env` file contains variables for all programs. Using `make` automatically references these values.
Program requirements are also documented in help text. 
```bash
> go run ./cmd/admin -h
Usage: admin [options] [arguments]

OPTIONS
  --web-debug-port/$ADMIN_WEB_DEBUG_PORT                            <string>    (default: 6060)
  --web-production/$ADMIN_WEB_PRODUCTION                            <bool>      (default: false)
  --web-read-timeout/$ADMIN_WEB_READ_TIMEOUT                        <duration>  (default: 5s)
  --web-write-timeout/$ADMIN_WEB_WRITE_TIMEOUT                      <duration>  (default: 5s)
  --web-shutdown-timeout/$ADMIN_WEB_SHUTDOWN_TIMEOUT                <duration>  (default: 5s)
  --web-address/$ADMIN_WEB_ADDRESS                                  <string>    (default: localhost)
  --web-port/$ADMIN_WEB_PORT                                        <string>    (default: 4001)
  --cognito-app-client-id/$ADMIN_COGNITO_APP_CLIENT_ID              <string>    (required)
  --cognito-user-pool-client-id/$ADMIN_COGNITO_USER_POOL_CLIENT_ID  <string>    (required)
  --cognito-region/$ADMIN_COGNITO_REGION                            <string>    (default: eu-central-1)
  --postgres-user/$ADMIN_POSTGRES_USER                              <string>    (required)
  --postgres-password/$ADMIN_POSTGRES_PASSWORD                      <string>    (required)
  --postgres-host/$ADMIN_POSTGRES_HOST                              <string>    (required)
  --postgres-port/$ADMIN_POSTGRES_PORT                              <int>       (required)
  --postgres-db/$ADMIN_POSTGRES_DB                                  <string>    (required)
  --postgres-disable-tls/$ADMIN_POSTGRES_DISABLE_TLS                <bool>      (default: false)
  --help/-h                                                         
  display this help message
```

> __TIP__  
> 
> 1. Using `make` is the easiest way to get started. However, if you choose to run go binaries directly, you can export the `.env` file variables to avoid using CLI flags:  
> ```bash
> export $(grep -v '^#' .env | xargs)
> go run ./cmd/{PROGRAM}
>```
> 
> 2. Enable bash-completion of the makefile targets. Add this in your `~/.bash_profile` file or `~/.bashrc` file.
> ```bash
> complete -W "\`grep -oE '^[a-zA-Z0-9_.-]+:([^=]|$)' ?akefile | sed 's/[^a-zA-Z0-9_.-]*$//'\`" make
> ```

