package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a value in Redis by key",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]
		ctx := context.Background()
		err := redisClient.Set(ctx, key, value, 0).Err()
		if err != nil {
			fmt.Printf("Error setting key: %v\n", err)
			return
		}
		fmt.Println("Value set successfully")
	},
	ValidArgsFunction: setCmdCompletion,
}

func setCmdCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// List all possible suggestions
	suggestions := []string{"key1", "key2", "key3", "value1", "value2", "value3"}
	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	rootCmd.AddCommand(setCmd)
}
