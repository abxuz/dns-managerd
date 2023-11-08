package dao

import (
	"context"
	"os"
)

var Local = &dLocal{}

type dLocal struct {
}

func (d *dLocal) GetConfig(ctx context.Context, k string) ([]byte, error) {
	return os.ReadFile(k)
}

func (d *dLocal) SetConfig(ctx context.Context, k string, v []byte) error {
	return os.WriteFile(k, v, 0655)
}
