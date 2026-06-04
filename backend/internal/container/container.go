package container

import (
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/infrastructure/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Container struct {
	DB     *pgxpool.Pool
	Logger *zap.Logger

	UserRepository    *postgres.UserRepository
	ProjectRepository *postgres.ProjectRepository
}
