package main

import (
	"github.com/keisuke-m123/godiagramgen/cmd/godiagramgen/classdiagram"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "godiagramgen",
		Short: "Godiagramgen is a golang class and package diagram generator",
	}
	root.AddCommand(
		classdiagram.NewClassDiagramGenCommand(),
	)
	if err := root.Execute(); err != nil {
		panic(err)
	}
}
