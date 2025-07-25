package services

import (
	"context"
	"os"
	"time"

	"github.com/burapha44/example/cache"
	"github.com/burapha44/example/config"
	"github.com/burapha44/example/db"

	"github.com/gofiber/fiber/v2"
)

type BaseService struct {
	InitializedAt time.Time
}

func NewBaseService() *BaseService {
	return &BaseService{
		InitializedAt: time.Now(),
	}
}

type DependencyStatus struct {
	Status         string `json:"status"`
	Message        string `json:"message,omitempty"`
	ResponseTimeMs int    `json:"response_time_ms,omitempty"`
}

type ServerInfoResponse struct {
	Status       string                      `json:"status"`
	Message      string                      `json:"message,omitempty"`
	Timestamp    time.Time                   `json:"timestamp"`
	Version      string                      `json:"version"`
	ServiceName  string                      `json:"service_name"`
	Environment  string                      `json:"environment"`
	Hostname     string                      `json:"hostname,omitempty"`
	Uptime       string                      `json:"uptime,omitempty"` // Or int for seconds
	Dependencies map[string]DependencyStatus `json:"dependencies,omitempty"`
}

func (s *BaseService) ServerInfo(c *fiber.Ctx) ServerInfoResponse {
	hostname, _ := os.Hostname()
	uptime := time.Since(s.InitializedAt).String()

	dbStatus := DependencyStatus{Status: "DOWN", Message: "Not configured"}
	if config.Conf.PostgresUser != "" && db.PostgresConn != nil {
		ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
		defer cancel()
		start := time.Now()
		if err := db.PostgresConn.Ping(ctx); err != nil {
			dbStatus = DependencyStatus{Status: "DOWN", Message: err.Error()}
		} else {
			dbStatus = DependencyStatus{Status: "UP", Message: "Connected", ResponseTimeMs: int(time.Since(start).Milliseconds())}
		}
	} else if config.Conf.PostgresUser != "" {
		dbStatus = DependencyStatus{Status: "N/A", Message: "Database disabled by config"}
	}

	redisStatus := DependencyStatus{Status: "DOWN", Message: "Not configured"}
	if config.Conf.RedisHost != "" && cache.GetCacheClient() != nil {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()

		start := time.Now()
		if err := cache.GetCacheClient().Ping(ctx).Err(); err != nil {
			redisStatus = DependencyStatus{Status: "DOWN", Message: err.Error()}
		} else {
			redisStatus = DependencyStatus{Status: "UP", Message: "Connected", ResponseTimeMs: int(time.Since(start).Milliseconds())}
		}
	} else if config.Conf.RedisHost != "" {
		redisStatus = DependencyStatus{Status: "N/A", Message: "Cache disabled by config"}
	}

	overallStatus := "UP"
	if dbStatus.Status != "UP" || redisStatus.Status != "UP" {
		overallStatus = "DEGRADED"
	}
	if dbStatus.Status == "DOWN" && config.Conf.PostgresUser != "" {
		overallStatus = "DOWN"
	}
	return ServerInfoResponse{
		Status:      overallStatus,
		Message:     "All systems operational",
		Timestamp:   time.Now().UTC(),
		Version:     config.Conf.Version,
		ServiceName: config.Conf.ServiceName,
		Environment: config.Conf.Environment,
		Hostname:    hostname,
		Uptime:      uptime,
		Dependencies: map[string]DependencyStatus{
			"database": dbStatus,
			"redis":    redisStatus,
		},
	}
}
