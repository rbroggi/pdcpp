package options

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

//flag variable names
const (
	Tls_c                  = "tls"
	Cert_file_c            = "cert_file"
	Key_file_c             = "key_file"
	Port_c                 = "port"
	PortHttp_c             = "http_port"
	GpNum_c                = "gpNum"
	Ca_file_c              = "ca_file"
	Server_host_override_c = "server_host_override"
	Timeout_c              = "timeout"
	Address_s              = "addr"
)

type ClientOpt struct {
	Tls                *bool
	CaFile             *string
	ServerHostOverride *string
}

func GetClientOptions(f ClientOpt) ([]grpc.DialOption, error) {
	var opts []grpc.DialOption
	if *f.Tls {
		if *f.CaFile == "" {
			*f.CaFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*f.CaFile, *f.ServerHostOverride)
		if err != nil {
			return nil, fmt.Errorf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	return opts, nil
}
