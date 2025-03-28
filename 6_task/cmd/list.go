package cmd

import (
	"fmt"

	"github.com/ranmerc/gophercises/task/db"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Args:  cobra.ExactArgs(0),
	Short: "Lists all the tasks not yet done",
	RunE: func(cmd *cobra.Command, args []string) error {
		sqlStatement := `SELECT description FROM "todos" ORDER BY time`

		rows, err := db.DB.Query(sqlStatement)
		if err != nil {
			return fmt.Errorf("failed to read: %w", err)
		}

		defer rows.Close()

		todos := []string{}

		for rows.Next() {
			todo := ""

			if err := rows.Scan(&todo); err != nil {
				return fmt.Errorf("failed to read: %w", err)
			}

			todos = append(todos, todo)
		}

		if len(todos) > 0 {
			fmt.Println(`You have the following tasks:`)
			for i, todo := range todos {
				line := fmt.Sprintf("%d. %s", i+1, todo)
				fmt.Println(line)
			}
		} else {
			fmt.Println("You are all caught up!")
		}

		return nil
	},
}
