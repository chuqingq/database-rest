# create
curl -X POST 'http://127.0.0.1:8080/mydb/mycol' -d '[{"mykey":"myvalue1"},{"mykey":"myvalue2"}]'

# retrieve
curl -X GET 'http://127.0.0.1:8080/mydb/mycol?filter={"mykey":"myvalue1"}'

# update
curl -X PATCH 'http://127.0.0.1:8080/mydb/mycol?filter={"mykey":"myvalue1"}' -d '{"mykey":"myvalue3"}'

# delete
curl -X DELETE 'http://127.0.0.1:8080/mydb/mycol?filter={"mykey":"myvalue3"}'
