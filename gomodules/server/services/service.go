package services

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

type ServiceContext struct {
	Ctx          context.Context
	Logger       runtime.Logger
	DbConnection *sql.DB
	NakamaModule runtime.NakamaModule
}
