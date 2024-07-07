package main

import (
    "fmt"
    "os"

    "github.com/bridge-yt/helm-bridge-plugin/pkg/resources"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    log "github.com/sirupsen/logrus"
)

var (
    apiURL string
)

var registerCmd = &cobra.Command{
    Use:   "register",
    Short: "Register resources with the Bridge service",
    Run: func(cmd *cobra.Command, args []string) {
        resources.Register(apiURL)
    },
}

var translateCmd = &cobra.Command{
    Use:   "translate",
    Short: "Translate placeholders in Helm templates",
    Run: func(cmd *cobra.Command, args []string) {
        resources.Translate(apiURL)
    },
}

var deployCmd = &cobra.Command{
    Use:   "deploy",
    Short: "Translate placeholders and register resources",
    Run: func(cmd *cobra.Command, args []string) {
        resources.Translate(apiURL)
        resources.Register(apiURL)
    },
}

func main() {
    rootCmd := &cobra.Command{Use: "bridge"}
    rootCmd.AddCommand(registerCmd)
    rootCmd.AddCommand(translateCmd)
    rootCmd.AddCommand(deployCmd)
    registerCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "API URL for Bridge service (required)")
    translateCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "API URL for Bridge service (required)")
    deployCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "API URL for Bridge service (required)")
    viper.BindPFlag("api_url", registerCmd.PersistentFlags().Lookup("api-url"))
    viper.BindPFlag("api_url", translateCmd.PersistentFlags().Lookup("api-url"))
    viper.BindPFlag("api_url", deployCmd.PersistentFlags().Lookup("api-url"))

    cobra.OnInitialize(initConfig)

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func initConfig() {
    configFile := viper.GetString("config")
    if configFile != "" {
        viper.SetConfigFile(configFile)
    } else {
        viper.AddConfigPath(".")
        viper.SetConfigName("config")
    }

    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err == nil {
        log.Infof("Using config file: %s", viper.ConfigFileUsed())
    } else {
        log.Warnf("No config file found. Using default settings")
    }

    apiURL = viper.GetString("api_url")
}
