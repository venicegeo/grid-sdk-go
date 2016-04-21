// Copyright 2016, RadiantBlue Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grid

import (
	"net/http"
	"testing"
)

func TestCheckResponse(t *testing.T) {
	r := http.Response{StatusCode: 200}
	err := CheckResponse(&r)
	if err != nil {
		t.Error(err)
	}

	r = http.Response{StatusCode: 299}
	err = CheckResponse(&r)
	if err != nil {
		t.Error(err)
	}
}

func TestNew(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Error(err)
	}
	// what checks on the grid client
}

func TestLookup(t *testing.T) {
	g, _ := New()
	_, _, err := g.Lookup("")
	if err == nil {
		t.Error("Should have received error")
	}
}

func TestCreateConfigFile(t *testing.T) {
	_, err := CreateConfigFile()
	if err != nil {
		t.Error(err)
	}
	// test that file got created
}

func TestGetConfig(t *testing.T) {
	_, err := GetConfig()
	if err != nil {
		t.Error(err)
	}
	// surely there is more we could test
}
