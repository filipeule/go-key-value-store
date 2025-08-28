package frontend

import (
	"errors"
	"fmt"
	"key-value-store/core"
)

type FrontEnd interface {
	Start(kv *core.KeyValueStore) error
}

func NewFrontEnd(frontend string) (FrontEnd, error) {
	switch frontend {
	case "rest":
		return &RestFrontEnd{}, nil
	case "grpc":
		return &GRPCFrontEnd{}, nil
	case "":
		return nil, errors.New("frontend type not defined")
	default:
		return nil, fmt.Errorf("no such frontend %s", frontend)
	}
}