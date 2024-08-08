package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a value from Redis by key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		ctx := context.Background()
		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			fmt.Printf("Error getting key: %v\n", err)
			return
		}
		fmt.Printf("Value: %s\n", val)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
