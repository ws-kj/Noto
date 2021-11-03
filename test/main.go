package main

import (
    "log"
    "github.com/jroimartin/gocui"
    "strconv"
)

var (
    tui    *gocui.Gui
    shell  *gocui.View
//    oldX   int
//    oldY   int
    maxX   int
    maxY   int
    boundX int
    boundY int
    first  bool //dumb hack

    inx   int
    input string
)

func InitTui() {
    first = true
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
            cx, cy := shell.Cursor()
            if !(cy == boundY && cx <= boundX) {
                v.EditDelete(true)
            }
        case key == gocui.KeyArrowLeft:
            cx, cy := shell.Cursor()
            if !(cy == boundY && cx <= boundX) {
                v.MoveCursor(-1, 0, false)
            }
        case key == gocui.KeyArrowRight:
            v.MoveCursor(1, 0, false)
            
        case key == gocui.KeyEnter:
            process_input()
    }
}

func process_input() {  
    var input string
    lines := shell.ViewBufferLines()
    for i := boundY; i < len(lines); i++ {
        t := 0
        if i == boundY {
            t = boundX
        }
        for j := t; j < len(lines[i]); j++ {
            input += string(lines[i][j])
        }
    }
    Shprintln("\n" + input)
    prompt()
}

func prompt() {
    Shprint("Noto ?> ")
}

func VPrint(v *gocui.View, text string) {
    v.Write([]byte(text))
}

func shell_manager(g *gocui.Gui) error {
    if first { maxX, maxY = g.Size() } 
    
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
    /*
    if oldX != maxX {
        ox, oy := shell.Cursor()
        shell.SetCursor(resize_pos(ox, oy, oldX, oldY, maxX, maxY))
    }
    oldX = maxX
    oldY = maxY
    */
    if first {
        first = false
        Shprintln("Noto v0.1 -- Dimensions: " + strconv.Itoa(maxX) + ", " + strconv.Itoa(maxY))
        prompt()
   }

    return nil
}

func resize_pos(ox int, oy int, omx int, omy int, nmx int, nmy int) (int, int) {

    l := (omx * oy) + (ox+1)
    nx := l % nmx
    ny := (l - nx) / nmx

    if omx < nmx {  } //no idea why this worka
    if omx > nmx {  }  
    //Shprintln(strconv.Itoa(nx) + " " + strconv.Itoa(ny))
    return nx, ny
}

func Shprintln(t string) {
    cx, cy := shell.Cursor()
    r := []rune(t)
    for i := 0; i < len(r); i++ {
        if r[i] == '\n' { cy++; }
        if cx+len(r) >= maxX {
            if (cx+i) % maxX == 0 {
                cy++;
            }
        }
    }
    VPrint(shell, t+"\n")
    shell.SetCursor(0, cy+1)
    boundX, boundY = shell.Cursor()
}

func Shprint(t string) {
    cx, cy := shell.Cursor()
    
    r := []rune(t)
    for i := 0; i < len(r); i++ {
        if r[i] == '\n' { cy++; cx = 0 }
        if cx+len(r) >= maxX {
            if (cx+i) % maxX == 0 {
                cy++;
            }
        }
    }
   
    VPrint(shell, t)

    shell.SetCursor(cx + len(r), cy)
    boundX, boundY = shell.Cursor()
}

func quit(_ *gocui.Gui, _ *gocui.View) error {
    return gocui.ErrQuit
}

func main() {
    InitTui()
}
