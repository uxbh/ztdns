// Copyright © 2017 uxbh
// This file is part of github.com/uxbh/ztdns.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ztdns",
	Short: "Zerotier DNS Server",
	Long: `ztDNS is a dedicated DNS server for ZeroTier networks.
This application will serve DNS requests for the members of a ZeroTier
network for both A (IPv4) and AAAA (IPv6) requests`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().Bool("debug", false, "enable debug messages")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ztdns.toml)")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".ztdns") // name of config file (without extension)
	viper.AddConfigPath(".")      // adding current directory as first search path
	viper.AddConfigPath("$HOME")  // adding home directory as second search path

	viper.SetEnvPrefix("ztdns")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
