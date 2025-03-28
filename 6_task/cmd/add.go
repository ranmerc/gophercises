package cmd

import (
	"fmt"

	"github.com/ranmerc/gophercises/task/db"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <task-message>",
	Args:  cobra.ExactArgs(1),
	Short: "Adds a task",
	RunE: func(cmd *cobra.Command, args []string) error {
		task := args[0]

		sqlStatement := `INSERT INTO todos (description) VALUES($1)`
		if _, err := db.DB.Exec(sqlStatement, task); err != nil {
			return fmt.Errorf("failed to add tasks %w", err)
		}

		text := fmt.Sprintf(`Added "%s" to your task list.`, task)
		fmt.Println(text)

		return nil
	},
}
