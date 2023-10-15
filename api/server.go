package api

import (
	"context"
	"embed"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	db "github.com/zcubbs/tlz/db/sqlc"
	"github.com/zcubbs/tlz/internal/util"
	"github.com/zcubbs/tlz/pb"
	"github.com/zcubbs/tlz/pkg/token"
	"github.com/zcubbs/tlz/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"io/fs"
	"net"
	"net/http"
)

type Server struct {
	pb.UnimplementedTlzServer
	store           db.Store
	tokenMaker      token.Maker
	cfg             util.Config
	embedAssets     []EmbedAssetsOpts
	taskDistributor worker.TaskDistributor
}

type EmbedAssetsOpts struct {
	// The directory to embed.
	// Defaults to "assets".
	Dir    embed.FS
	Path   string
	Prefix string
}

func NewServer(store db.Store, taskDistributor worker.TaskDistributor,
	cfg util.Config, embedOpts ...EmbedAssetsOpts) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(cfg.Auth.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create new tokenMaker: %w", err)
	}

	s := &Server{
		store:           store,
		tokenMaker:      tokenMaker,
		cfg:             cfg,
		embedAssets:     embedOpts,
		taskDistributor: taskDistributor,
	}

	return s, nil
}

func (s *Server) StartGrpcServer() {
	grpcLogger := grpc.UnaryInterceptor(GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterTlzServer(grpcServer, s)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.cfg.GrpcServer.Port))
	if err != nil {
		log.Fatal("cannot listen", "error", err, "port", s.cfg.GrpcServer.Port)
	}

	log.Info("🟢 starting gRPC server", "port", s.cfg.GrpcServer.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("cannot start server: %w", err)
	}
}

func (s *Server) StartHttpGateway() {
	jsonOpts := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOpts)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := pb.RegisterTlzHandlerServer(ctx, grpcMux, s)
	if err != nil {
		log.Fatal("cannot register handler server", "error", err)
	}

	apiPath := "/api"
	mux := http.NewServeMux()
	mux.Handle(apiPath+"/", http.StripPrefix(apiPath, grpcMux))
	log.Info("serving API Gateway", "path", apiPath)

	for _, opt := range s.embedAssets {
		log.Info("serving embedded assets", "path", opt.Path)
		sub, err := fs.Sub(opt.Dir, opt.Prefix)
		if err != nil {
			log.Fatal("cannot serve embedded assets", "error", err)
		}
		dir := http.FileServer(http.FS(sub))
		mux.Handle(opt.Path, http.StripPrefix(opt.Path, dir))
	}

	log.Info("🟢 starting HTTP Gateway server", "port", s.cfg.HttpServer.Port)
	handler := HttpLogger(mux)

	withCors := cors.New(cors.Options{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"ACCEPT", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(handler)

	httpSrv := &http.Server{
		Handler:     withCors,
		Addr:        fmt.Sprintf(":%d", s.cfg.HttpServer.Port),
		ReadTimeout: s.cfg.HttpServer.ReadHeaderTimeout,
	}

	if err := httpSrv.ListenAndServe(); err != nil {
		log.Fatal("cannot start HTTP Gateway server", "error", err)
	}
}
