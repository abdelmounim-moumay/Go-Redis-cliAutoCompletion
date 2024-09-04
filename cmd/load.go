package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load [file]",
	Short: "Load configuration from a JSON file into Redis",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalf("Error reading file: %v\n", err)
		}

		var config Config
		err = json.Unmarshal(data, &config)
		if err != nil {
			log.Fatalf("Error unmarshalling JSON: %v\n", err)
		}

		ctx := context.Background()

		for _, sourceIP := range config.AllowedSourceIPs {
			key := sourceIP.Name
			value, err := json.Marshal(sourceIP)
			if err != nil {
				log.Fatalf("Error marshalling sourceIP: %v\n", err)
			}

			// Stocker l'objet dans Redis
			err = redisClient.Set(ctx, key, value, 0).Err()
			if err != nil {
				log.Fatalf("Error setting key '%s' in Redis: %v\n", key, err)
			}
			fmt.Printf("Loaded '%s' into Redis\n", key)
		}

		fmt.Println("Configuration loaded successfully into Redis.")
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
}
