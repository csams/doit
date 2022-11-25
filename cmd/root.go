package cmd

import (
	"context"
	"fmt"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bombsimon/logrusr/v3"
	"github.com/sirupsen/logrus"

	"github.com/csams/doit/cmd/migrate"
	"github.com/csams/doit/cmd/serve"
)

var (
	rootCmd = &cobra.Command{
		Use: "doit",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initConfig()
		},
	}

	rootLog = logrusr.New(logrus.New())
)

func init() {
	viper.SetEnvPrefix("doit")
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.config/doit/config.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.AddCommand(serve.NewCommand(rootLog.WithName("serve")))
	rootCmd.AddCommand(migrate.NewCommand(rootLog.WithName("migrate")))
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	cfgFile := viper.GetString("config")

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfgDir := path.Join(home, ".config", "doit")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(cfgDir)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
	viper.SetDefault("storage.config.path", path.Join(cfgDir, "data.json"))
}

func Execute() {
	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		rootLog.V(0).Info("Unhandled", "error", err)
		os.Exit(1)
	}
}
