package main

import (
	"fmt"
	"os"
)

func main() {
	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	a.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}
