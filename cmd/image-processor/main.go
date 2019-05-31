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
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

const (
	defaultUnixSocket = "unix:///var/run/docker.sock"
	defaultAPIVersion = "1.18"
	dockerSocketPath  = "/var/run/docker.sock"
)

var level = flag.String("loglevel", "info", "default log level: debug, info, warn, error, fatal, panic")
var dockerUsername = flag.String("username", "", "Registry username for pulling and pushing")
var dockerPassword = flag.String("password", "", "Registry password for pulling and pushing")
var registryHost = flag.String("registry-host", "", "Hostname of the registry being monitored")
var repository = flag.String("repository", "", "Repository on the registry to pull and push")
var interval = flag.String("run-every", "2m", "The refresh interval in minutes")

var (
	dockerClient *docker.Client
	dockerHost   string
)

func main() {
	// Parse the command line flags.
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	lvl, err := log.ParseLevel(*level)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	log.SetLevel(lvl)

	// Ensure we have proper values.
	// if *dockerUsername == "" {
	// 	log.Fatalln("Missing username flag")
	// }

	// if *dockerPassword == "" {
	// 	log.Fatalln("Missing password flag")
	// }

	// if *registryHost == "" {
	// 	log.Fatalln("Missing registry-host flag")
	// }

	// if *repository == "" {
	// 	log.Fatalln("Missing repository flag")
	// }

	// client, err := NewDockerClient()
	// if err != nil {
	// 	log.Fatalf("error getting in the Docker client due to %+v", err)
	// }

	// images, err := client.ListImages(docker.ListImagesOptions{All: true})
	// if err != nil {
	// 	log.Fatalf("unable to list Docker images due to %+v", err)
	// }

	// log.Infof("images: %+v", images)

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Printf("%s %+v\n", image.ID[:10], image.RepoDigests)
	}
}

func base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func encodeAuthHeader(username string, password string) string {
	data := fmt.Sprintf("{ \"username\": \"%s\", \"password\": \"%s\" }", username, password)
	encoded := base64Encode(data)
	return encoded
}
