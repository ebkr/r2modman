package globals

import "fmt"

// WidgetQueue : Used to queue widget processing
type WidgetQueue struct {
	Action func()
}

var widgetQueue chan WidgetQueue

// NewWidgetQueue : Create a new queue item, and end the previous queue
func NewWidgetQueue() chan WidgetQueue {
	widgetQueue = make(chan WidgetQueue)
	go processor(widgetQueue)
	return widgetQueue
}

func processor(self chan WidgetQueue) {
	if self != widgetQueue {
		return
	}
	res := <-self
	fmt.Println("Adding widget")
	res.Action()
	processor(self)
}
