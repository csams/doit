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

	logincmd "github.com/csams/doit/cmd/login"
	"github.com/csams/doit/cmd/migrate"
	"github.com/csams/doit/cmd/serve"

	"github.com/csams/doit/pkg/login"
	"github.com/csams/doit/pkg/server"
	"github.com/csams/doit/pkg/storage"
)

var (
	rootCmd = &cobra.Command{
		Use: "doit",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initConfig()
		},
	}

	rootLog = logrusr.New(logrus.New())

	options = struct {
		Storage *storage.Options `mapstructure:"storage"`
		Login   *login.Options   `mapstructure:"login"`
		Server  *server.Options  `mapstructure:"serve"`
	}{
		storage.NewOptions(),
		login.NewOptions(),
		server.NewOptions(),
	}
)

func init() {
	viper.SetEnvPrefix("doit")
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.config/doit/config.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	serveCmd := serve.NewCommand(rootLog.WithName("serve"), options.Storage, options.Server)
	rootCmd.AddCommand(serveCmd)
	viper.BindPFlags(serveCmd.Flags())

	migrateCmd := migrate.NewCommand(rootLog.WithName("migrate"), options.Storage)
	rootCmd.AddCommand(migrateCmd)
	viper.BindPFlags(migrateCmd.Flags())

	loginCmd := logincmd.NewCommand(rootLog.WithName("login"), options.Login)
	rootCmd.AddCommand(loginCmd)
	viper.BindPFlags(loginCmd.Flags())
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
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(cfgDir)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	viper.Unmarshal(&options)
}

// Execute runs the root command
func Execute() {
	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		rootLog.V(0).Info("Unhandled", "error", err)
		os.Exit(1)
	}
}
