package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//
// db
//

const uri = "mongodb://localhost:27017"
const db = "mongoExample"
const collection = "students"

func ping(client *mongo.Client) {
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Successfully pinged!")
}

func insertStudent(students *mongo.Collection) {

	result, err := students.InsertOne(context.TODO(), generateStudent())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Inserted student: %s\n", result.InsertedID)
}

func viewStudent(students *mongo.Collection) {

}

func listStudents(students *mongo.Collection) {

	count, err := students.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total students: %d\n", count)
}

//
// student
//

type Address struct {
	Street string `bson:"street"`
	City   string `bson:"city"`
	State  string `bson:"state"`
}

type Student struct {
	FirstName string  `bson:"first_name"`
	LastName  string  `bson:"last_name,omitempty"`
	Address   Address `bson:"address"`
	Age       int     `bson:"age,omitempty"`
}

func randomElement(input []string) string {
	return input[len(input)-1]
}

func generateStudent() Student {
	firstNames := []string{"Ashley", "Brad", "Laura", "John", "Michael", "Alice", "Heater", "Alfred"}
	lastNames := []string{"Smith", "Johnson", "Peterson", "Simpson", "Armstrong"}
	streets := []string{"Belmont", "Main", "Lakewood", "Burnside", "Elmwood"}
	cities := []string{"Portland", "Gresham", "Vancouver"}
	states := []string{"OR", "WA", "CA"}

	streeAddress := fmt.Sprintf("%d %s", rand.Intn(1000)+1, randomElement(streets))

	return Student{
		FirstName: randomElement(firstNames),
		LastName:  randomElement(lastNames),
		Address:   Address{streeAddress, randomElement(cities), randomElement(states)},
		Age:       (rand.Intn(10) + 20),
	}

}

func main() {
	rand.Seed(time.Now().UnixNano())

	//
	// init connection
	//

	fmt.Println("connecting to:", uri)

	opts := options.Client()
	opts.SetConnectTimeout(5 * time.Second)
	opts.ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), opts)
	studentCollection := client.Database(db).Collection(collection)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	//
	// cli
	//

	cmd := flag.String("cmd", "ping", "the command to run")
	size := flag.Int("size", 1, "size of the operation")
	flag.Parse()

	switch *cmd {
	case "ping":
		ping(client)
	case "insert":
		for i := 0; i < *size; i++ {
			insertStudent(studentCollection)
		}
	case "view":
		viewStudent(studentCollection)
	case "list":
		listStudents(studentCollection)
	default:
		fmt.Println("unknown cmd:", cmd)
	}

}
