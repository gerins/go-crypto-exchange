package model

import "context"

type Engine interface {
	Execute(ctx context.Context, order Order) []Trade
}
