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
	"github.com/apigee/registry/connection"
	"github.com/apigee/registry/rpc"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"testing"
)

func TestManifestUpload(t *testing.T) {

	tests := []struct {
		desc     string
		project  string
		filePath string
		want     rpc.Manifest
	}{
		{
			desc:     "simple manifest upload",
			project:  "upload-manifest-demo",
			filePath: "../controller/test/manifest_test.yaml",
			want: rpc.Manifest{
				Name: "test-manifest",
				GeneratedResources: []*rpc.GeneratedResource{
					{
						Pattern: "apis/-/versions/-/specs/-/artifacts/lint-gnostic",
						Dependencies: []*rpc.Dependency{
							{
								Pattern: "$resource.spec",
								Filter:  "mime_type.contains('openapi')",
							},
						},
						Action: "compute lint $source0 --linter gnostic",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			ctx := context.Background()
			client, err := connection.NewClient(ctx)
			if err != nil {
				t.Fatalf("Setup: Failed to create client: %s", err)
			}

			err = client.DeleteProject(ctx, &rpc.DeleteProjectRequest{
				Name: "projects/" + test.project,
			})
			if err != nil && status.Code(err) != codes.NotFound {
				t.Fatalf("Setup: Failed to delete test project: %s", err)
			}

			_, err = client.CreateProject(ctx, &rpc.CreateProjectRequest{
				ProjectId: test.project,
				Project: &rpc.Project{
					DisplayName: "Demo",
					Description: "A demo catalog",
				},
			})
			if err != nil {
				t.Fatalf("Failed to create project %s: %s", test.project, err.Error())
			}

			args := append([]string{"upload", "manifest"}, test.filePath, "--project_id", test.project)
			rootCmd.SetArgs(args)
			if err = uploadManifestCmd.Execute(); err != nil {
				t.Fatalf("Execute() with args %v returned error: %s", args, err)
			}

			req := &rpc.GetArtifactContentsRequest{
				Name: "projects/" + test.project + "/artifacts/test-manifest",
			}

			manifest := rpc.Manifest{}
			body, err := client.GetArtifactContents(ctx, req)
			contents := body.GetData()
			err = proto.Unmarshal(contents, &manifest)

			// Verify the manifest definition is correct
			opts := cmp.Options{
				protocmp.Transform(),
			}

			if diff := cmp.Diff(test.want, manifest, opts); diff != "" {
				t.Errorf("GetArtifactContents returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}

}
