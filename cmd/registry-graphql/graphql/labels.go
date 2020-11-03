// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"errors"

	"github.com/apigee/registry/connection"
	"github.com/apigee/registry/rpc"
	"github.com/graphql-go/graphql"
)

var labelType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Label",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"created": &graphql.Field{
				Type: timestampType,
			},
			"updated": &graphql.Field{
				Type: timestampType,
			},
		},
	},
)

func representationForLabel(label *rpc.Label) map[string]interface{} {
	return map[string]interface{}{
		"id":      label.Name,
		"created": representationForTimestamp(label.CreateTime),
		"updated": representationForTimestamp(label.UpdateTime),
	}
}

func resolveLabels(p graphql.ResolveParams) (interface{}, error) {
	ctx := p.Context
	c, err := connection.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	req := &rpc.ListLabelsRequest{
		Parent: getParentFromParams(p),
	}
	pageToken, isFound := p.Args["after"].(string)
	if isFound {
		req.PageToken = pageToken
	}
	pageSize, isFound := p.Args["first"].(int)
	if isFound {
		req.PageSize = int32(pageSize)
	} else {
		pageSize = 50
	}
	var response *rpc.ListLabelsResponse
	edges := []map[string]interface{}{}
	for len(edges) < pageSize {
		response, err = c.GrpcClient().ListLabels(ctx, req)
		for _, label := range response.GetLabels() {
			edges = append(edges, representationForEdge(representationForLabel(label)))
		}
		req.PageToken = response.GetNextPageToken()
		if req.PageToken == "" {
			break
		}
	}
	return connectionForEdgesAndEndCursor(edges, response.GetNextPageToken()), nil
}

func resolveLabel(p graphql.ResolveParams) (interface{}, error) {
	ctx := p.Context
	c, err := connection.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	name, isFound := p.Args["id"].(string)
	if !isFound {
		return nil, errors.New("missing id field")
	}
	req := &rpc.GetLabelRequest{
		Name: name,
	}
	label, err := c.GetLabel(ctx, req)
	return representationForLabel(label), err
}
