// Copyright (c) 2017-2018 Samsung Electronics Co., Ltd All Rights Reserved
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
// limitations under the License
//

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	loads "github.com/go-openapi/loads"
	flag "github.com/spf13/pflag"

	"git.tizen.org/tools/boruta/http/client"
	"git.tizen.org/tools/weles/artifacts"
	"git.tizen.org/tools/weles/controller"
	"git.tizen.org/tools/weles/manager"
	"git.tizen.org/tools/weles/parser"
	"git.tizen.org/tools/weles/server"
	"git.tizen.org/tools/weles/server/operations"
)

var (
	borutaAddress            string
	borutaRefreshPeriod      time.Duration
	artifactDBName           string
	artifactDBLocation       string
	artifactDownloadQueueCap int
	activeWorkersCap         int
	notifierChannelCap       int
)

// This file was generated by the swagger tool.
// Make sure to regenerate server only with the Makefile which has --exclude-main option set.

func exitOnErr(ctx string, err error) {
	if err != nil {
		log.Fatal(ctx, err)
	}
}

func main() {

	swaggerSpec, err := loads.Embedded(server.SwaggerJSON, server.FlatSwaggerJSON)
	exitOnErr("failed to load embedded swagger spec", err)

	var srv *server.Server // make sure init is called
	var apiDefaults server.APIDefaults

	flag.Int32Var(&apiDefaults.PageLimit, "page-limit", 0, "Default limit of page size returned by Weles API. If set to 0 pagination will be turned off")
	flag.StringVar(&borutaAddress, "boruta-address", "http://127.0.0.1:8487", "Boruta address. Must contain protocol.")
	flag.DurationVar(&borutaRefreshPeriod, "boruta-refresh-period", 2*time.Second, "Boruta refresh period")
	flag.StringVar(&artifactDBName, "db-file", "weles.db", "name of *.db file. Should be located in --db-location")
	flag.StringVar(&artifactDBLocation, "db-location", "/tmp/weles/", "location of *.db file and place where Weles will store artifacts.")
	//TODO: when cyberdryads or testlab instance will be present, performance tests should be done
	// to set default values of below:
	flag.IntVar(&artifactDownloadQueueCap, "artifact-download-queue-cap", 100, "Capacity of artifact download queue")
	flag.IntVar(&activeWorkersCap, "active-workers-cap", 16, "Maximum number of active workers.")
	flag.IntVar(&notifierChannelCap, "notifier-channel-cap", 100, "Notifier channel capacity.")

	//TODO: input validation

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Usage:\n")
		fmt.Fprint(os.Stderr, "  weles-server [OPTIONS]\n\n")

		title := "Weles"
		fmt.Fprint(os.Stderr, title+"\n\n")
		desc := "This is a Weles server.   You can find out more about Weles at [http://tbd.tbd](http://tbd.tbd)."
		if desc != "" {
			fmt.Fprintf(os.Stderr, desc+"\n\n")
		}
		fmt.Fprintln(os.Stderr, flag.CommandLine.FlagUsages())
	}
	// parse the CLI flags
	flag.Parse()

	var yap parser.Parser
	am, err := artifacts.NewArtifactManager(
		artifactDBName,
		artifactDBLocation,
		notifierChannelCap,
		activeWorkersCap,
		artifactDownloadQueueCap)
	exitOnErr("failed to initialize ArtifactManager ", err)
	bor := client.NewBorutaClient(borutaAddress)
	djm := manager.NewDryadJobManager()
	jm := controller.NewJobManager(am, &yap, bor, borutaRefreshPeriod, djm)

	api := operations.NewWelesAPI(swaggerSpec)
	// get server with flag values filled out
	srv = server.NewServer(api)

	defer srv.Shutdown()

	apiDefaults.Managers = server.NewManagers(jm, am)

	srv.WelesConfigureAPI(&apiDefaults)
	err = srv.Serve()
	exitOnErr("failed to serve the API", err)
}
