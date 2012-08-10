package main

import (
	"container/list"
	"github.com/acieroid/go-sfml"
)

var Mgr *MessageManager

type Message struct {
	text          string
	duration      int64
	remainingTime int
}

type MessageManager struct {
	messages *list.List
	bg       sfml.RectangleShape
	w, h     int
}

func InitMessages(width, height int) {
	bg := sfml.NewRectangleShape()
	bg.SetSize(float32(width), 24)
	bg.SetFillColor(sfml.FromRGBA(0, 255, 0, 150))
	bg.SetPosition(0, float32(height/2-12))
	Mgr = &MessageManager{list.New(), bg, width, height}
}

func AddMessage(message string) {
	AddMessageWithDuration(message, 150)
}

func AddMessageWithDuration(message string, msDuration int) {
	Mgr.messages.PushBack(&Message{message, int64(msDuration) * 1e6, msDuration})
}

func getMessageToDraw(delta int) (string, bool) {
	for Mgr.messages.Len() > 0 {
		el := Mgr.messages.Front()
		message := el.Value.(*Message)
		message.remainingTime -= delta
		if message.remainingTime > 0 {
			return message.text, true
		} else {
			Mgr.messages.Remove(el)
		}
	}
	return "", false
}

func DrawMessages(delta int, win sfml.RenderWindow) {
	w, h := win.Size()
	message, has := getMessageToDraw(delta)
	if has {
		win.DrawRectangleShapeDefault(Mgr.bg)
		DrawTextBig(message, int(w/2), int(h/2), true, win)
	}
}
