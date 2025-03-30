package cmd

import (
	"log"

	"otus-highload/db"
	"otus-highload/redis"
	"otus-highload/utils"

	"github.com/spf13/cobra"
)

var cacheKeeperCmd = &cobra.Command{
	Use:   "feeder",
	Short: "Start the feed updater",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Starting feeder...")

		redis.InitRedis()
		db.InitDB()
		defer func() {
			if db.MasterDB != nil {
				db.MasterDB.Close()
			}
			for _, slaveDB := range db.SlaveDBs {
				if slaveDB != nil {
					slaveDB.Close()
				}
			}
			log.Println("Database connections closed.")
		}()
		if err := redis.CreateFeederGroup(); err != nil {
			log.Fatalf("Failed to create feeder group: %v", err)
		}

		utils.EnqueueTasksForAllUsers()
		redis.StartFeedTasksConsumer()
	},
}

func init() {
	rootCmd.AddCommand(cacheKeeperCmd)
}
