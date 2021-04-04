package main

import (
	"encoding/json"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	kafka "github.com/segmentio/kafka-go"
	"github.com/joho/godotenv"
)


// Database için istek struct
type ResponseLog struct {
	_id		  string
	ReqType   string
	RespTime  int64
	Timestamp int64
}
  
var (
	RequestLogger *log.Logger
	kafkaWriter *kafka.Writer
	collection *mongo.Collection
)


// Functions

func getRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET Request -> OK")
}

func postRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST Request -> OK")
}

func putRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT Request -> OK")
}

func deleteRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE Request -> OK")
}


func chartData(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.D{{"timestamp", bson.D{{"$gt", time.Now().Add(time.Hour*-1).Unix()}}}})
	if err != nil { log.Fatal(err) }
	defer cur.Close(ctx)
	if err := cur.Err(); err != nil {
	log.Fatal(err)
	}
	var logsFiltered []bson.M
	if err = cur.All(ctx, &logsFiltered); err != nil {
		log.Fatal(err)
	}
	fmt.Println(logsFiltered)
	b, _ := json.Marshal(logsFiltered)
	w.Write(b)
}

func logHandler(fn http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()
		time.Sleep(time.Duration(rand.Intn(3000))*time.Millisecond)
		fn(w,r)
		respTime := (time.Now().UnixNano()-timeStart.UnixNano())/1000000
        //Logs.txt adında dosya yoksa oluştur varsa içine gir.
		file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		
		logString := fmt.Sprintf("%s, %d, %d", r.Method, respTime, timeStart.Unix())
		RequestLogger = log.New(file, logString , 0)
		RequestLogger.Println()

		responseLog := ResponseLog{
			ReqType : r.Method,
			RespTime : respTime,
			Timestamp : timeStart.Unix(),
		}
		b, _ := json.Marshal(responseLog)
		msg := kafka.Message{
			Value: []byte(b),
		}
		err = kafkaWriter.WriteMessages(context.Background(), msg)
		if err != nil {
			fmt.Println(err)
		}
    }
}

func getMongoCollection(mongoURL, dbName, collectionName string) *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB ... !!")

	db := client.Database(dbName)
	collection := db.Collection(collectionName)
	return collection
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaURL},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func kafkaRead(){
	err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }
	// get Mongo db Collection using environment variables.
	mongoURL := os.Getenv("mongoURL")
	dbName := os.Getenv("dbName")
	collectionName := os.Getenv("collectionName")
	collection = getMongoCollection(mongoURL, dbName, collectionName)

	// get kafka reader using environment variables.
	kafkaURL := os.Getenv("kafkaURL")
	topic := os.Getenv("topic")
	groupID := os.Getenv("groupID")
	reader := getKafkaReader(kafkaURL, topic, groupID)
	kafkaWriter = newKafkaWriter(kafkaURL, topic)

	defer reader.Close()

	fmt.Println("start consuming ... !!")

	for {
		var responseLog ResponseLog
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(msg.Value, &responseLog)
		insertResult, err := collection.InsertOne(context.Background(), responseLog)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	//Handlings
	myRouter.HandleFunc("/", logHandler(getRequest)).Methods("GET")
	myRouter.HandleFunc("/", logHandler(postRequest)).Methods("POST")
	myRouter.HandleFunc("/", logHandler(putRequest)).Methods("PUT")
	myRouter.HandleFunc("/", logHandler(deleteRequest)).Methods("DELETE")
	myRouter.HandleFunc("/getChart", chartData).Methods("GET")
	fileServer := http.FileServer(http.Dir("./public/"))
	myRouter.Handle("/public/", http.StripPrefix("/public", fileServer))
	corsObj:=handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":10000", handlers.CORS(corsObj)(myRouter)))
}

func main() {
	go kafkaRead()
	handleRequests()
}
