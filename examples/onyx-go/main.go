package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/OnyxDevTools/onyx-cli/examples/onyx-go/onyx"
)

func main() {
	ctx := context.Background()

	// Let the SDK resolve config from env or onyx-database.json; allow override via ONYX_CONFIG_PATH.
	cfg := onyx.Config{}

	db, err := onyx.New(ctx, cfg)
	if err != nil {
		log.Fatalf("init onyx client: %v", err)
	}

	id := "cli-go-e2e"
	// Clean slate
	_, _ = db.Users().DeleteByID(ctx, id)

	created, err := db.Users().Save(ctx, onyx.User{
		Id:       id,
		Username: "go-e2e",
		Email:    "go-e2e@example.com",
	})
	if err != nil {
		log.Fatalf("create: %v", err)
	}
	fmt.Printf("Created: %+v\n", created)

	fetched, err := db.Users().FindByID(ctx, id)
	if err != nil {
		log.Fatalf("get: %v", err)
	}
	fmt.Printf("Fetched: %+v\n", fetched)

	updates := onyx.NewUserUpdates().
		SetUsername("go-e2e-updated").
		SetEmail("go-e2e-updated@example.com")
	_, err = db.Users().Where(onyx.Eq("id", id)).SetUserUpdates(updates).Update(ctx)
	if err != nil {
		log.Fatalf("update: %v", err)
	}
	fmt.Println("Updated.")

	deleted, err := db.Users().DeleteByID(ctx, id)
	if err != nil {
		log.Fatalf("delete: %v", err)
	}
	fmt.Printf("Deleted: %+v\n", deleted)

	// Small delay to keep output readable when run in fast CI.
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Go example CLI+SDK compatibility test passed.")
}
