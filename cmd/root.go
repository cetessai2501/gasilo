package cmd

import (
        "fmt"
        "log"
        "github.com/spf13/cobra"
        "os"
)

var developer string

var rootCmd = &cobra.Command{
        Use:   "my_flags",
        Short: "A brief description of your application",
        Long:  `A longer description.`,
       
}

func Execute() {
        if err := rootCmd.Execute(); err != nil {
                fmt.Println(err)
                os.Exit(1)
        }
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show heft.io client version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(RootCmd.Use + " " + VERSION)
	},
}


func init() {
        cobra.OnInitialize(initConfig)
        rootCmd.PersistentFlags().StringVar(&developer, "developer", "Unknown Developer!", "Developer name.")
        RootCmd.AddCommand(versionCmd)
        
}

func initConfig() {
        developer, _ := rootCmd.Flags().GetString("developer")
        log.Printf("eee", developer)       
        if developer != "" {
                fmt.Println("Developer:", developer)
        }
}
