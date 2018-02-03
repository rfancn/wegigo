package cmd

import (
	"log"
	"github.com/spf13/cobra"
	"github.com/rfancn/wegigo/pkg/wxmp"
)

var wxmpCmdArg = &wxmp.WxmpCmdArgument{}

var wxmpCmd = &cobra.Command{
	Use:   "wxmp",
	Short: "Run Wechat Media Platform App server",
	Long: "Wegigo proxy recevie outside http requests and routed to MessageBroker",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Run Server")
		wxmp.Run(wxmpCmdArg)
	},
}

func init() {
	rootCmd.AddCommand(wxmpCmd)

	wxmpCmd.Flags().StringVarP(&wxmpCmdArg.ServerUrl, "serverUrl", "", "http://127.0.0.1:80", "server url")
	wxmpCmd.Flags().StringVarP(&wxmpCmdArg.EtcdUrl, "etcdUrl", "", "http://127.0.0.1:2379", "etcd url")

	wxmpCmd.Flags().StringVarP(&wxmpCmdArg.RabbitmqUrl, "rabbitmqUrl", "", "amqp://guest:guest@127.0.0.1:5672/", "rabbitmq url")

	wxmpCmd.Flags().StringVarP(&wxmpCmdArg.AssetDir, "assetDir", "a", "", "external asset dir")
	wxmpCmd.Flags().StringVarP(&wxmpCmdArg.AppPluginDir, "appPluginDir", "p", "apps", "app plugin dir")

	wxmpCmd.Flags().IntVarP(&wxmpCmdArg.AppConcurrency, "appConcurrency", "", 1, "concurrency number")

}
