package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"otus-highload/db"
	"otus-highload/redis"

	"github.com/spf13/cobra"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Start the Redis consumer",
	Run: func(cmd *cobra.Command, args []string) {
		mode, _ := cmd.Flags().GetString("mode")
		if mode != "feeder" && mode != "notifier" {
			log.Fatalf("Invalid mode: %s. Use 'feeder' or 'notifier'", mode)
		}

		log.Printf("Starting consumer in %s mode...", mode)
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

		ctx, cancel := context.WithCancel(context.Background())
		go handleShutdown(cancel)

		if mode == "feeder" {
			if err := redis.CreateFeederGroup(); err != nil {
				log.Fatalf("Failed to create feeder group: %v", err)
			}
			redis.StartFeedTasksConsumer()
		} else {
			if err := redis.CreateNotificationGroup(); err != nil {
				log.Fatalf("Failed to create notification group: %v", err)
			}
			redis.StartNotificationConsumer()
		}

		<-ctx.Done()
		log.Println("Consumer stopped gracefully.")
	},
}

func init() {
	consumerCmd.Flags().String("mode", "feeder", "Mode of consumer: feeder or notifier")
	rootCmd.AddCommand(consumerCmd)
}

func handleShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Println("Received shutdown signal, stopping...")
	cancel()
}
