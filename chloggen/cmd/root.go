// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/spf13/cobra"
	"go.opentelemetry.io/build-tools/chloggen/internal/chlog"
)

var (
	changesDirectory string
	changelog        string
	template         string
	ctx              chlog.Context
)

var rootCmd = &cobra.Command{
	Use:   "chloggen",
	Short: "Updates CHANGELOG.MD to include all new changes",
	Long:  `chloggen is a tool used to automate the generation of CHANGELOG files using individual yaml files as the source.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ctx = chlog.New(chlog.RepoRoot(), changelog, changesDirectory, template)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().StringVarP(&changesDirectory, "changes-directory", "", "unreleased", "directory containing change files")
	rootCmd.PersistentFlags().StringVarP(&changelog, "changelog", "", "CHANGELOG.md", "destination changelog file")
	rootCmd.PersistentFlags().StringVarP(&template, "template", "", "TEMPLATE.yaml", "template file for new changes")

	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(validateCmd)
}
