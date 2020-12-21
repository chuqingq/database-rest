package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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

var coll *mongo.Collection

func colHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	db := vars["db"]
	col := vars["col"]
	switch r.Method {
	case http.MethodGet:
		cursor, err := coll.Find(context.TODO(), filter, opts)
		if err != nil {
			log.Printf("find error: %v", err)
			// TODO
		}
		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			log.Printf("all error: %v", err)
		}
		// TODO
	case http.MethodPost:
		// TODO docs
		res, err := coll.InsertMany(context.TODO(), docs, opts)
		if err != nil {
			log.Printf("insertmany error: %v", err)
			// TODO
		}
		log.Printf("insertIDS: %v", res.InsertedIDs)
		// TODO
	case http.MethodPatch:
		res, err := coll.UpdateMany(context.TODO(), filter, update)
		if err != nil {
			log.Printf("updatemany error: %v", err)
			// TODO
		}
		log.Printf("matchcount: %v", res.MatchedCount)
		// TODO
	case http.MethodDelete:
		res, err := coll.DeleteMany(context.TODO(), filter, opts)
		if err != nil {
			log.Printf("deletemany error: %v", err)
			// TODO
		}
		log.Printf("deletecount: %v", res.DeletedCount)
		// TODO
	default:
		log.Printf("unknown method: %v", r.Method)
	}
}

func main() {
	// TODO 读取配置文件；或者从命令行读取
	r := mux.NewRouter()
	r.HandleFunc("/{db}/{col}", colHandler)
	http.Handle("/", r)
}
