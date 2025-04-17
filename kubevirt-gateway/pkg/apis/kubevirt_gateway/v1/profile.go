//
// Copyright 2023 The Zetyun.GCP Authors.
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
//

package v1

// This is used to define some structures for request data,
// such as structures for handling user requests and response structures,
// or request and response structures for accessing other web services.
//
// Such as:
//
// // User requests.
// type UserStartJupyterServerRequest struct {
// 	   ServerName string `uri:"name"`
// 	   Profile    string `json:"profile"`
// }
//
// // Requests to Jupyterhub server.
// type CreateJupyterUserRequest struct {
// 	   UserNames []string `json:"usernames"`
// 	   Admin     bool     `json:"admin"`
// }

// Example Type
type Profiles []Profile

type Profile struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}
