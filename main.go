// Copyright © 2017 Ryan Fan <reg_info@qq.com>
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

package main

//go:generate go-bindata -o pkg/deploy/bindata.go -pkg deploy asset/pkg/deploy/...
//go:generate go-bindata -o pkg/wxmp/bindata.go -pkg wxmp asset/pkg/wxmp/...
//go:generate go-bindata -o sdk/server/bindata.go -pkg server asset/vendors/...

import "github.com/rfancn/wegigo/cmd"

func main() {
	cmd.Execute()
}
