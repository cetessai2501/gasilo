package cmd

import (
        "fmt"
        "log"
        "github.com/spf13/cobra"
        "os"
)

var developer string

var (
	// VERSION is set during build
	VERSION = "0.0.1"
)


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

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the  server",
	RunE:  runServer(),
}


func runServer(){
        app.NewServer()
	app.InitStores()
        app.StartServer()


}




var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show heft.io client version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(rootCmd.Use + " " + VERSION)
	},
}


func init() {
        cobra.OnInitialize(initConfig)
        rootCmd.PersistentFlags().StringVar(&developer, "developer", "Unknown Developer!", "Developer name.")
        rootCmd.AddCommand(versionCmd,serverCmd)
        
}

func initConfig() {
        developer, _ := rootCmd.Flags().GetString("developer")
        log.Printf("eee", developer)       
        if developer != "" {
                fmt.Println("Developer:", developer)
        }
}
