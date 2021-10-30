package main

import (
    "log"
    "github.com/jroimartin/gocui"
)

var (
    tui *gocui.Gui
    shell  *gocui.View
    oldX   int
    oldY   int
)

func InitGui() {
    g, err := gocui.NewGui(gocui.OutputNormal)
    if err != nil {
        log.Fatal(err)
    }
    defer g.Close()
    tui = g
    tui.SetManagerFunc(shell_manager)
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
        log.Fatal(err)
    }

    if err := tui.MainLoop(); err != nil && err != gocui.ErrQuit {
        log.Fatal(err)
    }
}

func shell_editor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
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

func shell_manager(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    sMinX := 0
    sMinY := 0
    sMaxY := maxY-1
    sMaxX := maxX-1

    var err error
    if shell, err = g.SetView("shell", sMinX, sMinY, sMaxX, sMaxY); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        } 
        g.SetCurrentView("shell")
        shell.Editor = gocui.EditorFunc(shell_editor)
        shell.Autoscroll = false
        shell.Wrap = true
        tui.Cursor = true
        shell.Editable = true
    }
    if oldX != maxX {
        ox, oy := shell.Cursor()
        l := (oy * oldX) + ox
        nx := l % maxX
        if oldX < maxX { nx -= 2 }  //why??
        if oldX > maxX { nx += 2 }  //why does this work? i do not know
        shell.SetCursor(nx, ((l - (l % maxX)) / maxX)) 
    }
    oldX = maxX
    oldY = maxY
    return nil
}

func pshell(t string) {
    set_text(shell, t)
}

func quit(_ *gocui.Gui, _ *gocui.View) error {
    return gocui.ErrQuit
}

func main() {
    InitGui()
}
