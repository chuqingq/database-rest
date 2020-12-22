package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
https://github.com/ltonetwork/mongodb-rest
Listing Databases: Format: GET /dbs
Listing Collections: Format:GET /<db>/

List Documents in a Collection: Format: GET /<db>/<collection>
List documents satisfying a query: Format:GET /<db>/<collection>?query={"key":"value"}
List documents with nested queries: Format:GET /<db>/<collection>?query={"key":{"second_key":{"_id":"value"}}}
Return document by id: Format GET /<db>/<collection>/id

Inserting documents: Format: POST /<db>/<collection>
Replacing a document: Format: PUT /<db>/<collection>/id
Updating a document: Format: PATCH /<db>/<collection>/id

Deleting a document by id: Format: DELETE /<db>/<collection>/id

Bulk write (insert, update and delete) Format: POST /<db>/bulk
*/

// https://pkg.go.dev/github.com/gorilla/mux

/*
insertone/insertmany
updateone/updatemany
findone/find
deleteone/deletemany

https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#pkg-functions
*/

var client *mongo.Client

func colHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	vars := mux.Vars(r)
	opts := options.Find().SetSort(bson.D{{"age", 1}}) // TODO
	coll := client.Database(vars["db"]).Collection(vars["col"])
	switch r.Method {
	case http.MethodGet:
		// filter
		filterStr := r.FormValue("filter")
		filter := map[string]interface{}{}
		json.Unmarshal([]byte(filterStr), &filter)
		// opts
		// optsStr := r.FormValue("opts")
		// opts := map[string]interface{}{}
		// json.Unmarshal([]byte(optsStr), &opts)
		// db
		cursor, err := coll.Find(context.TODO(), filter, opts)
		if err != nil {
			log.Printf("find error: %v", err)
			// TODO
			return
		}
		var res []bson.M
		if err = cursor.All(context.TODO(), &res); err != nil {
			log.Printf("all error: %v", err)
			// TODO
			return
		}
		// TODO
		b, err := json.Marshal(res)
		if err != nil {
			// TODO
			log.Printf("marshal error: %v", err)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			log.Printf("write error: %v", err)
			// TODO
			return
		}
	case http.MethodPost:
		// docs
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		docs := []interface{}{}
		json.Unmarshal(body, &docs)
		// db
		res, err := coll.InsertMany(context.TODO(), docs /*, opts*/)
		if err != nil {
			log.Printf("insertmany error: %v", err)
			// TODO
		}
		log.Printf("insertIDS: %v", res.InsertedIDs)
		b, err := json.Marshal(res)
		// TODO
		_, err = w.Write(b)
		if err != nil {
			log.Printf("write error: %v", err)
			// TODO
			return
		}
	case http.MethodPatch:
		// res, err := coll.UpdateMany(context.TODO(), filter, update)
		// if err != nil {
		// 	log.Printf("updatemany error: %v", err)
		// 	// TODO
		// }
		// log.Printf("matchcount: %v", res.MatchedCount)
		// TODO
	case http.MethodDelete:
		// res, err := coll.DeleteMany(context.TODO(), filter, opts)
		// if err != nil {
		// 	log.Printf("deletemany error: %v", err)
		// 	// TODO
		// }
		// log.Printf("deletecount: %v", res.DeletedCount)
		// TODO
	default:
		log.Printf("unknown method: %v", r.Method)
	}
}

func main() {
	// 命令行参数
	dbUrl := flag.String("dburl", "mongodb://localhost:27017", "mongodb connection url")
	flag.Parse()
	// mongoclient
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(*dbUrl))
	if err != nil {
		log.Printf("mongo.NewClient error: %v", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Printf("client.Connect error: %v", err)
		return
	}
	// http
	r := mux.NewRouter()
	r.HandleFunc("/{db}/{col}", colHandler)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
