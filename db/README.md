This project is using [atlas](https://atlasgo.io/) as the database setup/migration tool.

# Installing Atlas
Installation instructions can be found on the [Getting Started](https://atlasgo.io/getting-started/) page there.
Alternatively, you can run the below commands using the atlas docker container image by
`alias atlas="docker run -v $PWD/migrations:/migrations --network=host arigaio/atlas"` (which works as long as you are running your 'atlas' command from the `db/` directory).

# Making Schema changes
## Creating new migration steps
Details can be found in the Atlas documentation linked above, but in brief, if you need to make schema changes:
- run `make dev` to start the atlas-migration postgres container

- modify `schema.hcl` with your desired updates
- create the new migration .sql file by running `db/add_migration.sh <migration_name>`
  where `migration_name` is the name you wish to use for the new migration step.

See the Atlas docs linked above for the [atlas schema](https://atlasgo.io/atlas-schema/sql-resources) and more details and options.

### Clear Data In Database
If you'd like to clear all the data in your local Postgres to start from a clean slate, you can do so by running the following command _from the root of the project_
```bash
$ docker compose down && docker volume remove order-service_postgres-vol && docker compose up
```

## Testing new migration steps
Automated tests will apply all of the migration steps, but if you want to apply the migrations to your local db for development you can use
`db/apply_local.sh`

Alternatively, you can run atlas apply manually by doing
```bash
export POSTGRESQL_URL=postgres://postgres:example@localhost:5432/order_service?sslmode=disable
atlas migrate apply -u ${POSTGRESQL_URL}
```
