package backend

import (
	"context"
	"github.com/metajar/coreSNMP/internal/controller"
)

type CoreSNMPBackend interface {
	Init(ctx context.Context) error
	Close(ctx context.Context) error
	Put(ctx context.Context) error
	Get(ctx context.Context) error
	Update(ctx context.Context) error
	TestWrite(ctx context.Context, c controller.CoreSNMPResource) error
}
