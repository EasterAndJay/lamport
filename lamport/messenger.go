package lamport

import(
  "encoding/gob"
  "fmt"
  "math"
  "net"
  "sort"
)

const (
  REQUEST = iota
  REPLY
  RELEASE
)

type Messenger struct {
  pid int
  *Connector
  clock int
  replyCount int

  queue Queue
  likeLock chan int
}

func (m *Messenger) SendMessage(msg Message, conn net.Conn) error {
  encoder := gob.NewEncoder(conn)
  m.UpdateClock(-1)
  err := encoder.Encode(msg)
  if err != nil {
    fmt.Println(err)
  }
  return err
}

func (m *Messenger) RecvMessage(conn net.Conn) (Message, error) {
  msg := Message{}
  decoder := gob.NewDecoder(conn)
  err := decoder.Decode(&msg)
  if err != nil {
    fmt.Println(err)
    return msg, err
  }
  m.UpdateClock(msg.Clock)
  return msg, nil
}

func (m *Messenger) Reply(conn net.Conn) {
  msg := Message{REPLY, m.pid, m.clock}
  m.SendMessage(msg, conn)
}

func (m *Messenger) Request(conn net.Conn) {
  msg := Message{REQUEST, m.clock, m.pid}
  m.Enqueue(msg)
  m.SendMessage(msg, conn)
}

func (m *Messenger) Release(conn net.Conn) {
  msg := Message{RELEASE, m.clock, m.pid}
  m.SendMessage(msg, conn)
}

func (m *Messenger) Enqueue(msg Message) {
  m.queue = append(m.queue, msg)
  sort.Sort(m.queue)
}

func (m *Messenger) UpdateClock(peerClock int) {
  m.clock = int(math.Max(float64(m.clock), float64(peerClock))) + 1
  fmt.Printf("Client %d: Updated clock to %d\n", m.pid, m.clock)
}

func (m *Messenger) ProcessMsg(senderPid int, conn net.Conn) {
  for {
    msg, err := m.RecvMessage(conn)
    if err != nil {
      panic(err)
    }
    switch msg.MsgType {
    case REQUEST:
      fmt.Printf("Client %d: Request Message received from Client %d\n", m.pid, senderPid)
      m.Enqueue(msg)
      m.Reply(conn)
    case REPLY:
      fmt.Printf("Client %d: Reply Message received from Client %d\n", m.pid, senderPid)
      m.replyCount += 1
      if m.replyCount == len(m.connections) && m.queue[0].Pid == m.pid {
        m.replyCount = 0
        m.likeLock <- 1
      }
    case RELEASE:
      fmt.Printf("Client %d: Release Message received from Client %d\n", m.pid, senderPid)
      m.queue = m.queue[1:]
      if m.replyCount == len(m.connections) && m.queue[0].Pid == m.pid {
        m.replyCount = 0
        m.likeLock <- 1
      }
    default:

    }
  }
}
