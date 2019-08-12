package cmd

import (
      
        "log"
        "github.com/spf13/cobra"
  
        
       
)


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


func init() {
        rootCmd.AddCommand(serverCmd)
        cmd.Execute()
           
}



