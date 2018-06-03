// Package master is the central component, which (for the moment) creates the commands and queries.
// It redirects and supervises everything.
package main

import "os"

func main() {
	app := Server{}

	app.Initialize()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	app.Run(":" + port)
}
