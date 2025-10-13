package cli

import (
	"github.com/spf13/cobra"
)

func isVerbose(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}
	if cmd.Flag("verbose") != nil && cmd.Flag("verbose").Changed {
		return true
	}
	return false
}
