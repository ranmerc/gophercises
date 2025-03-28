package cmd

import (
	"fmt"
	"strconv"

	"github.com/ranmerc/gophercises/task/db"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do <task-number>",
	Args:  cobra.ExactArgs(1),
	Short: "Mark a task done",
	RunE: func(cmd *cobra.Command, args []string) error {
		taskNumber, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task number: %w", err)
		}

		sqlStatement := `
			SELECT id, description from todos 
			ORDER BY time 
			LIMIT 1 OFFSET $1`

		row := db.DB.QueryRow(sqlStatement, taskNumber-1)

		var (
			taskID      string
			description string
		)

		if err := row.Scan(&taskID, &description); err != nil {
			return fmt.Errorf("invalid task number: %w", err)
		}

		sqlStatement = `DELETE FROM todos WHERE id=$1`
		if _, err := db.DB.Exec(sqlStatement, taskID); err != nil {
			return fmt.Errorf("failed to delete: %w", err)
		}

		text := fmt.Sprintf(`You have completed the %q task.`, description)
		fmt.Println(text)

		return nil
	},
}
