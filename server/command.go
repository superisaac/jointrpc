package server

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	//"fmt"
	bridge "github.com/superisaac/rpctube/bridge"
	client "github.com/superisaac/rpctube/client"
	datadir "github.com/superisaac/rpctube/datadir"
	intf "github.com/superisaac/rpctube/intf/tube"
	tube "github.com/superisaac/rpctube/tube"
	handler "github.com/superisaac/rpctube/tube/handler"

	grpc "google.golang.org/grpc"
	credentials "google.golang.org/grpc/credentials"
)

func CommandStartServer() {
	serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
	pBind := serverFlags.String("b", "", "The grpc server address and port")
	pDatadir := serverFlags.String("d", "", "The datadir to store configs")
	pCertFile := serverFlags.String("cert", "", "tls cert file")
	pKeyFile := serverFlags.String("key", "", "tls key file")
	//httpBind := serverFlags.String("httpd", "127.0.0.1:50056", "http address and port")

	serverFlags.Parse(os.Args[2:])
	if *pDatadir != "" {
		datadir.SetDatadir(*pDatadir)
	}

	//go StartHTTPd(*httpBind)
	cfg := datadir.GetConfig()
	var opts []grpc.ServerOption

	// server bind
	bind := *pBind
	if bind == "" {
		bind = cfg.Server.Bind
	}

	// tls settings
	certFile := *pCertFile
	if certFile == "" {
		certFile = cfg.Server.TLS.CertFile
	}

	keyFile := *pKeyFile
	if keyFile == "" {
		keyFile = cfg.Server.TLS.KeyFile
	}

	if certFile != "" || keyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			panic(err)
		}
		opts = append(opts, grpc.Creds(creds))
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	StartServer(ctx, bind, opts...)
}

func StartServer(ctx context.Context, bind string, opts ...grpc.ServerOption) {
	lis, err := net.Listen("tcp", bind)
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	} else {
		log.Infof("entry server listen at %s", bind)
	}

	tube.Tube().Start(ctx)

	handlerCtx, _ := context.WithCancel(ctx)
	handler.Builtin().Start(handlerCtx)

	cfg := datadir.GetConfig()	
	if len(cfg.Cluster.StaticPeers) > 0 {
		// generate server entry from peers
		var serverEntries []client.ServerEntry
		for _, peer := range cfg.Cluster.StaticPeers {
			serverEntries = append(serverEntries, client.ServerEntry{
				Address:  peer.Address,
				CertFile: peer.CertFile,
			})
		}
		go bridge.StartNewBridge(ctx, serverEntries)
	}

	//var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	//hello.RegisterHelloServer(grpcServer, s)
	intf.RegisterJSONRPCTubeServer(grpcServer, NewJSONRPCTubeServer())
	grpcServer.Serve(lis)
}
