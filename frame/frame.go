package frame

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

type GinOption func(g *gin.Engine)

type GRPCOption func(g *grpc.Server)

type IService interface {
	RegisterHttp(...GinOption)
	RegisterGRPC(...GRPCOption)
	Run(port int)
}

type service struct {
	c cmux.CMux

	// 通讯 ginSvc grpcSvc 微服务  irisSvc web服务
	ginSvc  *gin.Engine
	grpcSvc *grpc.Server
	// 监控 jaeger、prometheus
	//tracing jaegerex.ILinkTracing
}

func (s *service) RegisterHttp(opts ...GinOption) {
	if len(opts) > 0 {
		s.ginSvc = gin.Default()
		s.ginSvc.Use(func(ctx *gin.Context) {})
		for _, o := range opts {
			o(s.ginSvc)
		}
	}
}

func (s *service) RegisterGRPC(opts ...GRPCOption) {
	if len(opts) > 0 {
		s.grpcSvc = grpc.NewServer(
			grpc.StreamInterceptor(
				grpc_prometheus.StreamServerInterceptor,
			),
			grpc.UnaryInterceptor(
				grpc_prometheus.UnaryServerInterceptor,
			),
		)
		for _, o := range opts {
			o(s.grpcSvc)
		}
		grpc_prometheus.Register(s.grpcSvc)
	}
}

func (s *service) Run(port int) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println("tcp listen failed: ", err.Error())
	}
	if s.c == nil {
		s.c = cmux.New(l)
	}
	if s.ginSvc != nil {
		go s.ginSvc.RunListener(s.c.Match(cmux.HTTP1Fast()))
	}
	if s.grpcSvc != nil {
		go s.grpcSvc.Serve(
			s.c.Match(
				cmux.HTTP2HeaderField("content-type", "application/grpc"),
			),
		)
	}
	go func() {
		if err = s.c.Serve(); err != nil {
			log.Println("cmux service start failed: ", err.Error())
		}
	}()

	fmt.Println("服务已启动, 端口: ", port, "...")
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	<-ctx.Done()
	s.c.Close()
	log.Println("shutdown service")
}

func New() IService {
	return &service{}
}
