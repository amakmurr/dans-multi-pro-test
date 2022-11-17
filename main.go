package main

import (
	"fmt"
	"github.com/amakmurr/dans-multi-pro-test/internal"
	_ "github.com/jackc/pgx/v4/stdlib"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -config oapi_codegen.yaml ./api/openapi3.yml
func main() {
	config, err := internal.NewConfig("config.yml")
	if err != nil {
		fmt.Println(err)
		return
	}

	server, err := internal.NewServer(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = server.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
}
