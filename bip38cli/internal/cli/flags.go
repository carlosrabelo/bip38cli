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

func outputFormat(cmd *cobra.Command) string {
	if cmd == nil {
		return "text"
	}
	if cmd.Flag("output-format") != nil && cmd.Flag("output-format").Changed {
		return cmd.Flag("output-format").Value.String()
	}
	return "text"
}

func isCompressed(cmd *cobra.Command) bool {
	if cmd == nil {
		return true // default to compressed
	}
	if cmd.Flag("compressed") != nil && cmd.Flag("compressed").Changed {
		val, _ := cmd.Flags().GetBool("compressed")
		return val
	}
	return true // default to compressed
}
