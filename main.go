package main

import (
	"context"
	"fmt"
	"github.com/metajar/coreSNMP/internal/backend"
	"github.com/metajar/coreSNMP/internal/controller"
	"github.com/metajar/coreSNMP/storage/mongodb"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

func main() {
	store := mongodb.MongoBackend{
		Host:     "127.0.0.1:27017",
		Database: "coreSNMP",
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := store.Init(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	err = store.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println(err)
	}

	newResource := controller.CoreSNMPResource{
		IP:"1.1.1.1",
		DeviceName:"test.router.dallas.network",
		Message:"does it matter",
	}

	err = backend.WriteTest(ctx, &store, newResource)
}
