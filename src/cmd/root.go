package cmd

import (
	"fmt"
	"go-template/src/core/log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "app",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var configFile string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is config.yaml)")
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath("cfg/")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetDefault("Log.Level", "debug")
	viper.SetDefault("Log.Color", true)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("unable to read config: %v\n", err)
		os.Exit(1)
	}

	//if err := readConfigFromHashicorp(); err != nil {
	//	fmt.Printf("unable to read config from hashicorp: %v\n", err)
	//	os.Exit(1)
	//}
}

//func readConfigFromHashicorp() error {
//	hConf, err := hashicorp.InitConfig()
//	if err != nil {
//		return errors.Wrap(err, "init hashicorp config")
//	}
//
//	logger, err := getLogger()
//	if err != nil {
//		return errors.Wrap(err, "get logger")
//	}
//
//	hClient, err := hashicorp.New(hConf, logger)
//	if err != nil {
//		return errors.Wrap(err, "create hashicorp client")
//	}
//
//	allData := viper.AllSettings()
//	err = hClient.GetDataFromField(&allData)
//	if err != nil {
//		return errors.Wrap(err, "get data from hashicorp")
//	}
//
//	for key, value := range allData {
//		viper.Set(key, value)
//	}
//
//	return nil
//}

func getLogger() (log.Logger, error) {

	configLogger, err := log.InitConfig()
	if err != nil {
		return nil, err
	}

	logger, err := log.NewLogger(configLogger, log.InstanceZapLogger)
	if err != nil {
		return nil, err
	}
	return logger, nil
}
