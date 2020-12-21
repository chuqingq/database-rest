package main

import (
	"net/http"
)

/*
Listing Databases: Format: GET /dbs
Listing Collections: Format:GET /<db>/
List Documents in a Collection: Format: GET /<db>/<collection>
List documents satisfying a query: Format:GET /<db>/<collection>?query={"key":"value"}
List documents with nested queries: Format:GET /<db>/<collection>?query={"key":{"second_key":{"_id":"value"}}}
Return document by id: Format GET /<db>/<collection>/<id>

Inserting documents: Format: POST /<db>/<collection>
Replacing a document: Format: PUT /<db>/<collection>/id
Updating a document: Format: PATCH /<db>/<collection>/id

Deleting a document by id: Format: DELETE /<db>/<collection>/id

Bulk write (insert, update and delete) Format: POST /<db>/bulk
*/

func main() {

}