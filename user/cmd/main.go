package main

import (
	"fmt"
	"net/http"

	"github.com/keslerliv/user/config"
	"github.com/keslerliv/user/internal/routes"
	"github.com/keslerliv/user/pkg/db"
)

func main() {
	config.LoadConfig()

	fmt.Println("started")

	// Start database connection
	conn, err := db.OpenConnection()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Make initial database migration
	db.MakeMigration(conn)

	// Start routes
	r := routes.LoadRoutes()

	// Start server
	http.ListenAndServe(fmt.Sprintf(":%s", config.Env.Port), r)
}
