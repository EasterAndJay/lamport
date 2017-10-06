package lamport

import(
  "fmt"
  "net"
  "strconv"
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

  queue []int
  likeLock chan int
}



func writeInt(i int, conn net.Conn) (int, error) {
  return conn.Write([]byte(strconv.Itoa(i)))
}

func (m *Messenger) Reply(conn net.Conn) {
  writeInt(REPLY, conn)
}

func (m *Messenger) Request(conn net.Conn) {
  writeInt(REQUEST, conn)
}

func (m *Messenger) Release(conn net.Conn) {
  writeInt(RELEASE, conn)
}

func (m *Messenger) ProcessMsg(senderPid int, conn net.Conn) {
  buff := make([]byte, 1024)
  for {
    conn.Read(buff)
    fmt.Printf("Client %d: Message received from Client %d\n", m.pid, senderPid)
    m.UpdateClock()
    msgType, _ := strconv.Atoi(string(buff[:]))
    switch msgType {
    case REQUEST:
      m.Reply(conn)
    case REPLY:
      m.replyCount += 1
      if m.replyCount == len(m.connections) && m.queue[0] == m.pid {
        m.replyCount = 0
        m.likeLock <- 1
      }
    case RELEASE:
      m.queue = m.queue[1:]
      if m.replyCount == len(m.connections) && m.queue[0] == m.pid {
        m.replyCount = 0
        m.likeLock <- 1
      }
    default:

    }
  }

}

func (m *Messenger) UpdateClock() {
  m.clock += 1
  fmt.Printf("Client %d: Updated clock to %d\n", m.pid, m.clock)
}