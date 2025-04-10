package cmd

import (
	"log"
	"net/http"
	"os"

	"otus-highload/db"
	"otus-highload/redis"
	"otus-highload/router"

	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting web server...")

		redis.InitRedis()
		db.InitDB()
		defer func() {
			if db.MasterDB != nil {
				db.MasterDB.Close()
			}
			if db.SlaveDB != nil {
				db.SlaveDB.Close()
			}
			log.Println("Database connections closed.")
		}()
		go redis.StartNotificationConsumer()
		r := router.NewRouter()
		port := os.Getenv("GO_PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Listening on port %s", port)
		log.Fatal(http.ListenAndServe(":"+port, r))
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
}
