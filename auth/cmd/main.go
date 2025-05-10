// cmd/main.go
package main

import (
	"auth/config"
	"auth/internal/app"
	authInfra "auth/internal/infrastructure/auth"
	hashpass "auth/internal/infrastructure/hashPass"
	"auth/internal/infrastructure/persistence/sqlite"
	"auth/internal/presentation/grpc"
	sqliteconn "auth/pkg/db/sqlite_conn"
	"auth/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// ctx := context.Background()

	cfg := config.GetConfig()
	log := logger.InitLogger(cfg)

	log.Info("GetConfig", zap.Any("config", cfg))

	// pgPool, err := db.NewPG(cfg, log, ctx)
	// if err != nil {
	// 	log.Fatal("db connection failed", zap.Error(err))
	// }

	sqliteConn, err := sqliteconn.NewSQLiteDB(cfg, log)
	if err != nil {
		log.Fatal("db connection failed", zap.Error(err))
	}

	jwtSvc := authInfra.NewJWTService(
		cfg.JWT.JWTSecret,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
	)
	passHasher := hashpass.NewPassHasher()

	// userRepo := postgres.NewUserPGRepository(pgPool.DB)
	userRepo := sqlite.NewUserSQLiteRepository(sqliteConn)

	userSvc := app.NewUserService(userRepo, jwtSvc, passHasher)

	grpcServer := grpc.NewServer(log, cfg, userSvc)

	if err := grpcServer.Start(); err != nil {
		log.Fatal("gRPC server failed", zap.Error(err))
	}
}
