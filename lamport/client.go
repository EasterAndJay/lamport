package lamport

import (
  "fmt"
  "net"
  "sync"
  "time"
)

type Client struct {
  likes int
  post string
  *Messenger
}

func NewClient(pid int, post string, n int) *Client {
  return &Client{
    0,
    post,
    &Messenger{
      pid,
      &Connector{
        &sync.RWMutex{},
        make(map[int]net.Conn),
      },
      0,
      0,
      make(Queue, 0, n),
      make(chan int, 1),
    },
  }
}

func (c *Client) Run(n int) {
  fmt.Println("Running")
  go c.AcceptConnections(c.pid, n)
  go c.InitiateConnections(c.pid, n)
  for len(c.connections) != n - 1 {

  }
  fmt.Printf("All connections made from client: %d\n", c.pid)
  go c.RecvMsgs()
  for {
    fmt.Printf("client: %d | Post content: %s | LIKE count: %d\n", c.pid, c.post, c.likes)
    time.Sleep(time.Second * 5)
    c.Like()
  }
}

func (c *Client) RecvMsgs() {
  for senderPid, conn := range c.connections {
    go c.ProcessMsg(senderPid, conn)
  }
}

func (c *Client) Like() {
  c.RequestLock()
  c.likes += 1
  fmt.Printf("LIKES: %d\n", c.likes)
  c.ReleaseLock()
}

func (c *Client) RequestLock() {
  for _, conn := range c.connections {
    c.Request(conn)
  }
  <-c.likeLock
}

func (c *Client) ReleaseLock() {
  for _, conn := range c.connections {
    c.Release(conn)
  }
  c.likeLock <- 1
}
