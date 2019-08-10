package cmd

import (
        "fmt"

        "github.com/spf13/cobra"
)

var countCmd = &cobra.Command{
        Use:   "count",
        Short: "A brief description of your command",
        Long:  `A longer description of count command.`,
        Run: func(cmd *cobra.Command, args []string) {
                fmt.Println("count called")
                number, _ := cmd.Flags().GetInt("number")
                for i := 0; i < number; i++ {
                        fmt.Print(i, " ")
                }
                fmt.Println()

                developer, _ := rootCmd.Flags().GetString("developer")
                if developer != "" {
                        fmt.Println("From count command - Developer:", developer)
                }
        },
}

func init() {
        rootCmd.AddCommand(countCmd)
        countCmd.Flags().Int("number", 10, "A help for number")
}
