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
	"RAE_AGE2_JAVAWEB/pdcpp/options"
	pb "RAE_AGE2_JAVAWEB/pdcpp/pdc_trade"
	"RAE_AGE2_JAVAWEB/pdcpp/service"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	"log"
	"net"
	"net/http"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the pdc cache server",
	Long: `Starts the server to the specified address/port throught the related flag. The max number of gps is also set for switching purposes.
	Example:

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fl, err := retreiveFlags(cmd)
		if err != nil {
			return fmt.Errorf("Couldn't parse flags")
		}
		if *fl.httpProxy {
			return serveHTTP(fl)
		}
		return startServing(fl)
	},
}

func init() {
	//adds this command as child of rootCmd
	rootCmd.AddCommand(serveCmd)
	//Parsing flags
	serveCmd.Flags().Bool(options.TLS, false, "Connection uses TLS if true, else plain TCP")
	serveCmd.Flags().Bool(options.HTTPReverseProxy, false, "Starts an HTTP reverse proxy server which connects to the GRPC server")
	serveCmd.Flags().String(options.CertFile, "", "The TLS cert file")
	serveCmd.Flags().String(options.KeyFile, "", "The TLS key file")
	serveCmd.Flags().IntP(options.Port, "p", 12889, "The server port - default 12889")
	serveCmd.Flags().Int(options.PortHTTP, 22888, "The http reverse proxy server port - default 22888")
	serveCmd.Flags().Int(options.GpNUM, 52, "The server port")
}

//Store flags
type serveFlags struct {
	tls       *bool
	certFile  *string
	keyFile   *string
	port      *int
	httpPort  *int
	gpNum     *int
	httpProxy *bool
}

func retreiveFlags(cmd *cobra.Command) (serveFlags, error) {
	tls, err := cmd.Flags().GetBool(options.TLS)
	if err != nil {
		return serveFlags{}, fmt.Errorf("Error while parsing flag %v", options.TLS)
	}
	certFile, err := cmd.Flags().GetString(options.CertFile)
	if err != nil {
		return serveFlags{}, fmt.Errorf("Error while parsing flag %v", options.CertFile)
	}
	keyFile, err := cmd.Flags().GetString(options.KeyFile)
	if err != nil {
		return serveFlags{}, fmt.Errorf("Error while parsing flag %v", options.KeyFile)
	}
	port, err := cmd.Flags().GetInt(options.Port)
	if err != nil {
		return serveFlags{}, fmt.Errorf("Error while parsing flag %v", options.Port)
	}
	httpPort, err := cmd.Flags().GetInt(options.PortHTTP)
	if err != nil {
		return serveFlags{}, fmt.Errorf("Error while parsing flag %v", options.PortHTTP)
	}
	gpNum, err := cmd.Flags().GetInt(options.GpNUM)
	if err != nil {
		return serveFlags{}, fmt.Errorf("Error while parsing flag %v", options.GpNUM)
	}
	httpRevProxy, err := cmd.Flags().GetBool(options.HTTPReverseProxy)
	if err != nil {
		return serveFlags{}, fmt.Errorf("Error while parsing flag %v", options.HTTPReverseProxy)
	}
	return serveFlags{
		tls:       &tls,
		certFile:  &certFile,
		keyFile:   &keyFile,
		port:      &port,
		httpPort:  &httpPort,
		gpNum:     &gpNum,
		httpProxy: &httpRevProxy,
	}, nil
}

func startServing(f serveFlags) error {

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *f.port))
	if err != nil {
		fmt.Errorf("failed to listen on port %v: %v", *f.port, err)
	}
	var opts []grpc.ServerOption
	if *f.tls {
		if *f.certFile == "" {
			*f.certFile = testdata.Path("server.pem")
		}
		if *f.keyFile == "" {
			*f.keyFile = testdata.Path("server.key")
		}
		creds, err := credentials.NewServerTLSFromFile(*f.certFile, *f.keyFile)
		if err != nil {
			fmt.Errorf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPdcTradePricesServiceServer(grpcServer, service.GetPriceProviderService(*f.gpNum))
	defer log.Printf("serving at port %v", *f.port)
	return grpcServer.Serve(lis)
}

func serveHTTP(f serveFlags) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterPdcTradePricesServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("0.0.0.0:%d", *f.port), opts)
	if err != nil {
		return err
	}

	defer log.Printf("http reverse proxy serving at port %v", *f.httpPort)
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *f.httpPort), mux)
}
