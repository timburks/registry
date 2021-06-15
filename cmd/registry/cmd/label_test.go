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
	"github.com/apigee/registry/connection"
	"github.com/apigee/registry/rpc"
	"github.com/google/go-cmp/cmp"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestLabel(t *testing.T) {
	var err error

	projectId := "label-test"
	projectName := "projects/" + projectId
	apiId := "sample"
	apiName := projectName + "/apis/" + apiId
	versionId := "1.0.0"
	versionName := apiName + "/versions/" + versionId
	specId := "openapi.json"
	specName := versionName + "/specs/" + specId

	// Create a registry client.
	ctx := context.Background()
	registryClient, err := connection.NewClient(ctx)
	if err != nil {
		t.Fatalf("error creating client: %+v", err)
	}
	defer registryClient.Close()
	// Clear the test project.
	err = registryClient.DeleteProject(ctx, &rpc.DeleteProjectRequest{
		Name: projectName,
	})
	if err != nil && status.Code(err) != codes.NotFound {
		t.Fatalf("error deleting test project: %+v", err)
	}
	// Create the test project.
	_, err = registryClient.CreateProject(ctx, &rpc.CreateProjectRequest{
		ProjectId: projectId,
		Project: &rpc.Project{
			DisplayName: "Test",
			Description: "A test catalog",
		},
	})
	if err != nil {
		t.Fatalf("error creating project %s", err)
	}
	// Create a sample api.
	_, err = registryClient.CreateApi(ctx, &rpc.CreateApiRequest{
		Parent: projectName,
		ApiId:  apiId,
		Api:    &rpc.Api{},
	})
	if err != nil {
		t.Fatalf("error creating api %s", err)
	}
	// Create a sample version.
	_, err = registryClient.CreateApiVersion(ctx, &rpc.CreateApiVersionRequest{
		Parent:       apiName,
		ApiVersionId: versionId,
		ApiVersion:   &rpc.ApiVersion{},
	})
	if err != nil {
		t.Fatalf("error creating version %s", err)
	}
	// Create a sample spec.
	_, err = registryClient.CreateApiSpec(ctx, &rpc.CreateApiSpecRequest{
		Parent:    versionName,
		ApiSpecId: specId,
		ApiSpec: &rpc.ApiSpec{
			MimeType: "application/x.openapi;version=3.0.0",
			Contents: []byte(`{"openapi": "3.0.0", "info": {"title": "test", "version": "v1"}, "paths": {}}`),
		},
	})
	if err != nil {
		t.Fatalf("error creating spec %s", err)
	}

	testCases := []struct {
		comment  string
		args     []string
		expected map[string]string
	}{
		{"add some labels",
			[]string{"a=1", "b=2"},
			map[string]string{"a": "1", "b": "2"}},
		{"remove one label and overwrite the other",
			[]string{"a=3", "b-", "--overwrite"},
			map[string]string{"a": "3"}},
		{"changing a label without --overwrite should be ignored",
			[]string{"a=4"},
			map[string]string{"a": "3"}},
	}
	// test labels for APIs.
	for _, tc := range testCases {
		cmd := &cobra.Command{}
		cmd.SetArgs(append([]string{"label", apiName}, tc.args...))
		labelCmd(cmd)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() with args %+v returned error: %s", tc.args, err)
		}
		api, err := registryClient.GetApi(ctx, &rpc.GetApiRequest{
			Name: apiName,
		})
		if err != nil {
			t.Errorf("error getting api %s", err)
		} else {
			if diff := cmp.Diff(api.Labels, tc.expected); diff != "" {
				t.Errorf("labels were incorrectly set %+v", api.Labels)
			}
		}
	}
	// test labels for versions.
	for _, tc := range testCases {
		cmd := &cobra.Command{}
		cmd.SetArgs(append([]string{"label", versionName}, tc.args...))
		labelCmd(cmd)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() with args %+v returned error: %s", tc.args, err)
		}
		version, err := registryClient.GetApiVersion(ctx, &rpc.GetApiVersionRequest{
			Name: versionName,
		})
		if err != nil {
			t.Errorf("error getting version %s", err)
		} else {
			if diff := cmp.Diff(version.Labels, tc.expected); diff != "" {
				t.Errorf("labels were incorrectly set %+v", version.Labels)
			}
		}
	}
	// test labels for specs.
	for _, tc := range testCases {
		cmd := &cobra.Command{}
		cmd.SetArgs(append([]string{"label", specName}, tc.args...))
		labelCmd(cmd)
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Execute() with args %+v returned error: %s", tc.args, err)
		}
		spec, err := registryClient.GetApiSpec(ctx, &rpc.GetApiSpecRequest{
			Name: specName,
		})
		if err != nil {
			t.Errorf("error getting api %s", err)
		} else {
			if diff := cmp.Diff(spec.Labels, tc.expected); diff != "" {
				t.Errorf("labels were incorrectly set %+v", spec.Labels)
			}
		}
	}

	// Delete the test project.
	if false {
		req := &rpc.DeleteProjectRequest{
			Name: projectName,
		}
		err = registryClient.DeleteProject(ctx, req)
		if err != nil {
			t.Fatalf("failed to delete test project: %s", err)
		}
	}
}
