import express, { application } from "express";

const app = express();
const port = 3000;

app.get("/", (req, res) => {
  res.send("hello world");
});

// Start the server
app.listen(port, () => {
  console.log(`Server is running on port ${port}`);
});
