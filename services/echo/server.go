package echo

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"

	pb "github.com/anvh2/be-echo/grpc-gen/echo"
	"github.com/anvh2/be-echo/plugins/monitor"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

// Server -
type Server struct{}

//NewServer -
func NewServer() *Server {
	return &Server{}
}

// Run -
func (s *Server) Run() error {
	port := viper.GetInt("service.grpc_port")
	// create a listener on TCP port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)

	}
	// init middleware for server
	uIntOpt := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		monitor.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor))

	// create a gRPC server
	grpcServer := grpc.NewServer(uIntOpt)
	// attach the Echo service to the server
	pb.RegisterEchoServiceServer(grpcServer, s)

	// enable prometheus
	httpPort := viper.GetInt("service.http_port")
	mux := http.NewServeMux()
	mux.Handle("/debug/vars", expvar.Handler())
	mux.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%d", httpPort), mux)

	// init monitor
	// monitor.InitCounterError()

	// start the server
	fmt.Printf("server is listening on port: %d", port)
	return grpcServer.Serve(lis)
}

// Echo -
func (s *Server) Echo(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {

	return &pb.EchoResponse{
		Msg: "Hello",
		Error: &pb.Error{
			Domain:  "zpi",
			Code:    -12013,
			Message: "OK",
		},
	}, nil
}
