"use strict";

const PORT = process.env.PORT || "8080";
const express = require("express");
const app = express();

app.get("/", (req, res) => {
  res.send("Hello world!");
});

app.listen(parseInt(PORT, 10), () => {
  console.log(`Listening for requests on http://localhost:${PORT}`);
});
