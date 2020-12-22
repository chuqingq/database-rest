// https://docs.mongodb.com/drivers/node/quick-start

const { MongoClient } = require("mongodb");
// Replace the uri string with your MongoDB deployment's connection string.
const uri = "mongodb://127.0.0.1:27017/";
const client = new MongoClient(uri);
async function run() {
  try {
    await client.connect();
    const database = client.db('mydb');
    const collection = database.collection('mycol');
    // Query for a movie that has the title 'Back to the Future'
    const query = { title: 'Back to the Future' };
    const movie = await collection.findOne(query);
    console.log(movie);
  } finally {
    // Ensures that the client will close when you finish/error
    await client.close();
  }
}
run().catch(console.dir);

// chuqq

const express = require('express')
const app = express()
const port = 8080

app.post('/:db/:col', (req, res) =>
	console.log(`db: ${req.params.db}, col: ${req.params.col}`)

	res.send('Hello World!')
)

app.listen(port, () => console.log(`Example app listening on port ${port}!`))
