package lamport

import(
  "fmt"
  "os"
  "net"
  "strconv"
  "sync"
)

const (
  REQUEST = iota
  REPLY
  RELEASE
)

const (
  CONN_HOST = "localhost"
  CONN_TYPE = "tcp"
  BASE_PORT = 5000
  CLIENT_COUNT = 5
)

type Messenger struct {
  pid int
  connLock *sync.RWMutex
  connections map[int]net.Conn

  clock int
  replyCount int

  queue []int
  likeLock chan int
}

func (m *Messenger) AcceptConnections(pid int) {
  port := strconv.Itoa(pid + BASE_PORT)
  l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+port)
  if err != nil {
    fmt.Println("Error listening:", err.Error())
    os.Exit(1)
  }
  defer l.Close()
  fmt.Println("Listening on " + CONN_HOST + ":" + port)
  for {
    if len(m.connections) == CLIENT_COUNT - 1 {
      break
    }
    conn, err := l.Accept()
    if err != nil {
      fmt.Println("Error accepting: ", err.Error())
      os.Exit(1)
    }
    buff := make([]byte, 1024)
    conn.Read(buff)
    clientPid, _ := strconv.Atoi(string(buff[:]))
    m.connLock.Lock()
    m.connections[clientPid] = conn
    m.connLock.Unlock()
  }
}

func (m *Messenger) connectionExists(peerPid int) bool {
  m.connLock.RLock()
  defer m.connLock.RUnlock()
  if _, ok := m.connections[peerPid]; ok {
    return true
  }
  return false
}

func (m *Messenger) addConnection(peerPid int, conn net.Conn) {
  m.connLock.Lock()
  m.connections[peerPid] = conn
  m.connLock.Unlock()
}

func (m *Messenger) connectToPeer(peerPid int) (net.Conn, error) {
  port := strconv.Itoa(BASE_PORT + peerPid)
  servAddr := CONN_HOST + ":" + port
  tcpAddr, err := net.ResolveTCPAddr(CONN_TYPE, servAddr)
  if err != nil {
    fmt.Println("ResolveTCPAddr failed:", err.Error())
    os.Exit(1)
  }
  return net.DialTCP(CONN_TYPE, nil, tcpAddr)
}

func (m *Messenger) writeInt(i int, conn net.Conn) (int, error) {
  return conn.Write([]byte(strconv.Itoa(i)))
}

func (m *Messenger) sendPid(pid int, conn net.Conn) {
  _, err := m.writeInt(pid, conn)
  if err != nil {
    fmt.Println("Write to server failed:", err.Error())
    os.Exit(1)
  }
}

func (m *Messenger) InitiateConnections(pid int) {
  peerPid := -1
  for {
    peerPid += 1
    if peerPid > CLIENT_COUNT - 1 {
      peerPid = 0
    }
    if peerPid == pid || m.connectionExists(peerPid) {
      continue
    }
    conn, err := m.connectToPeer(peerPid)
    if err != nil {
      fmt.Printf("Client %d: Dial failed: %v\n", pid, err.Error())
      continue
    }
    m.sendPid(pid, conn)
    m.addConnection(peerPid, conn)
  }
}

func (m *Messenger) Reply(conn net.Conn) {
  m.writeInt(REPLY, conn)
}

func (m *Messenger) Request(conn net.Conn) {
  m.writeInt(REQUEST, conn)
}

func (m *Messenger) Release(conn net.Conn) {
  m.writeInt(RELEASE, conn)
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