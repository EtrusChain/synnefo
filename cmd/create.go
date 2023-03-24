package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a todo",
	Long:  `This command will create todo`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("title", "t", "", "specify task title / heading")
	createCmd.Flags().StringP("description", "d", "", "specify task description")
}
