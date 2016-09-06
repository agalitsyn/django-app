// Copyright 2016 Gravitational Inc.
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

package rigging

import (
	"encoding/json"

	"github.com/gravitational/trace"
)

func NodesMatchingLabel(labelQuery string) (*NodeList, error) {
	var nodes NodeList

	cmd := KubeCommand("get", "nodes", "-l", labelQuery, "-o", "json")
	out, err := cmd.Output()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	err = json.Unmarshal(out, &nodes)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return &nodes, nil
}

func LabelNode(name string, label string) ([]byte, error) {
	cmd := KubeCommand("label", "node", name, label)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return out, trace.Wrap(err)
	}
	return out, nil
}
