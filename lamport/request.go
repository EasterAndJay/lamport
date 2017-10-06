package lamport

import(
  "strconv"
  "strings"
)
type Queue []Request

type Request struct {
  pid int
  clock int
}

func (q Queue) Len() int {
  return len(q)
}

func (q Queue) Swap(i, j int) {
  q[i], q[j] = q[j], q[i]
}

func (q Queue) Less(i, j int) bool {
  if q[i].clock == q[j].clock {
    return q[i].pid < q[j].pid
  }
  return q[i].clock < q[j].clock
}

func parseRequest(msgBody string) Request {
  split := strings.Split(msgBody, ",")
  pid, _ := strconv.Atoi(split[0])
  clock, _ := strconv.Atoi(split[1])
  return Request{clock, pid}
}