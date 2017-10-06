package main

import(
  "fmt"
  "github.com/easterandjay/lamport/lamport"
)

func main() {
  fmt.Println("Spawning 5 clients")
  var c *lamport.Client
  for i := 0; i < 5; i++ {
    c = lamport.NewClient(i, "The post content HERE")
    go c.Run()
  }
  for {

  }
}