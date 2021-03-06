// Copyright 2021 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"

	"github.com/apigee/registry/cmd/registry/core"
	"github.com/apigee/registry/connection"
	"github.com/apigee/registry/rpc"
	"github.com/apigee/registry/server/names"
	"github.com/spf13/cobra"
)

type exportCSVRow struct {
	ApiID        string
	VersionID    string
	SpecID       string
	ContentsPath string
}

func init() {
	exportCmd.AddCommand(exportCsvCmd)
	exportCmd.Flags().String("filter", "", "Common Expression Language filter for selecting specs.")
}

var exportCsvCmd = &cobra.Command{
	Use:   "csv [--filter expression] parent ...",
	Short: "Export API specs to a CSV file",
	Args: func(cmd *cobra.Command, args []string) error {
		for _, parent := range args {
			if re := names.VersionRegexp(); !re.MatchString(parent) {
				return fmt.Errorf("invalid parent argument %q: must match %q", parent, re)
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		filter, err := cmd.Flags().GetString("filter")
		if err != nil {
			log.Fatalf("Failed to get filter string from flags: %s", err)
		}

		ctx := context.TODO()
		client, err := connection.NewClient(ctx)
		if err != nil {
			log.Fatalf("Failed to create client: %s", err)
		}

		rows := make([]exportCSVRow, 0)
		for _, parent := range args {
			err := core.ListSpecs(ctx, client, names.VersionRegexp().FindStringSubmatch(parent), filter, func(spec *rpc.ApiSpec) {
				if !names.SpecRegexp().MatchString(spec.GetName()) {
					log.Printf("Failed to parse spec name %q: skipping spec row", spec.GetName())
					return
				}

				m := names.SpecRegexp().FindStringSubmatch(spec.GetName())
				rows = append(rows, exportCSVRow{
					ApiID:        m[2],
					VersionID:    m[3],
					SpecID:       m[4],
					ContentsPath: fmt.Sprintf("$APG_REGISTRY_ADDRESS/%s/contents", spec.GetName()),
				})
			})
			if err != nil {
				log.Fatalf("Failed to list specs: %s", err)
			}
		}

		w := csv.NewWriter(cmd.OutOrStdout())
		if err := w.Write([]string{"api_id", "version_id", "spec_id", "contents_path"}); err != nil {
			log.Fatalf("Failed to write CSV header: %s", err)
		}

		for _, row := range rows {
			if err := w.Write([]string{row.ApiID, row.VersionID, row.SpecID, row.ContentsPath}); err != nil {
				log.Fatalf("Failed to write CSV row %v: %s", row, err)
			}
		}

		w.Flush()
		if err := w.Error(); err != nil {
			log.Fatalf("Error occured while flushing writes to output: %s", err)
		}
	},
}
