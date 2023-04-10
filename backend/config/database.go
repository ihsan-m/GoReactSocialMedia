package config

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Client {
    // Replace the following string with your MongoDB connection URI
    mongoURI := "mongodb+srv://hihasan24:<bigmistake1223>@cluster0.nedz0xn.mongodb.net/test"

    // Set client options
    clientOptions := options.Client().ApplyURI(mongoURI)

    // Connect to MongoDB
    client, err := mongo.Connect(context.TODO(), clientOptions)

    if err != nil {
        log.Fatal(err)
    }

    // Check the connection
    err = client.Ping(context.TODO(), nil)

    if err != nil {
        log.Fatal(err)
    }

    log.Println("Connected to MongoDB")
    return client
}
