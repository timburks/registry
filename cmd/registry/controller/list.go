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

package controller

import (
	"context"
	"github.com/apigee/registry/cmd/registry/core"
	"github.com/apigee/registry/connection"
	"github.com/apigee/registry/server/names"
)

func ListResources(ctx context.Context, client connection.Client, pattern, filter string) ([]Resource, error) {
	var result []Resource
	var err error

	// First try to match collection names.
	if m := names.ApisRegexp().FindStringSubmatch(pattern); m != nil {
		err = core.ListAPIs(ctx, client, m, filter, GenerateApiHandler(&result))
	} else if m := names.SpecsRegexp().FindStringSubmatch(pattern); m != nil {
		err = core.ListSpecs(ctx, client, m, filter, GenerateSpecHandler(&result))
	} else if m := names.ArtifactsRegexp().FindStringSubmatch(pattern); m != nil {
		err = core.ListArtifacts(ctx, client, m, filter, false, GenerateArtifactHandler(&result))
	}

	// Then try to match resource names.
	if m := names.ApiRegexp().FindStringSubmatch(pattern); m != nil {
		err = core.ListAPIs(ctx, client, m, filter, GenerateApiHandler(&result))
	} else if m := names.SpecRegexp().FindStringSubmatch(pattern); m != nil {
		err = core.ListSpecs(ctx, client, m, filter, GenerateSpecHandler(&result))
	} else if m := names.ArtifactRegexp().FindStringSubmatch(pattern); m != nil {
		err = core.ListArtifacts(ctx, client, m, filter, false, GenerateArtifactHandler(&result))
	}

	if err != nil {
		return nil, err
	}

	return result, err
}
