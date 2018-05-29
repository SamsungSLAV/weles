/*
 *  Copyright (c) 2017-2018 Samsung Electronics Co., Ltd All Rights Reserved
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License
 */

package mock

//go:generate ../../bin/dev-tools/mockgen -package mock -destination=./jobscontroller.go git.tizen.org/tools/weles/controller JobsController

// TODO: fix relative path to absolute
// You need to have Boruta checked out next to Weles for below go generate to work.
// Reflect mode will not work in below case.
//go:generate ../../bin/dev-tools/mockgen -package mock -destination=./requests.go -source=../../../boruta/boruta.go

//go:generate ../../bin/dev-tools/mockgen -package mock -destination=./boruter.go git.tizen.org/tools/weles/controller Boruter

//go:generate ../../bin/dev-tools/mockgen -package mock -destination=./downloader.go git.tizen.org/tools/weles/controller Downloader

//go:generate ../../bin/dev-tools/mockgen -package mock -destination=./dryader.go git.tizen.org/tools/weles/controller Dryader

//go:generate ../../bin/dev-tools/mockgen -package mock -destination=./parser.go git.tizen.org/tools/weles/controller Parser
