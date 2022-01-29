package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func insertStudents(students *mongo.Collection, num *int) {
	studentList := []interface{}{}

	for i := 0; i < *num; i++ {
		studentList = append(studentList, generateStudent())
	}

	result, err := students.InsertMany(context.TODO(), studentList)
	if err != nil {
		panic(err)
	}

	ids := result.InsertedIDs
	for _, id := range ids {
		fmt.Printf("Inserted document with _id: %v\n", id)
	}
	fmt.Printf("Documents inserted: %v\n", len(ids))
}

func viewStudent(students *mongo.Collection, id *string) {
	var result Student
	objectId, err := primitive.ObjectIDFromHex(*id)
	if err != nil {
		panic(err)
	}
	err = students.FindOne(context.TODO(), bson.D{{"_id", objectId}}, options.FindOne()).Decode(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

func deleteStudent(students *mongo.Collection, id *string) {
	objectId, convErr := primitive.ObjectIDFromHex(*id)
	if convErr != nil {
		panic(convErr)
	}
	result, err := students.DeleteOne(context.TODO(), bson.D{{"_id", objectId}})
	if err != nil {
		panic(err)
	}
	if result.DeletedCount == 1 {
		fmt.Printf("confirmed deleted: %s\n", *id)
	} else {
		fmt.Printf("object not deleted or does not exist: %s\n", *id)
	}
}

func countStudents(students *mongo.Collection) {

	count, err := students.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total students: %d\n", count)
}

func listStudents(students *mongo.Collection, limit *int64, offset *int64) {
	options := options.Find().SetLimit(*limit).SetSkip(*offset)
	cursor, err := students.Find(context.TODO(), bson.D{}, options)
	if err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		var result Student

		if err := cursor.Decode(&result); err != nil {
			panic(err)
		}

		fmt.Println(result)
	}
	if err := cursor.Err(); err != nil {
		panic(err)
	}
}

func filterStudents(students *mongo.Collection, limit *int64, offset *int64, sort *string, dir *string, filter *string, term *string) {
	var direction int
	if strings.ToUpper(*dir) == "ASC" {
		direction = 1
	} else {
		direction = -1
	}

	options := options.Find().SetLimit(*limit).SetSkip(*offset).SetSort(bson.D{{*sort, direction}})
	cursor, err := students.Find(context.TODO(), bson.D{{*filter, *term}}, options)
	if err != nil {
		panic(err)
	}

	for cursor.Next(context.TODO()) {
		var result Student

		if err := cursor.Decode(&result); err != nil {
			panic(err)
		}

		fmt.Println(result)
	}
	if err := cursor.Err(); err != nil {
		panic(err)
	}
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
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FirstName string             `bson:"first_name"`
	LastName  string             `bson:"last_name,omitempty"`
	Address   Address            `bson:"address"`
	Age       int                `bson:"age,omitempty"`
}

func randomElement(input []string) string {
	return input[rand.Intn(len(input)-1)]
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
	num := flag.Int("num", 5, "number of students to generate")
	id := flag.String("id", "-", "an object id")
	limit := flag.Int64("limit", 5, "max results returned")
	offset := flag.Int64("offset", 0, "offset to begin return results")
	sort := flag.String("sort", "last_name", "field on which to sort")
	dir := flag.String("dir", "asc", "direction to sort results: 'asc' or 'desc'")
	filter := flag.String("filter", "first_name", "field on which to filter")
	term := flag.String("term", "-", "value to filter for")
	flag.Parse()

	switch *cmd {
	case "ping":
		ping(client)
	case "count":
		countStudents(studentCollection)
	case "insert":
		insertStudent(studentCollection)
	case "insertMany":
		insertStudents(studentCollection, num)
	case "view":
		viewStudent(studentCollection, id)
	case "list":
		listStudents(studentCollection, limit, offset)
	case "delete":
		deleteStudent(studentCollection, id)
	case "filter":
		filterStudents(studentCollection, limit, offset, sort, dir, filter, term)
	default:
		fmt.Printf("unknown cmd: %s\n", *cmd)
	}

}
