package app

import "context"

type WebApp interface {
	Start(port string) error
	Stop(ctx context.Context) error
}
