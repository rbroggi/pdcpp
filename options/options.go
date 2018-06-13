package options

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

//flag variable names
const (
	TLS                = "tls"
	CertFile           = "cert_file"
	KeyFile            = "key_file"
	Port               = "port"
	PortHTTP           = "http_port"
	GpNUM              = "gpNum"
	CaFile             = "ca_file"
	ServerHostOverride = "server_host_override"
	Timeout            = "timeout"
	Address            = "addr"
	HTTPReverseProxy   = "http_reverse_proxy"
)

//ClientOpt is a struct containing authentication info for client connections
type ClientOpt struct {
	TLS                *bool
	CaFile             *string
	ServerHostOverride *string
}

//GetClientOptions returns a []grpc.DialOption describing the connection options for a
//grpc client
func GetClientOptions(f ClientOpt) ([]grpc.DialOption, error) {
	var opts []grpc.DialOption
	if *f.TLS {
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
