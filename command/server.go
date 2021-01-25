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
	"github.com/superisaac/jointrpc/datadir"
	//intf "github.com/superisaac/jointrpc/intf/jointrpc"
	"github.com/superisaac/jointrpc/server"
	service "github.com/superisaac/jointrpc/service"
	"github.com/superisaac/jointrpc/service/builtin"
	"github.com/superisaac/jointrpc/service/neighbor"
	"github.com/superisaac/jointrpc/service/vars"
	//misc "github.com/superisaac/jointrpc/misc"
	"github.com/superisaac/jointrpc/rpcrouter"
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
	pHttpBind := serverFlags.String("http_bind", "", "http address and port")

	serverFlags.Parse(os.Args[2:])
	if *pDatadir != "" {
		datadir.SetDatadir(*pDatadir)
	}

	router := rpcrouter.NewRouter("server")
	router.Config.ParseDatadir()
	router.Config.SetupLogger()

	var opts []grpc.ServerOption
	var httpOpts []server.HTTPOptionFunc
	// server bind
	bind := *pBind
	if bind == "" {
		bind = router.Config.Server.Bind
	}

	httpBind := *pHttpBind
	if httpBind == "" {
		httpBind = router.Config.Server.HttpBind
	}

	// tls settings
	certFile := *pCertFile
	if certFile == "" {
		certFile = router.Config.Server.TLS.CertFile
	}

	keyFile := *pKeyFile
	if keyFile == "" {
		keyFile = router.Config.Server.TLS.KeyFile
	}

	if certFile != "" && keyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			panic(err)
		}
		opts = append(opts, grpc.Creds(creds))
		httpOpts = append(httpOpts, server.WithTLS(certFile, keyFile))
	}

	rootCtx := server.ServerContext(context.Background(), router)

	srvs := []service.IService{
		builtin.NewBuiltinService(),
		neighbor.NewNeighborService(),
		vars.NewVarsService(),
	}

	go func() {
		// start services after grpc server starts
		time.Sleep(100 * time.Millisecond)
		for _, srv := range srvs {
			service.TryStartService(rootCtx, srv)
		}
	}()
	if httpBind != "" {
		log.Infof("http server starts at %s", httpBind)
		go server.StartHTTPServer(rootCtx, httpBind, httpOpts...)
	}
	log.Infof("server starts at %s", bind)
	server.StartGRPCServer(rootCtx, bind, opts...)
}
