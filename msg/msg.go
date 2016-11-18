package msg

type Message interface {
	Encode() []byte
	Decode([]byte) error
}

type MessageQueue struct {
	msgChan chan Message
}

func (mq *MessageQueue) Init() {
	mq.msgChan = make(chan Message, 2000)
}

func (mq *MessageQueue) Push(msg Message) {
	mq.msgChan <- msg
}

func (mq *MessageQueue) Pop() Message {
	msg := <-mq.msgChan
	return msg
}
