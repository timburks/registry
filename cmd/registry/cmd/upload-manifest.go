// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"github.com/apigee/registry/cmd/registry/controller"
	"github.com/apigee/registry/cmd/registry/core"
	"github.com/apigee/registry/connection"
	"github.com/apigee/registry/rpc"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"log"
)

func init() {
	uploadCmd.AddCommand(uploadManifestCmd)
	uploadManifestCmd.Flags().String("project_id", "", "ProjectID this manifest should be associated with.")
}

var uploadManifestCmd = &cobra.Command{
	Use:   "manifest FILE_PATH --project_id=value",
	Short: "Upload a dependency manifest",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		flagset := cmd.LocalFlags()
		manifestPath := args[0]
		if manifestPath == "" {
			log.Fatal("Please provide manifest_path")
		}

		project, err := flagset.GetString("project_id")
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
		if project == "" {
			log.Fatal("Please specify a project")
		}

		manifest, err := controller.ReadManifestProto(manifestPath)
		if err != nil {
			log.Fatal(err.Error())
		}
		manifestData, err := proto.Marshal(manifest)

		ctx := context.TODO()
		client, err := connection.NewClient(ctx)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		artifact := &rpc.Artifact{
			Name:     "projects/" + project + "/artifacts/" + manifest.Name,
			MimeType: core.MimeTypeForMessageType("google.cloud.apigee.registry.applications.v1alpha1.Manifest"),
			Contents: manifestData,
		}
		err = core.SetArtifact(ctx, client, artifact)
		if err != nil {
			log.Fatal(err.Error())
		}

		return
	},
}
