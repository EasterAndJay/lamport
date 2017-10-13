package lamport

// import(
// )
type Queue []Message

type Message struct {
  msgType int
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
