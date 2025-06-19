package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"tablelink_test/internal/delivery/grpc"
	"tablelink_test/internal/repository"
	"tablelink_test/internal/usecase"
	"tablelink_test/pkg"
	"time"

	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

type Config struct {
	App struct {
		GrpcPort  int    `yaml:"grpc_port"`
		JwtSecret string `yaml:"jwt_secret"`
	} `yaml:"app"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
}

func main() {
	cfg := loadConfig()

	db, err := pkg.NewDB(pkg.DBConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
	})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	redisClient := pkg.NewRedis(pkg.RedisConfig{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	jwtManager := pkg.NewJWTManager(cfg.App.JwtSecret, time.Hour*24)

	userRepo := repository.NewUserRepositoryPostgres(db)
	roleRepo := repository.NewRoleRepositoryPostgres(db)
	roleRightRepo := repository.NewRoleRightRepositoryPostgres(db)
	sessionRepo := repository.NewSessionRepositoryRedis(redisClient)

	authUC := usecase.NewAuthUsecase(userRepo, sessionRepo, jwtManager)
	userUC := usecase.NewUserUsecase(userRepo, roleRepo, roleRightRepo, sessionRepo, jwtManager)

	grpcServer := grpc.NewServer()
	authHandler := &grpc.AuthServiceServer{AuthUC: authUC}
	userHandler := &grpc.UserServiceServer{UserUC: userUC}

	// Register gRPC service
	// TODO: Import package proto hasil generate dan register service
	// authpb.RegisterAuthServiceServer(grpcServer, authHandler)
	// userpb.RegisterUserServiceServer(grpcServer, userHandler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.App.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Server ready on port", cfg.App.GrpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func loadConfig() *Config {
	f, err := os.Open("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to open config: %v", err)
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %v", err)
	}
	return &cfg
}
