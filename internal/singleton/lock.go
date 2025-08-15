package singleton

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Guard struct {
	client *mongo.Client
	col    *mongo.Collection
	owner  string
	stop   chan struct{}
}

func EnsureSingletonOrExit(ctx context.Context, mongoURI, dbName, lockID string) *Guard {
	if mongoURI == "" || dbName == "" {
		log.Println("[singleton] mongo uri/db not provided; skipping lock (NOT SAFE)")
		return &Guard{stop: make(chan struct{})}
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil { log.Fatalf("[singleton] connect mongo: %v", err) }
	col := client.Database(dbName).Collection("locks")
	_, _ = col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "expiresAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	owner := fmt.Sprintf("%s/%s/%d", hostname(), runtime.GOOS, os.Getpid())
	now := time.Now()
	doc := bson.M{"_id": lockID, "owner": owner, "createdAt": now, "expiresAt": now.Add(2 * time.Minute)}
	_, err = col.InsertOne(ctx, doc)
	if err != nil {
		filter := bson.M{"_id": lockID, "expiresAt": bson.M{"$lt": now}}
		update := bson.M{"$set": bson.M{"owner": owner, "expiresAt": time.Now().Add(2 * time.Minute)}}
		res, err2 := col.UpdateOne(ctx, filter, update)
		if err2 != nil || res.ModifiedCount == 0 {
			log.Fatalf("[singleton] another instance holds the lock (%s). Exit.", lockID)
		}
	}
	g := &Guard{client: client, col: col, owner: owner, stop: make(chan struct{})}
	go func() { // keepalive
		t := time.NewTicker(30 * time.Second); defer t.Stop()
		for {
			select {
			case <-t.C:
				_, err := col.UpdateByID(context.Background(), lockID,
					bson.M{"$set": bson.M{"expiresAt": time.Now().Add(2 * time.Minute), "owner": owner}})
				if err != nil { log.Printf("[singleton] keepalive failed: %v", err) }
			case <-g.stop:
				return
			}
		}
	}()
	go func() { // release on SIGTERM
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch
		_ = g.Release(context.Background(), lockID)
		os.Exit(0)
	}()
	log.Printf("[singleton] lock acquired by %s", owner)
	return g
}

func (g *Guard) Release(ctx context.Context, lockID string) error {
	close(g.stop)
	_, _ = g.col.DeleteOne(ctx, bson.M{"_id": lockID, "owner": g.owner})
	if g.client != nil { _ = g.client.Disconnect(ctx) }
	log.Printf("[singleton] lock released by %s", g.owner)
	return nil
}

func hostname() string {
	h, _ := os.Hostname()
	if h == "" { h = "host" }
	return h + "-" + strconv.Itoa(os.Getpid())
}

// ForceRelease — силой удаляет документ лока.
func ForceRelease(ctx context.Context, mongoURI, dbName, lockID string) error {
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil { return err }
    defer client.Disconnect(ctx)
    col := client.Database(dbName).Collection("locks")
    _, err = col.DeleteOne(ctx, bson.M{"_id": lockID})
    return err
}

// CurrentOwner — посмотреть, кто держит лок.
func CurrentOwner(ctx context.Context, mongoURI, dbName, lockID string) (owner string, expiresAt time.Time, err error) {
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil { return "", time.Time{}, err }
    defer client.Disconnect(ctx)
    col := client.Database(dbName).Collection("locks")
    var doc struct {
        Owner     string    `bson:"owner"`
        ExpiresAt time.Time `bson:"expiresAt"`
    }
    if err := col.FindOne(ctx, bson.M{"_id": lockID}).Decode(&doc); err != nil {
        return "", time.Time{}, err
    }
    return doc.Owner, doc.ExpiresAt, nil
}
