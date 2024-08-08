package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
)

var redisClient *redis.Client

var rootCmd = &cobra.Command{
	Use:     "rediscli",
	Short:   "A CLI to interact with Redis",
	Version: "0.1",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cobra.OnInitialize(initConfig)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.CompletionOptions.DisableNoDescFlag = true
	rootCmd.CompletionOptions.DisableDescriptions = true

	getCmd.ValidArgsFunction = RedisKeysCompletion
	setCmd.ValidArgsFunction = RedisKeysCompletion
}

func initConfig() {
	// Any additional initialization code
}

func RedisKeysCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ctx := context.Background()
	keys, err := redisClient.Keys(ctx, "*").Result()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	return keys, cobra.ShellCompDirectiveNoFileComp
}
