package main

import (
	"context"
	"fmt"
	"time"

	"testgin/database"

	"github.com/fxtlabs/date"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Kiosk struct {
	ID    string `json:"id"`
	State string `json:"state"`
	Name  string `json:"name"`
}

type Doctor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	day := date.TodayUTC()
	fmt.Printf("var1 = %T\n", day)
	var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
	var kioskCollection *mongo.Collection = database.OpenCollection(database.Client, "kiosks")
	var result Kiosk
	kioskCollection.FindOne(ctx, bson.M{}).Decode(&result)
	fmt.Printf("Found a kiosk: %+v\n", result)

	var dresult Doctor
	var doctorCollection *mongo.Collection = database.OpenCollection(database.Client, "doctors")
	doctorCollection.FindOne(ctx, bson.M{"kc_username": "doctor0"}).Decode(&dresult)
	fmt.Printf("Found a doctor: %+v\n", dresult)

}
