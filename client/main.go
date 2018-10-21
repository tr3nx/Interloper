package main

import (
	"bytes"
	"fmt"
	"github.com/jroimartin/gocui"
	"io"
	"log"
	"net"
)

var (
	keychain = []byte("PaJJM8JjLmMnhR6Emfexi3Oa37bYDQph")
	done     = make(chan bool)
	incoming = make(chan string)
	outgoing = make(chan string)
	host     = "localhost"
	port     = "1337"
	conn     net.Conn
)

func sendMessage(g *gocui.Gui, v *gocui.View) error {
	go func() {
		buf := v.ViewBuffer()
		v.Clear()
		v.SetCursor(0, 0)
		data, err := encrypt([]byte(buf), keychain)
		if err != nil {
			return
		}
		message := packMessage(data)
		if len(message) > 0 {
			b := bytes.NewReader(message)
			io.Copy(conn, b)
		}
	}()
	return nil
}

func receiveMessages(conn net.Conn) {
	for {
		b := unpackMessage(conn)
		message, _ := decrypt(b, keychain)
		if len(message) > 0 {
			incoming <- string(message)
		}
	}
}

func renderMessages(g *gocui.Gui) {
	for {
		select {
		case <-done:
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("messages")
				if err != nil {
					return err
				}
				fmt.Fprintf(v, "Connection closed...")
				return nil
			})
			return

		case message := <-incoming:
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("messages")
				if err != nil {
					return err
				}
				fmt.Fprint(v, message)
				return nil
			})
		}
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if mv, err := g.SetView("messages", 0, 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		mv.Title = "Messages"
		mv.Autoscroll = true
		mv.Wrap = true
	}

	if sv, err := g.SetView("send", 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		sv.Title = "Send"
		sv.Editable = true
		if _, err := g.SetCurrentView("send"); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	return gocui.ErrQuit
}

func main() {
	var err error
	conn, err = net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		panic(err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, sendMessage); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	go renderMessages(g)
	go receiveMessages(conn)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
