package main

import (
	"fmt"
	"os"

	"github.com/ranmerc/gophercises/task/cmd"
	"github.com/ranmerc/gophercises/task/db"
)

func main() {
	if err := db.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to open DB connection: %v", err)
		os.Exit(1)
	}

	defer db.Close()

	cmd.Execute()
}
