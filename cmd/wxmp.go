package cmd

import (
	"log"
	"github.com/spf13/cobra"
	"github.com/rfancn/wegigo/pkg/wxmp"
)

var (
	etcdAddress string
	etcdPort int
	rabbitmqAddress string
	rabbitmqPort int
	assetDir string
	appsDir string
)

var wxmpCmd = &cobra.Command{
	Use:   "wxmp",
	Short: "Run Wechat Media Platform App server",
	Long: "Wegigo proxy recevie outside http requests and routed to MessageBroker",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Run Server")

		wxmp.Run(appsDir, assetDir, etcdAddress, etcdPort, rabbitmqAddress, rabbitmqPort)
	},
}

func init() {
	rootCmd.AddCommand(wxmpCmd)

	wxmpCmd.Flags().StringVarP(&etcdAddress, "etcdAddress", "", "127.0.0.1", "etcd address")
	wxmpCmd.Flags().IntVarP(&etcdPort, "etcdPort", "", 2379, "etcd port")
	wxmpCmd.Flags().StringVarP(&rabbitmqAddress, "rabbitmqAddress", "", "localhost", "rabbitmq address")
	wxmpCmd.Flags().IntVarP(&rabbitmqPort, "rabbitmqPort", "", 5672, "rabbitmq port")

	wxmpCmd.Flags().StringVarP(&assetDir, "assetDir", "a", "", "external asset root dir for wegigo server")
	wxmpCmd.Flags().StringVarP(&appsDir, "appsDir", "i", "apps", "app modules dir")
}
