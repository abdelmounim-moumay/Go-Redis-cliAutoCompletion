package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [name] [subparam]",
	Short: "Configure and save parameters based on the JSON file",
	Args:  cobra.MinimumNArgs(0), // Permet de ne pas exiger d'argument
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// Si aucun argument n'est fourni, affiche la liste des noms disponibles
			fmt.Println("Liste des noms disponibles:")
			for _, allowedIP := range config.AllowedSourceIPs {
				fmt.Println("- " + strings.TrimSpace(allowedIP.Name))
			}
			return
		}

		name := args[0]
		var selectedIP *AllowedSourceIP

		for _, allowedIP := range config.AllowedSourceIPs {
			if strings.TrimSpace(allowedIP.Name) == name {
				selectedIP = &allowedIP
				break
			}
		}

		if selectedIP == nil {
			// Si le nom n'est pas trouvé, affiche la liste des noms disponibles avant d'afficher l'erreur
			fmt.Println("Name not found. Liste des noms disponibles:")
			for _, allowedIP := range config.AllowedSourceIPs {
				fmt.Println("- " + strings.TrimSpace(allowedIP.Name))
			}
			return
		}

		// Afficher et sélectionner les sous-paramètres (le cas échéant)
		if len(args) > 1 {
			subparam := args[1]

			for _, site := range selectedIP.Sites {
				if fmt.Sprintf("site-%d", site.T1T7) == subparam {
					fmt.Printf("Selected Site: T1T7: %d\n", site.T1T7)
					for _, ip := range site.IPs {
						fmt.Printf("IP: %s, Port: %d, Role: %s\n", ip.IP, ip.Port, ip.Role)
					}
				}
			}

			for _, codec := range selectedIP.Codecs {
				if codec.Name == subparam {
					fmt.Printf("Selected Codec: Name: %s, ID: %d\n", codec.Name, codec.ID)
				}
			}
		} else {
			// Afficher la configuration entière
			fmt.Printf("Configuration for %s:\n", name)
			fmt.Printf("Monitor: %t\n", selectedIP.Monitor)
			fmt.Printf("Peer Name: %s\n", selectedIP.PeerName)
			for _, site := range selectedIP.Sites {
				fmt.Printf("Site T1T7: %d\n", site.T1T7)
				for _, ip := range site.IPs {
					fmt.Printf("IP: %s, Port: %d, Role: %s\n", ip.IP, ip.Port, ip.Role)
				}
			}
			for _, codec := range selectedIP.Codecs {
				fmt.Printf("Codec: %s, ID: %d\n", codec.Name, codec.ID)
			}
		}

		// Convertir le selectedIP en JSON
		jsonData, err := json.Marshal(selectedIP)
		if err != nil {
			fmt.Println("Error marshaling to JSON:", err)
			return
		}

		// Sauvegarder dans Redis
		ctx := context.Background()
		err = redisClient.Set(ctx, name, jsonData, 0).Err()
		if err != nil {
			fmt.Println("Error saving to Redis:", err)
			return
		}

		fmt.Println("Configuration saved to Redis successfully")
	},
	ValidArgsFunction: setCmdCompletion,
}

// Configuration de l'auto-complétion pour les sous-paramètres
func setCmdCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var completions []string
	if len(args) == 0 {
		// Première argument (nom) auto-complétion
		for _, allowedIP := range config.AllowedSourceIPs {
			if strings.HasPrefix(allowedIP.Name, toComplete) {
				completions = append(completions, strings.TrimSpace(allowedIP.Name))
			}
		}
	} else if len(args) == 1 {
		// Deuxième argument (subparam) auto-complétion
		name := args[0]
		var selectedIP *AllowedSourceIP

		for _, allowedIP := range config.AllowedSourceIPs {
			if strings.TrimSpace(allowedIP.Name) == name {
				selectedIP = &allowedIP
				break
			}
		}

		if selectedIP != nil {
			for _, site := range selectedIP.Sites {
				siteID := fmt.Sprintf("site-%d", site.T1T7)
				if strings.HasPrefix(siteID, toComplete) {
					completions = append(completions, siteID)
				}
			}

			for _, codec := range selectedIP.Codecs {
				if strings.HasPrefix(codec.Name, toComplete) {
					completions = append(completions, codec.Name)
				}
			}
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	rootCmd.AddCommand(setCmd)
}
