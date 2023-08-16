# TLZ;


[![test](https://github.com/zcubbs/tlz/actions/workflows/test.yml/badge.svg)](https://github.com/zcubbs/tlz/actions/workflows/test.yaml)
[![release](https://github.com/zcubbs/tlz/actions/workflows/release.yml/badge.svg)](https://github.com/zcubbs/tlz/actions/workflows/release.yaml)
[![vulnerability-scan](https://github.com/zcubbs/tlz/actions/workflows/scan.yml/badge.svg)](https://github.com/zcubbs/tlz/actions/workflows/scan.yaml)


This application is designed to monitor the status and health of SSL/TLS certificates. In today's digital world, ensuring the validity and security of your SSL/TLS certificates is crucial. With our application, you can streamline this process, keeping an eye on your certificates and their expiry dates.

**TLZ;** does more than just monitor - it actively notifies you of any pending expiry dates, ensuring you have ample time to renew and prevent any potential downtime or security issues. By proactively managing your SSL/TLS certificates, our application helps maintain the integrity and reliability of your digital infrastructure.

Whether you are managing a single website or a vast network, this application is designed to be your reliable companion in maintaining your SSL/TLS certificate health.

## CRON Jobs

### Cron job configuration

### About Patterns

Each job can be configured with a cron pattern. The pattern is a string of 5 or 6 fields separated by white space that
represents a set of times, normally as a schedule to execute some routine. The fields are as follows:

```
*    *    *    *    *    *
-    -    -    -    -    -
|    |    |    |    |    |
|    |    |    |    |    + year [optional]
|    |    |    |    +----- day of week (0 - 7) (Sunday=0 or 7)
|    |    |    +---------- month (1 - 12)
|    |    +--------------- day of month (1 - 31)
|    +-------------------- hour (0 - 23)
+------------------------- min (0 - 59)
```

> You can choose to run a job only once by passing "-" as the cron pattern.
> Exemple: `TLZ_CRON_CHECK_CERTIFICATE_VALIDITY_CRON_PATTERN=-`

## Database Support

This application supports both **PostgreSQL** and **SQLite** databases, adapting to various use cases and deployment scenarios.

### PostgreSQL

PostgreSQL is a powerful, open-source object-relational database system. To configure the application for PostgreSQL, set `database.postgres.enabled` to `true` and provide the necessary connection details.

### SQLite

SQLite is a self-contained, serverless, zero-configuration SQL database engine, ideal for simpler setup requirements. To configure the application for SQLite, set `database.sqlite.enabled` to `true` and provide the necessary connection details.


## Configuration

### Environment Variables Configuration Documentation

| Environment Variable                           | Description                                 | Default Value                  |
|------------------------------------------------|---------------------------------------------|--------------------------------|
| `DEBUG`                                        | Application debug mode                      | `false`                        |
| `HTTP_SERVER_PORT`                             | HTTP server port                            | `8000`                         |
| `HTTP_SERVER_ALLOW_ORIGINS`                    | Allowed origins for CORS                    | `*`                            |
| `HTTP_SERVER_ALLOW_HEADERS`                    | Allowed headers for CORS                    | `Origin, Content-Type, Accept` |
| `HTTP_SERVER_TZ`                               | Time zone for the HTTP server               | `UTC`                          |
| `HTTP_SERVER_ENABLE_PRINT_ROUTES`              | Enable route printing                       | `false`                        |
| `DATABASE_POSTGRES_ENABLED`                    | Enable Postgres database                    | `false`                        |
| `DATABASE_POSTGRES_HOST`                       | Postgres host                               | `""`                           |
| `DATABASE_POSTGRES_PORT`                       | Postgres port                               | `""`                           |
| `DATABASE_POSTGRES_USERNAME`                   | Postgres username                           | `""`                           |
| `DATABASE_POSTGRES_PASSWORD`                   | Postgres password                           | `""`                           |
| `DATABASE_POSTGRES_DATABASE`                   | Postgres database name                      | `""`                           |
| `DATABASE_POSTGRES_SSL_MODE`                   | SSL mode for Postgres                       | `""`                           |
| `DATABASE_POSTGRES_VERBOSE`                    | Verbose mode for Postgres                   | `""`                           |
| `DATABASE_SQLITE_ENABLED`                      | Enable SQLite database                      | `true`                         |
| `DATABASE_SQLITE_DB_NAME`                      | SQLite database name                        | `"DefaultDbName"`              |
| `CRON_CHECK_CERTIFICATE_VALIDITY_ENABLED`      | Enable certificate validity check cron job  | `true`                         |
| `CRON_CHECK_CERTIFICATE_VALIDITY_CRON_PATTERN` | Cron pattern for certificate validity check | `"*/10 * * * * *"`             |
| `CRON_SEND_MAIL_NOTIFICATION_ENABLED`          | Enable mail notification cron job           | `true`                         |
| `CRON_SEND_MAIL_NOTIFICATION_CRON_PATTERN`     | Cron pattern for mail notification          | `"*/10 * * * * *"`             |
| `NOTIFICATION_MAIL_SMTP_ENABLED`               | Enable SMTP for mail notification           | `true`                         |
| `NOTIFICATION_MAIL_SMTP_HOST`                  | SMTP host for mail notification             | `"localhost"`                  |
| `NOTIFICATION_MAIL_SMTP_PORT`                  | SMTP port for mail notification             | `1025`                         |
| `NOTIFICATION_MAIL_SMTP_USERNAME`              | SMTP username for mail notification         | `""`                           |
| `NOTIFICATION_MAIL_SMTP_PASSWORD`              | SMTP password for mail notification         | `""`                           |
| `NOTIFICATION_MAIL_SMTP_FROM`                  | SMTP from address for mail notification     | `"no-reply@tlz"`               |

### YAML Configuration Documentation

| YAML Path                                      | Description                                 | Default Value                  |
|------------------------------------------------|---------------------------------------------|--------------------------------|
| `debug`                                        | Application debug mode                      | `false`                        |
| `http_server.port`                             | HTTP server port                            | `8000`                         |
| `http_server.allow_origins`                    | Allowed origins for CORS                    | `*`                            |
| `http_server.allow_headers`                    | Allowed headers for CORS                    | `Origin, Content-Type, Accept` |
| `http_server.tz`                               | Time zone for the HTTP server               | `UTC`                          |
| `http_server.enable_print_routes`              | Enable route printing                       | `false`                        |
| `database.postgres.enabled`                    | Enable Postgres database                    | `false`                        |
| `database.postgres.host`                       | Postgres host                               | `""`                           |
| `database.postgres.port`                       | Postgres port                               | `""`                           |
| `database.postgres.username`                   | Postgres username                           | `""`                           |
| `database.postgres.password`                   | Postgres password                           | `""`                           |
| `database.postgres.database`                   | Postgres database name                      | `""`                           |
| `database.postgres.ssl_mode`                   | SSL mode for Postgres                       | `""`                           |
| `database.postgres.verbose`                    | Verbose mode for Postgres                   | `""`                           |
| `database.sqlite.enabled`                      | Enable SQLite database                      | `true`                         |
| `database.sqlite.db_name`                      | SQLite database name                        | `"DefaultDbName"`              |
| `cron.check_certificate_validity.enabled`      | Enable certificate validity check cron job  | `true`                         |
| `cron.check_certificate_validity.cron_pattern` | Cron pattern for certificate validity check | `"*/10 * * * * *"`             |
| `cron.send_mail_notification.enabled`          | Enable mail notification cron job           | `true`                         |
| `cron.send_mail_notification.cron_pattern`     | Cron pattern for mail notification          | `"*/10 * * * * *"`             |
| `notification.mail.smtp.enabled`               | Enable SMTP for mail notification           | `true`                         |
| `notification.mail.smtp.host`                  | SMTP host for mail notification             | `"localhost"`                  |
| `notification.mail.smtp.port`                  | SMTP port for mail notification             | `1025`                         |
| `notification.mail.smtp.username`              | SMTP username for mail notification         | `""`                           |
| `notification.mail.smtp.password`              | SMTP password for mail notification         | `""`                           |
| `notification.mail.smtp.from`                  | SMTP from address for mail notification     | `"no-reply@tlz"`               |

## Development

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.20 or higher)
- [Task](https://taskfile.dev/#/installation)
- [Docker](https://docs.docker.com/get-docker/) 
- Windows:
  - [Scoop](https://scoop.sh/) 
