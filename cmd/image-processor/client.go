/*
Copyright (C) 2019 Synopsys, Inc.
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"os"
	"path"

	docker "github.com/fsouza/go-dockerclient"
)

// NewDockerClient returns the Docker client
func NewDockerClient() (*docker.Client, error) {
	apiVersion := getenv("DOCKER_API_VERSION", defaultAPIVersion)
	dockerHost = os.Getenv("DOCKER_HOST")
	if dockerHost == "" {
		dockerHost = defaultUnixSocket
	}

	if os.Getenv("DOCKER_CERT_PATH") == "" {
		certPath := os.Getenv("DOCKER_CERT_PATH")
		tlsVerify := os.Getenv("DOCKER_TLS_VERIFY") != ""

		if tlsVerify && certPath != "" {
			cert := path.Join(certPath, "cert.pem")
			key := path.Join(certPath, "key.pem")
			ca := path.Join(certPath, "ca.pem")
			return docker.NewVersionedTLSClient(dockerHost, cert, key, ca, apiVersion)
		}
	}
	return docker.NewVersionedClient(dockerHost, apiVersion)
}

func getenv(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultVal
	}
	return val
}