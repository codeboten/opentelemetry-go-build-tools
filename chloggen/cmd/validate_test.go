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
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/build-tools/chloggen/internal/entry"
)

func TestValidateE2E(t *testing.T) {
	tests := []struct {
		name    string
		entries []*entry.Entry
		wantErr string
	}{
		{
			name:    "all_valid",
			entries: getSampleEntries(),
		},
		{
			name: "invalid_change_type",
			entries: func() []*entry.Entry {
				return append(getSampleEntries(), &entry.Entry{
					ChangeType: "fake",
					Component:  "receiver/foo",
					Note:       "Add some bar",
					Issues:     []int{12345},
				})
			}(),
			wantErr: "'fake' is not a valid 'change_type'",
		},
		{
			name: "missing_component",
			entries: func() []*entry.Entry {
				return append(getSampleEntries(), &entry.Entry{
					ChangeType: entry.BugFix,
					Component:  "",
					Note:       "Add some bar",
					Issues:     []int{12345},
				})
			}(),
			wantErr: "specify a 'component'",
		},
		{
			name: "missing_note",
			entries: func() []*entry.Entry {
				return append(getSampleEntries(), &entry.Entry{
					ChangeType: entry.BugFix,
					Component:  "receiver/foo",
					Note:       "",
					Issues:     []int{12345},
				})
			}(),
			wantErr: "specify a 'note'",
		},
		{
			name: "missing_issue",
			entries: func() []*entry.Entry {
				return append(getSampleEntries(), &entry.Entry{
					ChangeType: entry.BugFix,
					Component:  "receiver/foo",
					Note:       "Add some bar",
					Issues:     []int{},
				})
			}(),
			wantErr: "specify one or more issues #'s",
		},
		{
			name: "all_invalid",
			entries: func() []*entry.Entry {
				sampleEntries := getSampleEntries()
				for _, e := range sampleEntries {
					e.ChangeType = "fake"
				}
				return sampleEntries
			}(),
			wantErr: "'fake' is not a valid 'change_type'",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := setupTestDir(t, tc.entries)
			cmd := validateCmd
			cmd.Flags().Set("changelog", ctx.ChangelogMD)
			// 	"--template", ctx.TemplateYAML,
			// 	"--changes-directory", ctx.UnreleasedDir,
			// })
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			err := cmd.Execute()

			// err := validate(ctx)
			// _, err := ioutil.ReadAll(b)
			if tc.wantErr != "" {
				require.Regexp(t, tc.wantErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
