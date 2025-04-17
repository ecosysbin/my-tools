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

import "github.com/gin-gonic/gin"

type GCPContextHandlerFunc func(*GCPContext)

func GCPContextWrapper(h GCPContextHandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cc := &GCPContext{Context: ctx}
		h(cc)
	}
}

// GCPContext is a wrapper around gin.Context, adding some custom methods.
// If needed, new methods can continue to be added here.
type GCPContext struct {
	*gin.Context
}

func (gcpc *GCPContext) SetUesrName(name string) {
	gcpc.Set("username", name)
}

func (gcpc *GCPContext) GetUesrName() string {
	return gcpc.GetString("username")
}
