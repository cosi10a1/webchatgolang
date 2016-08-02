package chat
type Message struct {
	Author string `json:"author"`
	Body string `json:"body"`
}

func (message *Message) String() string{
	return message.Author+" says "+message.Body
	
}