package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func isVerbose(cmd *cobra.Command) bool {
	if cmd == nil {
		return viper.GetBool("verbose")
	}
	if cmd.Flag("verbose") != nil && cmd.Flag("verbose").Changed {
		return true
	}
	return viper.GetBool("verbose")
}
