// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	opt "RAE_AGE2_JAVAWEB/pdcpp/options"
	pb "RAE_AGE2_JAVAWEB/pdcpp/pdc_trade"
	"fmt"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [trade_id]",
	Short: "Checks existence of trade_id in the cache",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, err := cmd.Flags().GetString(opt.Address_s)
		if err != nil {
			return fmt.Errorf("Error while parsing %v", opt.Address_s)
		}
		tls, err := cmd.Flags().GetBool(opt.Tls_c)
		if err != nil {
			return fmt.Errorf("Error while parsing flag %v", opt.Tls_c)
		}
		ca_file, err := cmd.Flags().GetString(opt.Ca_file_c)
		if err != nil {
			return fmt.Errorf("Error while parsing flag %v", opt.Ca_file_c)
		}
		oHost, err := cmd.Flags().GetString(opt.Server_host_override_c)
		if err != nil {
			return fmt.Errorf("Error while parsing %v", opt.Server_host_override_c)
		}
		timeout, err := cmd.Flags().GetDuration(opt.Timeout_c)
		if err != nil {
			return fmt.Errorf("Error while parsing %v", opt.Timeout_c)
		}
		clientOpt := opt.ClientOpt{
			Tls:                &tls,
			CaFile:             &ca_file,
			ServerHostOverride: &oHost,
		}
		opts, err := opt.GetClientOptions(clientOpt)
		if err != nil {
			return fmt.Errorf("fail to get options: %v", err)
		}
		conn, err := grpc.Dial(addr, opts...)
		if err != nil {
			return fmt.Errorf("fail to dial: %v", err)
		}
		defer conn.Close()
		client := pb.NewPdcTradePricesServiceClient(conn)
		ctx, _ := context.WithTimeout(context.Background(), timeout*time.Second)
		result, err := client.CheckTrade(ctx, &pb.TradeId{
			TradeId: args[0],
		})
		if err != nil {
			fmt.Errorf("Error in RPC call %v", err)
		}
		fmt.Print((*result).GetValue())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().Bool(opt.Tls_c, false, "Connection uses TLS if true, else plain TCP")
	checkCmd.Flags().String(opt.Ca_file_c, "", "The file containning the CA root cert file")
	checkCmd.Flags().String(opt.Server_host_override_c, "", "The server name use to verify the hostname returned by TLS handshake")
	checkCmd.Flags().String(opt.Address_s, "localhost:12889", "The address of the server <hostname:port>")
	checkCmd.Flags().Duration(opt.Timeout_c, 600, "Duration in seconds to timeout the request")
}
