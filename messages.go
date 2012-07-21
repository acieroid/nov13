package main

import (
	"container/list"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
)

var Mgr *MessageManager

type Message struct {
	text     string
	duration int64
	remainingTime int
}

type MessageManager struct {
	messages  *list.List
	bg        *sdl.Surface
	w, h int
}

func InitMessages(width, height int) {
	Mgr = &MessageManager{list.New(), nil, width, height}
	Mgr.bg = sdl.CreateRGBSurface(sdl.HWSURFACE,
		width, 24, 32, 0, 0, 0, 0)
	Mgr.bg.FillRect(&sdl.Rect{0, 0, uint16(width), 24}, 0x0000FF00)
	Mgr.bg.SetAlpha(sdl.SRCALPHA, 150)
}

func AddMessage(message string) {
	AddMessageWithDuration(message, 500)
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


func DrawMessages(delta int, surf *sdl.Surface) {
	message, has := getMessageToDraw(delta)
	if has {
		surf.Blit(&sdl.Rect{0, int16(Mgr.h/2 - 12), 0, 0},
			Mgr.bg, nil)
		DrawTextBig(message, Mgr.w/2 - 100, Mgr.h/2 - 10, surf)
	}
}
