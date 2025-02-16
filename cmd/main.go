package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Dev-AustinPeter/spamhaus-take-home-task/cmd/api"
	"github.com/Dev-AustinPeter/spamhaus-take-home-task/config"
	"github.com/Dev-AustinPeter/spamhaus-take-home-task/constants"
	"github.com/Dev-AustinPeter/spamhaus-take-home-task/service"
	"github.com/Dev-AustinPeter/spamhaus-take-home-task/utils"
)

func main() {

	// Load stored data on startup
	utils.LoadData(constants.DATA_FILE)
	// Start background processes
	go utils.StartBatchSave(constants.DATA_FILE)
	go service.StartBackgroundFetch()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("[INFO] Shutting down, saving data...")
		utils.SaveData(constants.DATA_FILE)
		os.Exit(0)
	}()

	server := api.NewAPIServer(":" + config.Envs.Port)
	if err := server.Run(); err != nil {
		log.Fatalf("[ERROR] Server exited with error: %v", err)
	}
}
