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

package chlog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	changelogMD   = "CHANGELOG.md"
	unreleasedDir = "unreleased"
	templateYAML  = "TEMPLATE.yaml"
)

// chlogContext enables tests by allowing them to work in an test directory
type Context struct {
	rootDir       string
	ChangelogMD   string
	UnreleasedDir string
	TemplateYAML  string
}

func New(root, changelog, changesDir, template string) Context {
	return Context{
		rootDir:       root,
		ChangelogMD:   changelog,
		UnreleasedDir: changesDir,
		TemplateYAML:  template,
	}
}

func NewContext(rootDir string) Context {
	return Context{
		rootDir:       rootDir,
		ChangelogMD:   filepath.Join(rootDir, changelogMD),
		UnreleasedDir: filepath.Join(rootDir, unreleasedDir),
		TemplateYAML:  filepath.Join(rootDir, unreleasedDir, templateYAML),
	}
}

func (c *Context) Unreleased(unreleasedDir string) {
	c.UnreleasedDir = filepath.Join(c.rootDir, unreleasedDir)
	c.TemplateYAML = filepath.Join(c.rootDir, unreleasedDir, templateYAML)
}

var DefaultCtx = NewContext(RepoRoot())

func RepoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("FAIL: Could not determine current working directory")
	}
	return dir
}

func moduleDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("FAIL: Could not determine module directory")
	}
	return filepath.Dir(filename)
}
