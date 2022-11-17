package cmd

import (
	"context"
	"fmt"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/csams/doit/cmd/add"
	"github.com/csams/doit/cmd/modify"
	"github.com/csams/doit/cmd/remove"
	"github.com/csams/doit/cmd/search"
	"github.com/csams/doit/cmd/start"
	"github.com/csams/doit/cmd/stop"
	store "github.com/csams/doit/pkg/storage/factory"
	_ "github.com/csams/doit/pkg/storage/file"
)

var (
	rootCmd = &cobra.Command{
		Use: "doit",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			initConfig()

			ctx := cmd.Context()
			st, err := store.New(viper.GetViper())
			if err != nil {
				return err
			}
			cmd.SetContext(context.WithValue(ctx, store.ContextKey, st))
			return nil
		},
	}
)

func init() {
	viper.SetEnvPrefix("doit")
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.config/doit/config.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	rootCmd.AddCommand(add.NewCommand())
	rootCmd.AddCommand(remove.NewCommand())
	rootCmd.AddCommand(modify.NewCommand())
	rootCmd.AddCommand(start.NewCommand())
	rootCmd.AddCommand(stop.NewCommand())
	rootCmd.AddCommand(search.NewCommand("search"))
	rootCmd.AddCommand(search.NewCommand("list"))
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
		os.Exit(1)
	}
}
