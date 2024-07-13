package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	domainpb "github.com/hts0000/ops-backend/domain/api/gen/v1"
	"github.com/hts0000/ops-backend/shared/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	// 需要匿名导入errdetails，用以获取errdetails中rpc错误映射到http错误的方法
	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
)

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseEnumNumbers: true, // 返回常量时，返回对应的number而不是string
					UseProtoNames:  true, // 返回JSON时，字段名格式为xxx_xxx格式
				},
			},
		),
		// runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, _ protoreflect.ProtoMessage) error {
		// 	// 从ctx中拿到request
		// 	// 如果没有，说明不是需要处理的请求
		// 	r, ok := ctx.Value(middleware.PassRequestKey{}).(*http.Request)
		// 	if !ok {
		// 		return nil
		// 	}

		// 	// 处理登录请求时，将token设置到header的cookie中
		// 	if r.URL.Path == "/v1/auth/login" {
		// 		// 拿到auth.Login处理之后，设置到ctx里的token
		// 		token := w.Header().Get(textproto.CanonicalMIMEHeaderKey(runtime.MetadataHeaderPrefix + "token"))
		// 		if token == "" {
		// 			return fmt.Errorf("cannot get token")
		// 		}

		// 		// 删除Grpc-Metadata-Token
		// 		w.Header().Del("Grpc-Metadata-Token")

		// 		// 设置token到cookie
		// 		http.SetCookie(w, &http.Cookie{
		// 			Name:  "token",
		// 			Value: token,
		// 			// HttpOnly: true, // 不允许在js中访问cookie
		// 			Expires: time.Now().Add(time.Hour * 24),
		// 			// Path:     "/",  // 设置只允许在Path中携带Cookie
		// 			// Secure:   true, // 只能在https请求中发送cookie
		// 		})
		// 	}

		// 	return nil
		// }),
		// runtime.WithOutgoingHeaderMatcher(func(key string) (string, bool) {
		// 	if key == textproto.CanonicalMIMEHeaderKey(runtime.MetadataHeaderPrefix+"Set-Cookie") {
		// 		fmt.Println("success match")
		// 		return "Set-Cookie", true
		// 	}
		// 	if strings.ToLower(key) == "x-test-header" {
		// 		return key, true
		// 	}
		// 	return runtime.DefaultHeaderMatcher(key)
		// }),
	)

	serverConfig := []struct {
		name         string
		addr         string
		registerFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	}{
		{
			name:         "domain",
			addr:         "localhost:18083",
			registerFunc: domainpb.RegisterDomainServiceHandlerFromEndpoint,
		},
	}

	for _, s := range serverConfig {
		err := s.registerFunc(
			c, mux, s.addr,
			[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
		)
		if err != nil {
			logger.Fatal("cannot register service", zap.String("service", s.name), zap.Error(err))
		}
	}

	// handler := middleware.NewHandler(
	// 	mux,
	// 	middleware.Debug,
	// 	middleware.Cors,
	// 	middleware.PassRequest,
	// )

	addr := ":18080"
	srv := &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mux,
	}

	logger.Info("grpc gateway started", zap.String("addr", addr))
	logger.Fatal("cannot listen and server", zap.Error(srv.ListenAndServe()))
}
