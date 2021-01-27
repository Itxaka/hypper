package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	helmAction "helm.sh/helm/v3/pkg/action"
)

var globalUsage = `Usage: hypper command

A package manager built on Helm charts and Helm itself.
`

func newRootCmd(actionConfig *helmAction.Configuration, out io.Writer, args []string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:          "hypper",
		Short:        "A package manager built on Helm charts and Helm itself",
		Long:         globalUsage,
		SilenceUsage: false,
	}

	flags := cmd.PersistentFlags()
	settings.AddFlags(flags)

	cmd.AddCommand(
		newInstallCmd(actionConfig, out),
	)
	err := flags.Parse(args)
	if err != nil {
		_ = log.Output(2, fmt.Sprintf("failed while parsing flags for %s", args))
		os.Exit(1)
	}

	if settings.NoColors {
		color.NoColor = true // disable colorized output
	}
	return cmd, nil
}
