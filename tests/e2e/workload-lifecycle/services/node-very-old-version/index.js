// Import the express module
var express = require('express');

// Create an instance of the express application
var app = express();

// Define a route for the root URL that sends "Hello, World!" as the response
app.get('/', function (req, res) {
  res.send('Hello, World!');
});

// Start the server and listen on port 3000
app.listen(3000, function () {
  console.log('Server running at http://127.0.0.1:3000/');
});
