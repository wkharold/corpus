package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdCatalog = &cobra.Command{
	Use:   "catalog",
	Short: "Catalog an e-mail corpus",
	Long:  `Catalog an e-mail corpus so that it can be served to clients`,
	Run:   catalogRun,
}

func init() {
	RootCmd.AddCommand(cmdCatalog)
}

func catalogRun(cmd *cobra.Command, args []string) {
	Configure()
	fmt.Printf("catalogRun [serverdir=%s]\n", viper.GetString("serverdir"))
}
