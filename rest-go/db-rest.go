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
/*
{
	"func": "insertmany|find|updatemany|deletemany",
	"filter": {
	},
	"docs": [
	],
	"opts": {"Limit":2,"Sort":{"mykey":1}}
}
*/

type body struct {
	Func string `json:func`
	Docs []interface{} `json:docs`
	Filter interface{} `json:filter`
	Update interface{} `json:update`
	Opts json.RawMessage `json:opts`
}

func colHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
    w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
    w.Header().Set("content-type", "application/json")             //返回数据格式是json
	if r.Method != http.MethodPost {
		log.Printf("method not post")
		return
	}
	// r.ParseForm()
	vars := mux.Vars(r)
	// opts := options.Find().SetSort(bson.D{{"age", 1}}) // TODO
	coll := client.Database(vars["db"]).Collection(vars["col"])
	// body
	defer r.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read body error: %v", err)
		return
	}
	log.Printf("bodyBytes: %v", string(bodyBytes))
	var body body
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		log.Printf("json.Unmarshal error: %v", err)
		return
	}
	log.Printf("body: %v", body)
	// -----
	switch body.Func {
	case "insertmany":
		// opts
		opts := options.InsertMany()
		err := json.Unmarshal(body.Opts, opts)
		if err != nil {
			log.Printf("unmarshal opts error: %v", err)
			// return
		}
		// db
		res, err := coll.InsertMany(context.TODO(), body.Docs, opts)
		if err != nil {
			log.Printf("insertmany error: %v", err)
			return
		}
		log.Printf("insertIDS: %v", res.InsertedIDs)
		b, err := json.Marshal(res)
		if err != nil {
			log.Printf("marshal error: %v", err)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			log.Printf("write error: %v", err)
			return
		}
	case "find":
		// opts
		opts := &options.FindOptions{}
		err := json.Unmarshal(body.Opts, opts)
		if err != nil {
			log.Printf("unmsharl opts error: %v", err)
			// return
		}
		log.Printf("opts: %v", opts)
		// db
		cursor, err := coll.Find(context.TODO(), body.Filter, opts)
		if err != nil {
			log.Printf("find error: %v", err)
			return
		}
		var res []bson.M
		if err = cursor.All(context.TODO(), &res); err != nil {
			log.Printf("all error: %v", err)
			return
		}
		// res
		b, err := json.Marshal(res)
		if err != nil {
			log.Printf("marshal error: %v", err)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			log.Printf("write error: %v", err)
			return
		}
	case "updatemany":
		res, err := coll.UpdateMany(context.TODO(), body.Filter, body.Update)
		if err != nil {
			log.Printf("updatemany error: %v", err)
			return
		}
		log.Printf("matchcount: %v", res.MatchedCount)
		b, err := json.Marshal(res)
		if err != nil {
			log.Printf("marshal error: %v", err)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			log.Printf("write error: %v", err)
			return
		}
	case "deletemany":
		// opts
		opts := &options.DeleteOptions{}
		err := json.Unmarshal(body.Opts, opts)
		if err != nil {
			log.Printf("unmsharl opts error: %v", err)
			// return
		}
		log.Printf("opts: %v", opts)
		// db
		res, err := coll.DeleteMany(context.TODO(), body.Filter, opts)
		if err != nil {
			log.Printf("deletemany error: %v", err)
			return
		}
		log.Printf("deletecount: %v", res.DeletedCount)
		// res
		b, err := json.Marshal(res)
		if err != nil {
			log.Printf("marshal error: %v", err)
			return
		}
		_, err = w.Write(b)
		if err != nil {
			log.Printf("write error: %v", err)
			return
		}
	default:
		log.Printf("unknown method: %v", r.Method)
	}
}

func main() {
	// 命令行参数
	dbUrl := flag.String("dburl", "mongodb://localhost:27017", "mongodb connection url")
	// TODO url prefix
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
