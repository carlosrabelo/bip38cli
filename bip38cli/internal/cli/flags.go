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
		return true
	}
	if f := cmd.LocalFlags().Lookup("compressed"); f != nil && f.Changed {
		if val, err := cmd.LocalFlags().GetBool("compressed"); err == nil {
			return val
		}
	}
	for c := cmd; c != nil; c = c.Parent() {
		if f := c.PersistentFlags().Lookup("compressed"); f != nil && f.Changed {
			if val, err := c.PersistentFlags().GetBool("compressed"); err == nil {
				return val
			}
		}
	}
	return true
}
