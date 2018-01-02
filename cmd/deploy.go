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

package cmd

import (
	"github.com/rfancn/wegigo/pkg/deploy"
	"github.com/spf13/cobra"
)



//deployCfgFile is a ansible hosts file definition
var	inventory string
var serverAddress string
var serverPort int
var serverAssetDir string
var timeout int

var deployCmd = &cobra.Command{
	Use:   "deploy [INVENTORY]...",
	Short: "Run wegigo deploy server or deploy with config file",
	Long: `If no INVENTORY then run as deploy server,
Otherwise deploy with INVENTORY config file`,
	Run: func(cmd *cobra.Command, args []string) {
		argsLen := len(args)
		if argsLen == 0 {
			deploy.RunServerMode(serverAddress, serverPort, serverAssetDir, timeout)
		}else if argsLen == 1 {
			inventory = args[0]
			deploy.RunCliMode(inventory, timeout)
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	//deployCmd.Flags().StringVarP(&inventory, "inventory", "i", "", "ansible deployment inventory")
	//deployCmd.Flags().BoolVarP(&isServerMode, "server", "s", false, "Run deploy server")
	deployCmd.Flags().IntVarP(&timeout, "timeout", "t", 30, "timeout for deployment[minutes]")
	deployCmd.Flags().StringVarP(&serverAddress, "bind", "b", "127.0.0.1", "bind address for deploy server")
	deployCmd.Flags().IntVarP(&serverPort, "port", "p", 8443, "listen port for deploy server")
	deployCmd.Flags().StringVarP(&serverAssetDir, "assetDir", "a", "", "asset root dir for deploy server, it not specify, use internal one")
}

/**
// loadConfig reads in config file and ENV variables if set.
func loadConfig(cfgFile string) error {
	if deployCfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(deployCfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	return viper.ReadInConfig()
}
**/

