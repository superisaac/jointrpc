package command

import (
	"context"
	"flag"
	log "github.com/sirupsen/logrus"
	"time"
	//"net"
	"os"
	//"fmt"
	//grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	datadir "github.com/superisaac/jointrpc/datadir"
	//intf "github.com/superisaac/jointrpc/intf/jointrpc"
	server "github.com/superisaac/jointrpc/server"
	service "github.com/superisaac/jointrpc/service"
	builtin "github.com/superisaac/jointrpc/service/builtin"
	mirror "github.com/superisaac/jointrpc/service/mirror"
	vars "github.com/superisaac/jointrpc/service/vars"
	//misc "github.com/superisaac/jointrpc/misc"
	//"github.com/superisaac/jointrpc/rpcrouter"
	//handler "github.com/superisaac/jointrpc/rpcrouter/handler"
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

	cfg := datadir.NewConfig()
	cfg.ParseDatadir()
	cfg.SetupLogger()

	//go StartHTTPd(*httpBind)
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

	rootCtx := server.ServerContext(context.Background(), nil, nil)

	srvs := []service.IService{
		builtin.NewBuiltinService(),
		mirror.NewMirrorService(),
		vars.NewVarsService(),
	}

	go func() {
		log.Debugf("ddddd")
		time.Sleep(100 * time.Millisecond)
		for _, srv := range srvs {
			service.TryStartService(rootCtx, srv)
		}
	}()
	server.StartServer(rootCtx, bind, opts...)
}
