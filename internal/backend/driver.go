package backend

import (
	"context"
	"fmt"
	"github.com/metajar/coreSNMP/internal/controller"
)

func WriteToBackend(ctx context.Context, s CoreSNMPBackend) {
	fmt.Println("Writing to the storage array.")
}

func WriteTest(ctx context.Context, s CoreSNMPBackend, resource controller.CoreSNMPResource) error {
	err := s.TestWrite(ctx, resource)
	if err != nil {
		return err
	}
	return nil
}