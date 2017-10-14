package main

import(
  "flag"
  "fmt"
  "github.com/easterandjay/lamport/lamport"
)

func main() {
  n := flag.Int("n", 3, "The number of clients")
  pid := flag.Int("pid", -1, "The pid of this client")
  fmt.Printf("Spawning client with pid: %d", *pid)
  c := lamport.NewClient(*pid, "The post content HERE")
  c.Run(*n)
}