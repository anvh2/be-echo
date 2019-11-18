package integration

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/anvh2/be-echo/grpc-gen/echo"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var (
	echoClient = pb.NewEchoServiceClient(getLocalConn(55210))
)

func getLocalConn(port int) *grpc.ClientConn {
	conn, _ := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	return conn
}

func TestEcho(t *testing.T) {
	res, err := echoClient.Echo(context.Background(), &pb.EchoRequest{
		Msg: "Hello World",
	})

	assert.Nil(t, err)
	fmt.Printf("RES: %v\n", res)
}
