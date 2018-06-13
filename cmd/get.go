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
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [trade_id] [gp_index]",
	Short: "Get record identified by trade_id and gp",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("requires at least two arguments: trade_id and gp_index")
		}
		gpIdx, err := strconv.Atoi(args[1])
		if err != nil || (gpIdx < 0 || gpIdx > 51) {
			return fmt.Errorf("invalid gp index specified (0<gpIdx<51): %s", args[0])
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, err := cmd.Flags().GetString(opt.Address)
		if err != nil {
			return fmt.Errorf("Error while parsing %v", opt.Address)
		}
		tls, err := cmd.Flags().GetBool(opt.TLS)
		if err != nil {
			return fmt.Errorf("Error while parsing flag %v", opt.TLS)
		}
		caFile, err := cmd.Flags().GetString(opt.CaFile)
		if err != nil {
			return fmt.Errorf("Error while parsing flag %v", opt.CaFile)
		}
		oHost, err := cmd.Flags().GetString(opt.ServerHostOverride)
		if err != nil {
			return fmt.Errorf("Error while parsing %v", opt.ServerHostOverride)
		}
		timeout, err := cmd.Flags().GetDuration(opt.Timeout)
		if err != nil {
			return fmt.Errorf("Error while parsing %v", opt.Timeout)
		}
		clientOpt := opt.ClientOpt{
			TLS:                &tls,
			CaFile:             &caFile,
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
		result, err := client.GetIdcTradePrices(ctx, &pb.RecordId{
			TradeId: args[0],
			GpIndex: args[1],
		})
		if err != nil {
			fmt.Errorf("Error in RPC call %v", err)
		}
		fmt.Print(*result)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().Bool(opt.TLS, false, "Connection uses TLS if true, else plain TCP")
	getCmd.Flags().String(opt.CaFile, "", "The file containning the CA root cert file")
	getCmd.Flags().String(opt.ServerHostOverride, "", "The server name use to verify the hostname returned by TLS handshake")
	getCmd.Flags().String(opt.Address, "localhost:12889", "The address of the server <hostname:port>")
	getCmd.Flags().Duration(opt.Timeout, 600, "Duration in seconds to timeout the request")
}
