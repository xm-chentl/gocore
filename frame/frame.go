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
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
)

type GinOption func(g *gin.Engine)

type GRPCOption func(g *grpc.Server)

type IService interface {
	RegisterGin(...GinOption)
	RegisterGRPC(...GRPCOption)
	Run(port int)
}

type service struct {
	c cmux.CMux

	// 通讯
	ginSvc  *gin.Engine
	grpcSvc *grpc.Server

	// 监控 jaeger、prometheus
}

func (s *service) RegisterGin(opts ...GinOption) {
	if len(opts) > 0 {
		s.ginSvc = gin.Default()
		for _, o := range opts {
			o(s.ginSvc)
		}
	}
}

func (s *service) RegisterGRPC(opts ...GRPCOption) {
	if len(opts) > 0 {
		s.grpcSvc = grpc.NewServer()
		for _, o := range opts {
			o(s.grpcSvc)
		}
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

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	<-ctx.Done()
	s.c.Close()
	log.Println("shutdown service")
}

func New() IService {
	return &service{}
}
