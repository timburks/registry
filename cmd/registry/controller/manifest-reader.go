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
	"fmt"
	"github.com/apigee/registry/rpc"
	yaml2 "github.com/ghodss/yaml"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Dependency struct {
	Source string `yaml:"source"`
	Filter string `yaml:"filter"`
}

type ManifestEntry struct {
	Resource     string       `yaml:"resource"`
	Filter       string       `yaml:"filter"`
	Dependencies []Dependency `yaml:"dependencies"`
	Action       string       `yaml:"action"`
}

type Manifest struct {
	Project string          `yaml:"project"`
	Entries []ManifestEntry `yaml:"manifest"`
}

// TODO: Remove this function while cleaning up the controller commands
func ReadManifest(filename string) (*Manifest, error) {

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	m := &Manifest{}
	err = yaml.Unmarshal(buf, m)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return m, nil
}

// TODO: Add validation for pattern values and actions
func ReadManifestProto(filename string) (*rpc.Manifest, error) {

	yamlBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := yaml2.YAMLToJSON(yamlBytes)
	m := &rpc.Manifest{}
	err = protojson.Unmarshal(jsonBytes, m)

	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return m, nil
}
