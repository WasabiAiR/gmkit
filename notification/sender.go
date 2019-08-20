package notification

// Sender is an interface that sends an email.
type Sender interface {
	// Send is a function that sends an email
	Send(msg Message) error
}

// Message is basic structure for notification system.
type Message struct {
	To         []string
	From       string
	FromPretty string
	Body       string
	BodyHTML   string
	Subject    string
}
