package main

import (
	"fmt"
	"log"
	"os"

	termbox "github.com/nsf/termbox-go"
)

// DummySerialPort is a fake serial port.
type DummySerialPort struct {
	Msgs chan string
}

// NewDummySerialPort creates a new Dummy Serial Port.
func NewDummySerialPort(msgs chan string) *DummySerialPort {
	return &DummySerialPort{Msgs: msgs}
}

func (dsp DummySerialPort) Read(p []byte) (n int, err error) {
	msg := <-dsp.Msgs
	copy(p, []byte(msg))
	return len(p), nil
}

func (dsp DummySerialPort) Write(p []byte) (n int, err error) {
	dsp.Msgs <- string(p)
	return len(p), nil
}

// Close closes the dummy serial port.
func (dsp DummySerialPort) Close() error {
	close(dsp.Msgs)
	return nil
}

const (
	Program    = 0
	Red        = 1
	Green      = 2
	Blue       = 3
	Brightness = 4
)

// DummyInteractive allows the user to change the control panel values interactively.
func DummyInteractive(dummyMsgs chan string) {
	err := termbox.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer termbox.Close()

	msgVals := []int16{0, 255, 0, 0, 255}
	curIndex := 0

	log.Printf("==================================")
	log.Printf("Press:")
	log.Printf("P (Program)")
	log.Printf("R (Red)")
	log.Printf("G (Green)")
	log.Printf("B (Blue)")
	log.Printf("L (Brightness)")
	log.Printf("Then use up/down to change values. Press space to submit.")
	log.Printf("%v %v %v %v %v\n", msgVals[Program], msgVals[Red], msgVals[Green], msgVals[Blue], msgVals[Brightness])
	log.Printf("==================================")

	for {
		ev := termbox.PollEvent()

		switch ev.Ch {
		case 'p':
			curIndex = Program
		case 'r':
			curIndex = Red
		case 'g':
			curIndex = Green
		case 'b':
			curIndex = Blue
		case 'l':
			curIndex = Brightness
		case 'q':
			os.Exit(0)
		}

		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				fallthrough
			case termbox.MouseWheelUp:
				msgVals[curIndex]++
			case termbox.KeyArrowDown:
				fallthrough
			case termbox.MouseWheelDown:
				msgVals[curIndex]--
			case termbox.KeySpace:
				// Sends a new message over the websocket.
				if *spaceToSend {
					dummyMsgs <- formatDummyMessage(msgVals)
				}
			case termbox.KeyCtrlC:
				fallthrough
			case termbox.KeyCtrlZ:
				fallthrough
			case termbox.KeyCtrlD:
				fallthrough
			case termbox.KeyEsc:
				os.Exit(0)
			default:
				continue
			}
		case termbox.EventError:
			log.Fatalln(ev.Err)
		}

		msgVals[Program] = Max(Min(msgVals[Program], 4), 0)

		for i := Red; i <= Brightness; i++ {
			msgVals[i] = Max(Min(msgVals[i], 255), 0)
		}

		// We send on each change.
		if *spaceToSend == false {
			dummyMsgs <- formatDummyMessage(msgVals)
		}

		log.Printf("[%v]: %v %v %v %v %v\n",
			[]string{"Program", "Red", "Green", "Blue", "Brightness"}[curIndex],
			msgVals[Program], msgVals[Red], msgVals[Green], msgVals[Blue], msgVals[Brightness])
	}
}

func formatDummyMessage(msgVals []int16) string {
	return fmt.Sprintf("%v %v %v %v %v\n",
		msgVals[Program], msgVals[Red], msgVals[Green], msgVals[Blue], msgVals[Brightness])
}
