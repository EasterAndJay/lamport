package lamport

import (
  "net"
  "sync"
)

const (
  REQUEST = iota
  REPLY
  RELEASE
)

type Client struct {
  lock *sync.Mutex
  likes int32

  queue []int32
  connections []*net.Conn
  clock int32
  pid int32
  post string
}

func (c *Client) Like() {
  c.RequestLock()
  c.likes += 1
  c.ReleaseLock()
}

func (c *Client) RequestLock() {
  for conn := range connections {
    conn.Send(REQUEST)
  }
  c.lock.Lock()
}

func (c *Client) Reply(conn *net.Conn) {
  conn.Send(REPLY)
}

func (c *Client) ReleaseLock() {
  for conn := range connections {
    conn.Send(RELEASE)
  }
  c.lock.Unlock()
}

func (c *Client) RecvMsgs() {
  for conn := range connections {
    go c.ProcessMsg(conn)
  }
}

func (c *Client) ProcessMsg(conn *net.Conn) {
  conn.Recv()
  switch msg {
  case REPLY:

  case RELEASE:

  default:

  }

}

func (c *Client) UpdateClock() {
  c.clock += 1
}