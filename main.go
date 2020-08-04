package main

import (
	"github.com/huangc28/migrator/cmd"
)

// command to create migration file boilerplate
//
// Commands:
//
//   - Create migrations
//
//     `go run migrator.go create be29_create_user_table`
//
//     The above command will generate 2 files:
//
//     {{ PROCESS_WORKING_DIRECTORY }}/cmd/migrations/be29_create_user_table_202002161130/up.go
//     {{ PROCESS_WORKING_DIRECTORY }}/cmd/migrations/be29_create_user_table_202002161130/down.go
//
//     If `migration` directory isn't found in the directory where `go` command is being run, `migration` directory will be created. Migration files would be placed in migration directory.
//
//   - Run all outstanding migrations
//
//     `go run migrator.go migrate`
//
//     Tables that have been migrated are stored in DB so that we can check any outstanding migrations to be run.
//

func main() {
	cmd.Execute()
}
