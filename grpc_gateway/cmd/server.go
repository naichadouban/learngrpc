package cmd

import (
	"log"

	"github.com/naichadouban/learngrpc/grpc_gateway/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the grpc hello-word server",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recover error:%v", err)
			}
		}()
		server.Serve()
	},
}

func init() {
	serverCmd.Flags().StringVarP(&server.ServerPort, "port", "p", "8010", "server port")
	serverCmd.Flags().StringVarP(&server.CertPemPath, "cert-pem", "", "./certs/server.pem", "cert pem path")
	serverCmd.Flags().StringVarP(&server.CertKeyPath, "cert-key", "", "./certs/server.key", "cert key path")
	serverCmd.Flags().StringVarP(&server.CertName, "cert-name", "", "server", "server's hostname")
	rootCmd.AddCommand(serverCmd)
}
