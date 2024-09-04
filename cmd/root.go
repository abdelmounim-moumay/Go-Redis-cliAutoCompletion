package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
)

// Types pour la configuration (comme vous les avez définis)
type IP struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
	Role string `json:"role,omitempty"`
}

type Site struct {
	T1T7             int  `json:"t1t7"`
	SepID            int  `json:"sepId"`
	DownStreamPrefix int  `json:"downStreamPrefix"`
	IPs              []IP `json:"ips"`
}

type Codec struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type AllowedSourceIP struct {
	Name     string  `json:"name"`
	Monitor  bool    `json:"monitor"`
	Sites    []Site  `json:"sites"`
	Codecs   []Codec `json:"codecs,omitempty"`
	PeerName string  `json:"peer_name"`
}

type Config struct {
	VlanID           int               `json:"vlan_id"`
	SipProfile       int               `json:"sip_profile"`
	AllowedSourceIPs []AllowedSourceIP `json:"allowed_source_ips"`
}

var redisClient *redis.Client
var config Config

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

// Charge le fichier de configuration JSON
func loadConfig(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &config)
	return err
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

	// Ajout de getCmd et loadCmd
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(loadCmd) 
}

func initConfig() {
	err := loadConfig("config/config.json")
	if err != nil {
		fmt.Println("Error loading config file:", err)
		os.Exit(1)
	}
}
