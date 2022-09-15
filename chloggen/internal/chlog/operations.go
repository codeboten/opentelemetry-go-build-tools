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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"go.opentelemetry.io/build-tools/chloggen/internal/entry"
	"gopkg.in/yaml.v2"
)

type summary struct {
	Version         string
	BreakingChanges []string
	Deprecations    []string
	NewComponents   []string
	Enhancements    []string
	BugFixes        []string
}

func ReadEntries(ctx Context) ([]*entry.Entry, error) {
	entryYAMLs, err := filepath.Glob(filepath.Join(ctx.UnreleasedDir, "*.yaml"))
	if err != nil {
		return nil, err
	}

	entries := make([]*entry.Entry, 0, len(entryYAMLs))
	for _, entryYAML := range entryYAMLs {
		if filepath.Base(entryYAML) == filepath.Base(ctx.TemplateYAML) {
			continue
		}

		fileBytes, err := os.ReadFile(entryYAML)
		if err != nil {
			return nil, err
		}

		entry := &entry.Entry{}
		if err = yaml.Unmarshal(fileBytes, entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func DeleteEntries(ctx Context) error {
	entryYAMLs, err := filepath.Glob(filepath.Join(ctx.UnreleasedDir, "*.yaml"))
	if err != nil {
		return err
	}

	for _, entryYAML := range entryYAMLs {
		if filepath.Base(entryYAML) == filepath.Base(ctx.TemplateYAML) {
			continue
		}

		if err := os.Remove(entryYAML); err != nil {
			fmt.Printf("Failed to delete: %s\n", entryYAML)
		}
	}
	return nil
}

func GenerateSummary(version string, entries []*entry.Entry) (string, error) {
	s := summary{
		Version: version,
	}

	for _, e := range entries {
		switch e.ChangeType {
		case entry.Breaking:
			s.BreakingChanges = append(s.BreakingChanges, e.String())
		case entry.Deprecation:
			s.Deprecations = append(s.Deprecations, e.String())
		case entry.NewComponent:
			s.NewComponents = append(s.NewComponents, e.String())
		case entry.Enhancement:
			s.Enhancements = append(s.Enhancements, e.String())
		case entry.BugFix:
			s.BugFixes = append(s.BugFixes, e.String())
		}
	}

	s.BreakingChanges = sort.StringSlice(s.BreakingChanges)
	s.Deprecations = sort.StringSlice(s.Deprecations)
	s.NewComponents = sort.StringSlice(s.NewComponents)
	s.Enhancements = sort.StringSlice(s.Enhancements)
	s.BugFixes = sort.StringSlice(s.BugFixes)

	return s.String()
}

func (s summary) String() (string, error) {
	summaryTmpl := filepath.Join(moduleDir(), "summary.tmpl")

	tmpl := template.Must(
		template.
			New("summary.tmpl").
			Option("missingkey=error").
			ParseFiles(summaryTmpl))

	buf := bytes.Buffer{}
	if err := tmpl.Execute(&buf, s); err != nil {
		return "", fmt.Errorf("failed executing template: %w", err)
	}

	return buf.String(), nil
}
