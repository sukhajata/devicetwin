package grpchelper

import (
	"github.com/sukhajata/devicetwin/pkg/errorhelper"
	"google.golang.org/grpc"
)

// ConnectGRPC connects to a given address and returns the connection
func ConnectGRPC(address string) (*grpc.ClientConn, error) {
	for {
		retries := 0
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			retries++
			if retries > 5 {
				return nil, err
			}
			errorhelper.StartUpError(err)
			continue
		}

		return conn, nil
	}
}
