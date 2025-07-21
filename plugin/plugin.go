// Copyright 2019 The casbin Authors. All Rights Reserved.
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

package plugin

import (
	"log"
	"net/url"

	"github.com/casbin/casbin/v2"
	"github.com/docker/go-plugins-helpers/authorization"
)

// CasbinAuthZPlugin is the Casbin Authorization Plugin
type CasbinAuthZPlugin struct {
	// Casbin enforcer
	enforcer *casbin.Enforcer
}

// newPlugin creates a new casbin authorization plugin
func NewPlugin(casbinModel string, casbinPolicy string) (*CasbinAuthZPlugin, error) {
	plugin := &CasbinAuthZPlugin{}

	var err error
	plugin.enforcer, err = casbin.NewEnforcer(casbinModel, casbinPolicy)

	return plugin, err
}

// AuthZReq authorizes the docker client command.
// The command is allowed only if it matches a Casbin policy rule.
// Otherwise, the request is denied!
func (plugin *CasbinAuthZPlugin) AuthZReq(req authorization.Request) authorization.Response {
	// Parse request and the request body
	reqURI, _ := url.QueryUnescape(req.RequestURI)
	reqURL, _ := url.ParseRequestURI(reqURI)

	// Check if URL parsing failed
	if reqURL == nil {
		log.Println("Failed to parse request URI:", reqURI)
		return authorization.Response{Allow: false, Msg: "Invalid request URI"}
	}

	obj := reqURL.String()
	act := req.RequestMethod

	allowed, err := plugin.enforcer.Enforce(obj, act)
	if err != nil {
		log.Println("Enforce error:", err)
		return authorization.Response{Allow: false, Msg: "Authorization error: " + err.Error()}
	}

	if allowed {
		log.Println("obj:", obj, ", act:", act, "res: allowed")
		return authorization.Response{Allow: true}
	}

	log.Println("obj:", obj, ", act:", act, "res: denied")
	return authorization.Response{Allow: false, Msg: "Access denied by casbin plugin"}
}

// AuthZRes authorizes the docker client response.
// All responses are allowed by default.
func (plugin *CasbinAuthZPlugin) AuthZRes(req authorization.Request) authorization.Response {
	// Allowed by default.
	return authorization.Response{Allow: true}
}
