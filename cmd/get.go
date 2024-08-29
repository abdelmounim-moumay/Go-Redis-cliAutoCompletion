package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get details by name from Redis",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		ctx := context.Background()

		// Récupérer l'objet JSON stocké dans Redis
		jsonData, err := redisClient.Get(ctx, name).Result()
		if err != nil {
			fmt.Printf("Error getting details for name '%s': %v\n", name, err)
			return
		}

		var sourceIP AllowedSourceIP
		err = json.Unmarshal([]byte(jsonData), &sourceIP)
		if err != nil {
			fmt.Printf("Error unmarshalling JSON data: %v\n", err)
			return
		}
		fmt.Printf("Details for '%s':\n", name)
		fmt.Printf("  Monitor: %t\n", sourceIP.Monitor)
		fmt.Printf("  Peer Name: %s\n", sourceIP.PeerName)

		for _, site := range sourceIP.Sites {
			fmt.Printf("  Site: T1T7=%d, SepID=%d, DownStreamPrefix=%d\n", site.T1T7, site.SepID, site.DownStreamPrefix)
			for _, ip := range site.IPs {
				fmt.Printf("    IP: %s, Port: %d, Role: %s\n", ip.IP, ip.Port, ip.Role)
			}
		}

		for _, codec := range sourceIP.Codecs {
			fmt.Printf("  Codec: Name=%s, ID=%d\n", codec.Name, codec.ID)
		}
	},
	ValidArgsFunction: completeNames,
}

// Fonction pour compléter les noms disponibles dans Redis
func completeNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ctx := context.Background()

	keys, err := redisClient.Keys(ctx, "*"+toComplete+"*").Result()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	return keys, cobra.ShellCompDirectiveNoFileComp
}

func init() {

	rootCmd.AddCommand(getCmd)
}
