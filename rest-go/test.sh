# insertmany
curl -X POST 'http://127.0.0.1:8080/api/db' -d '
{
	"db": "mydb",
	"col": "mycol",
	"func": "insertmany",
	"docs": [
		{"mykey":"myvalue1","age":2},
		{"mykey":"myvalue2","age":1},
		{"mykey":"myvalue1","age":3}
	]
}'

# find
curl -X POST 'http://127.0.0.1:8080/api/db' -d '
{
	"db": "mydb",
	"col": "mycol",
	"func": "find",
	"filter": {"mykey":"myvalue1"}
}'

# updatemany
curl -X POST 'http://127.0.0.1:8080/api/db' -d '
{
	"db": "mydb",
	"col": "mycol",
	"func": "updatemany",
	"filter": {"mykey":"myvalue1"},
	"update": {"$set": {"age":4}}
}'

# deletemany
curl -X POST 'http://127.0.0.1:8080/api/db' -d '
{
	"db": "mydb",
	"col": "mycol",
	"func": "deletemany",
	"filter": {"mykey":"myvalue2"}
}'
