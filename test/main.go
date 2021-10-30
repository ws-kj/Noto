package main

import (
    "log"
//    "strconv"
    "github.com/jroimartin/gocui"
)

type Tui struct {
    Win    *gocui.Gui
    Disp   *gocui.View    
    Shell  *gocui.View
    Dbuf   string
    Sbuf   string
    Vert   bool
    oldX   int
    oldY   int
}

func (tui *Tui) init() {
    g, err := gocui.NewGui(gocui.OutputNormal)
    if err != nil {
        log.Fatal(err)
    }
    defer g.Close()
    tui.Win = g
    tui.Win.SetManagerFunc(tui.shell_manager)
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, tui.quit); err != nil {
        log.Fatal(err)
    }

    if err := tui.Win.MainLoop(); err != nil && err != gocui.ErrQuit {
        log.Fatal(err)
    }

}

func editor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
    switch {
        case ch != 0 && mod == 0:
            v.EditWrite(ch)
        case key == gocui.KeySpace:
            v.EditWrite(' ')
        case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
            v.EditDelete(true)
        case key == gocui.KeyArrowLeft:
            v.MoveCursor(-1, 0, false)
        case key == gocui.KeyArrowRight:
            v.MoveCursor(1, 0, false)
    }
}

func set_text(v *gocui.View, text string) {
    v.Clear()
    v.Write([]byte(text))
    v.SetCursor(len(text), 0)
}

func (tui *Tui) shell_manager(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    sMinX := 0
    sMinY := 0
    sMaxY := maxY-1
    sMaxX := maxX-1

    var err error
    if tui.Shell, err = g.SetView("shell", sMinX, sMinY, sMaxX, sMaxY); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        } 
        g.SetCurrentView("shell")
        tui.Shell.Editor = gocui.EditorFunc(editor)
        tui.Shell.Autoscroll = false
        tui.Shell.Wrap = true
        tui.Win.Cursor = true
        tui.Shell.Editable = true
    }
    if tui.oldX != maxX {
        ox, oy := tui.Shell.Cursor()
        l := (oy * tui.oldX) + ox
        nx := l % maxX
        if tui.oldX < maxX { nx -= 2 }  //why??
        if tui.oldX > maxX { nx += 2 }  //why does this work? i do not know
        tui.Shell.SetCursor(nx, ((l - (l % maxX)) / maxX)) 
    }
    tui.oldX = maxX
    tui.oldY = maxY
    return nil
}

func (tui *Tui) pshell(t string) {
    set_text(tui.Shell, t)
}


func (tui *Tui) quit(_ *gocui.Gui, _ *gocui.View) error {
    return gocui.ErrQuit
}

func main() {
    tui := Tui{}
    tui.init()
}
