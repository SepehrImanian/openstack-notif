package notifier

type Notifier interface {
	Send(text string)
}
