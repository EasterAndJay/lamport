package lamport

import(
  "fmt"
  "os"
  "net"
  "strconv"
  "sync"
)

const (
  CONN_HOST = "localhost"
  CONN_TYPE = "tcp"
  BASE_PORT = 5000
)

type Connector struct {
  connLock *sync.RWMutex
  connections map[int]net.Conn
}


func (cn *Connector) AcceptConnections(pid int) {
  port := strconv.Itoa(pid + BASE_PORT)
  l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+port)
  if err != nil {
    fmt.Println("Error listening:", err.Error())
    os.Exit(1)
  }
  defer l.Close()
  fmt.Println("Listening on " + CONN_HOST + ":" + port)
  for {
    if len(cn.connections) == CLIENT_COUNT - 1 {
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
    cn.connLock.Lock()
    cn.connections[clientPid] = conn
    cn.connLock.Unlock()
  }
}

func (cn *Connector) InitiateConnections(pid int) {
  peerPid := -1
  for {
    peerPid += 1
    if peerPid > CLIENT_COUNT - 1 {
      peerPid = 0
    }
    if peerPid == pid || cn.connectionExists(peerPid) {
      continue
    }
    conn, err := cn.connectToPeer(peerPid)
    if err != nil {
      fmt.Printf("Client %d: Dial failed: %v\n", pid, err.Error())
      continue
    }
    sendPid(pid, conn)
    cn.addConnection(peerPid, conn)
  }
}

func (cn *Connector) connectionExists(peerPid int) bool {
  cn.connLock.RLock()
  defer cn.connLock.RUnlock()
  if _, ok := cn.connections[peerPid]; ok {
    return true
  }
  return false
}

func (cn *Connector) addConnection(peerPid int, conn net.Conn) {
  cn.connLock.Lock()
  cn.connections[peerPid] = conn
  cn.connLock.Unlock()
}

func (cn *Connector) connectToPeer(peerPid int) (net.Conn, error) {
  port := strconv.Itoa(BASE_PORT + peerPid)
  servAddr := CONN_HOST + ":" + port
  tcpAddr, err := net.ResolveTCPAddr(CONN_TYPE, servAddr)
  if err != nil {
    fmt.Println("ResolveTCPAddr failed:", err.Error())
    os.Exit(1)
  }
  return net.DialTCP(CONN_TYPE, nil, tcpAddr)
}

func sendPid(pid int, conn net.Conn) {
  _, err := writeInt(pid, conn)
  if err != nil {
    fmt.Println("Write to server failed:", err.Error())
    os.Exit(1)
  }
}