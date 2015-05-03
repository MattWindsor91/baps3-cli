package main

import (
	"log"
	"time"

	"github.com/UniversityRadioYork/baps3-go"
	tb "github.com/nsf/termbox-go"
	ui "github.com/wlcx/spot/termboxui"
)

func drawChannel(x, y, w, h int, svcst baps3.ServiceState) {
	ui.Printlim(x, y, tb.ColorWhite, tb.ColorDefault, svcst.Identifier, w)
	progress := int(float64(svcst.Elapsed) / float64(svcst.Duration) * 100)
	drawProgressBar(x, y+h-2, w, progress)
	ui.Print(x, y+h-1, tb.ColorWhite, tb.ColorDefault, baps3.PrettyDuration(svcst.Elapsed))
	ui.Printr(x+w, y+h-1, tb.ColorWhite, tb.ColorDefault, baps3.PrettyDuration(svcst.Duration))
}

func drawChannels(x, y, w, h int, svcstates []baps3.ServiceState) {
	chanwidth := w / len(svcstates)
	for num, state := range svcstates {
		if num > 0 {
			ui.Drawbox(x+(chanwidth*num)-1, y, 1, h, "")
		}
		drawChannel(x+(chanwidth*num), y, chanwidth-1, h, state)
	}
}

func drawProgressBar(x, y, w, percent int) {
	filled, empty, pointer := '=', '-', '>'
	progress := int(float64(w) * (float64(percent) / 100.0))
	for i := 0; i < w; i++ {
		char := empty
		switch {
		case i < progress:
			char = filled
		case i == progress:
			char = pointer
		}
		tb.SetCell(x+i, y, char, tb.ColorWhite, tb.ColorDefault)
	}
}

func redraw() {
	tb.Clear(tb.ColorWhite, tb.ColorDefault)
	w, h := tb.Size()
	ui.Drawbar(0, 0, w, tb.ColorBlue)
	ui.Print(0, 0, tb.ColorWhite, tb.ColorBlue, "BAPS3-CLI")

	teststates := []baps3.ServiceState{
		baps3.ServiceState{
			"Channel 1",
			nil,
			baps3.StPlaying,
			time.Duration(time.Second * 120),
			time.Duration(time.Second * 360),
			"",
		},

		baps3.ServiceState{
			"Channel 2",
			nil,
			baps3.StPlaying,
			time.Duration(time.Second * 120),
			time.Duration(time.Second * 360),
			"",
		},
		baps3.ServiceState{
			"Channel 3",
			nil,
			baps3.StPlaying,
			time.Duration(time.Second * 120),
			time.Duration(time.Second * 360),
			"",
		},
	}
	drawChannels(0, 1, w, h-1, teststates)

	tb.Flush()
}

func main() {
	err := tb.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer tb.Close()

	eventCh := make(chan tb.Event)
	quit := false

	go func() {
		for {
			eventCh <- tb.PollEvent()
		}
	}()

	for !quit {
		redraw()
		select {
		case event := <-eventCh:
			switch event.Type {
			case tb.EventKey:
				if event.Ch == 'q' {
					quit = true
				}
			}
		}
	}
}
