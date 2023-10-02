package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-faster/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"github.com/iyear/tdl/pkg/utils"
)

func NewGen() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "gen",
		Short:  "A set of gen tools",
		Hidden: true,
	}

	cmd.AddCommand(NewGenDoc())

	return cmd
}

func NewGenDoc() *cobra.Command {
	var dir string

	cmd := &cobra.Command{
		Use:   "doc",
		Short: "Generate doc",
		RunE: func(cmd *cobra.Command, args []string) error {
			const frontmatter = `---
title: "%s"
bookHidden: true
---
`
			cmd.VisitParents(func(c *cobra.Command) {
				// Disable the "Auto generated by spf13/cobra on DATE"
				// as it creates a lot of diffs.
				c.DisableAutoGenTag = true
			})

			if !utils.FS.PathExists(dir) {
				if err := os.MkdirAll(dir, os.ModePerm); err != nil {
					return errors.Wrap(err, "mkdir")
				}
			}

			prepender := func(filename string) string {
				name := filepath.Base(filename)
				base := strings.TrimSuffix(name, path.Ext(name))
				return fmt.Sprintf(frontmatter, strings.Replace(base, "_", " ", -1))
			}

			linkHandler := func(name string) string {
				base := strings.TrimSuffix(name, path.Ext(name))
				return "/docs/more/cli/" + strings.ToLower(base) + "/"
			}

			fmt.Println("Generating command-line documentation in", dir, "...")
			err := doc.GenMarkdownTreeCustom(cmd.Root(), dir, prepender, linkHandler)
			if err != nil {
				return errors.Wrap(err, "gendoc")
			}
			fmt.Println("Done.")

			return nil
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", "", "dir to generate doc")

	_ = cmd.MarkFlagRequired("dir")

	return cmd
}
