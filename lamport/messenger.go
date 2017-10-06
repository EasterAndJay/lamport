package lamport

import(
  "fmt"
  "net"
  "sort"
  "strconv"
  "strings"
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

func writeInt(i int, conn net.Conn) (int, error) {
  return conn.Write([]byte(strconv.Itoa(i)))
}

func (m *Messenger) Reply(conn net.Conn) {
  writeInt(REPLY, conn)
}

func (m *Messenger) Request(conn net.Conn) {
  m.Enqueue(Request{m.clock, m.pid})
  msg := fmt.Sprintf("%d,%d,%d", REQUEST, m.pid, m.clock)
  conn.Write([]byte(msg))
}

func (m *Messenger) Release(conn net.Conn) {
  writeInt(RELEASE, conn)
}

func (m *Messenger) Enqueue(r Request) {
  m.queue = append(m.queue, r)
  sort.Sort(m.queue)
}

func (m *Messenger) UpdateClock() {
  m.clock += 1
  fmt.Printf("Client %d: Updated clock to %d\n", m.pid, m.clock)
}

func (m *Messenger) ProcessMsg(senderPid int, conn net.Conn) {
  buff := make([]byte, 1024)
  for {
    conn.Read(buff)
    fmt.Printf("Client %d: Message received from Client %d\n", m.pid, senderPid)
    m.UpdateClock()
    msg := string(buff[:])
    // TODO: BUGS BE HERE
    fmt.Println(msg)
    split := strings.SplitN(msg, ",", 1)
    msgType, _ := strconv.Atoi(split[0])
    switch msgType {
    case REQUEST:
      fmt.Println(split)
      msgBody := split[1]
      request := parseRequest(msgBody)
      m.Enqueue(request)
      m.Reply(conn)
    case REPLY:
      m.replyCount += 1
      if m.replyCount == len(m.connections) && m.queue[0].pid == m.pid {
        m.replyCount = 0
        m.likeLock <- 1
      }
    case RELEASE:
      m.queue = m.queue[1:]
      if m.replyCount == len(m.connections) && m.queue[0].pid == m.pid {
        m.replyCount = 0
        m.likeLock <- 1
      }
    default:

    }
  }
}
