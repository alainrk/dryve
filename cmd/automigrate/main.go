package main

import (
	"dryve/internal/config"
	"dryve/internal/datastruct"
	"dryve/internal/repository"
	"fmt"
)

func main() {
	f := "config.json"

	config := config.NewConfig(f)

	db, err := repository.NewDB(config.Database)
	if err != nil {
		fmt.Printf("database initialization failed with err %v\n", err)
	}

	tables := []any{
		&datastruct.File{},
	}

	err = repository.Automigrate(db, tables)
	if err != nil {
		fmt.Printf("automigration failed with err %v\n", err)
	}
}
