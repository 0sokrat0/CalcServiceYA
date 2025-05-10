package auth_test

import (
	"context"
	"net"
	"testing"

	"auth/config"
	"auth/internal/app"
	authInfra "auth/internal/infrastructure/auth"
	hashpass "auth/internal/infrastructure/hashPass"
	sqliteRepo "auth/internal/infrastructure/persistence/sqlite"
	"auth/internal/presentation/grpc/handlers"
	sqliteconn "auth/pkg/db/sqlite_conn"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	authpb "auth/pkg/gen/api"
)

const bufSize = 1024 * 1024

func dialer(t *testing.T) (*grpc.ClientConn, func()) {
	lis := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer()

	// === Инициализация exactly как в main.go ===
	cfg := config.GetConfig("../config/config.yaml")
	logger := zap.NewExample()
	sqliteDB, err := sqliteconn.NewSQLiteDB(cfg, logger)
	require.NoError(t, err)
	jwtSvc := authInfra.NewJWTService(cfg.JWT.JWTSecret, cfg.JWT.AccessTokenDuration, cfg.JWT.RefreshTokenDuration)
	passHasher := hashpass.NewPassHasher()
	userRepo := sqliteRepo.NewUserSQLiteRepository(sqliteDB)
	userSvc := app.NewUserService(userRepo, jwtSvc, passHasher)

	// === Регистрируем handler, а не всю обёртку g.NewServer ===
	authHandler := handlers.NewAuthHandler(userSvc)
	authpb.RegisterAuthServer(grpcServer, authHandler)

	go func() { require.NoError(t, grpcServer.Serve(lis)) }()

	conn, err := grpc.DialContext(
		context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	return conn, func() {
		conn.Close()
		grpcServer.Stop()
	}
}

func TestRegisterAndLogin(t *testing.T) {
	conn, cleanup := dialer(t)
	defer cleanup()

	client := authpb.NewAuthClient(conn)

	// 1) Register
	regReq := &authpb.RegisterRequest{
		Email:    "test@example.com",
		Password: "secret123",
	}
	regResp, err := client.Register(context.Background(), regReq)
	require.NoError(t, err)
	require.NotEmpty(t, regResp.UserId, "user_id should be returned")

	// 2) Login
	loginReq := &authpb.LoginRequest{
		Email:    "test@example.com",
		Password: "secret123",
	}
	loginResp, err := client.Login(context.Background(), loginReq)
	require.NoError(t, err)
	require.NotEmpty(t, loginResp.Access, "access token should be returned")
	require.NotEmpty(t, loginResp.Refresh, "refresh token should be returned")
}
