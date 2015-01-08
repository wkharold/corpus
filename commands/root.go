package commands

import (
	"fmt"
	"os"

	"github.com/cznic/ql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wkharold/corpus/catalog"
)

var RootCmd = &cobra.Command{
	Use:   "corpus",
	Short: "An e-mail corpus server",
	Long:  `Loads an e-mail corpus and serves bundles of messages to clients`,
	Run: func(cmd *cobra.Command, args []string) {
		Configure()

		listener := catalog.New(ServerDir)
		_, err := os.Open(ServerDir)
		if err != nil {
			switch err.(*os.PathError).Err.Error() {
			case "no such file or directory":
				if os.Mkdir(ServerDir, os.ModePerm) != nil {
					panic(fmt.Sprintf("Can't create server directory %s [%v]", ServerDir, err))
				}
			default:
				panic(fmt.Sprintf("Unexpected error [%v]", err))
			}
		}

		schema := ql.MustSchema((*catalog.MbxMsg)(nil), "", nil)

		Msgs = make(chan catalog.MbxMsg)
		Done = make(chan chan int)
		go listener(Msgs, Done)
	},
}

var MailDir, ServerDir string

var rootcmd *cobra.Command

const (
	defaultMailDir   = "/corpus/enron_mail_20110402/maildir"
	defaultServerDir = "/var/lib/corpusserver"
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&MailDir, "maildir", "M", "", "directory of mailboxes")
	RootCmd.PersistentFlags().StringVarP(&ServerDir, "serverdir", "S", defaultServerDir, "corpus server directory")
	rootcmd = RootCmd
}

func Configure() {
	viper.SetEnvPrefix("corpus")
	viper.BindEnv("maildir")
	viper.BindEnv("serverdir")

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/corpus.d")
	viper.AddConfigPath("$HOME/.corpus")
	viper.ReadInConfig()

	viper.SetDefault("maildir", "/corpus/enron_mail_20110402/maildir")
	viper.SetDefault("serverdir", defaultServerDir)

	if rootcmd.PersistentFlags().Lookup("maildir").Changed {
		viper.Set("maildir", MailDir)
	}

	if rootcmd.PersistentFlags().Lookup("serverdir").Changed {
		viper.Set("serverdir", ServerDir)
	}
}

func Execute() {
	RootCmd.Execute()
}
