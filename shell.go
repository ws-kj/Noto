package main

import (
    "fmt"
    "os"
    "path/filepath"
    "log"
    "strconv"
    "strings"
    "bufio"
)

type Color string
const (
    Reset  Color = "\033[0m"
    Red    Color = "\033[31m"
    Green  Color = "\033[32m"
    Yellow Color = "\033[33m"
    Blue   Color = "\033[34m"
    Purple Color = "\033[35m"
    Cyan   Color = "\033[36m"
    White  Color = "\033[37m"
)

type Shell struct {
    CurP     *Project
    CurS     *Source
    CurN     *Note
    Prompt   string
}

func Cprintln(c Color, s string) {
    fmt.Println(string(c), s, string(Reset))
}

func errlog(s string) {
    Cprintln(Red, "\n\t"+s+"\n")
}

func (shell *Shell) build() string {
    if shell.CurP == nil {
        return "Noto ?> ";
    }
    r := shell.CurP.Title;
    if shell.CurS == nil {
        return r + " ?> ";
    }
    r += " | "
    if shell.CurS.Title == "" {
        r += "Source "
        r += strconv.Itoa(shell.CurS.Id)
    } else {
        r += shell.CurS.Title
    }
    if shell.CurN == nil {
        r += " ?> "
        return r
    }
    r += " | "
    if shell.CurN.Title == "" {
        r += "Note "
        r += strconv.Itoa(shell.CurN.Id)
    } else {
        r += shell.CurN.Title
    }
    r += " ?> "
    return r
}

func (shell *Shell) prompt() {
    p := shell.build()
    fmt.Print(p)
    reader := bufio.NewReader(os.Stdin)
    in, _ := reader.ReadString('\n')
    shell.process(in)
    shell.prompt()
}

func (shell *Shell) process(input string) {
    args := strings.Fields(input)

    if len(args) == 0 {
        return
    }

    switch args[0] {
        case "list", "l":
            shell.c_list(args)
        case "new", "n":
            shell.c_new(args)
        case "use", "u":
            shell.c_use(args)
        case "edit", "e":
            shell.c_edit(args)
        case "cite", "c":
            shell.c_cite(args)
        case "info", "i":
            shell.c_info(args)
        case "back", "b", "..":
            shell.c_back(args)
        default:
            errlog("Command not recognized")
    }
}

func (shell *Shell) c_list(args []string) {
    if len(args) > 1 {
        switch args[1] {
            case "projects","project","p":
                shell.listp(args)
            case "sources","source","s":
                if shell.CurP != nil {
                    shell.lists(args)
                } else {
                    errlog("Use a project first")
                }
            case "notes","note","n":
                if shell.CurS != nil {
                    shell.listn(args)
                } else {
                    errlog("Use a source first")
                }
            default:
                errlog("Command not recognized")
        }
    } else {
        if shell.CurP == nil {
            shell.listp(args)
        } else if shell.CurS == nil {
            shell.lists(args)
        } else {
            shell.listn(args)
        }
    }
}

//DRY!!
func (shell *Shell) listp(args []string) {
    f, err := os.Open(filepath.Join(home(), ".noto"))
    if err != nil {
        log.Fatal(err)
    }
    files, err := f.Readdir(0)
    if err != nil {
        log.Fatal(err)
    }
    for _, proj := range files {
        if proj.IsDir() {
            fmt.Println(proj.Name())
        }
    }
}
func (shell *Shell) lists(args []string) {
    f, err := os.Open(filepath.Join(home(), ".noto", shell.CurP.Title))
    if err != nil {
        log.Fatal(err)
    }
    files, err := f.Readdir(0)
    if err != nil {
        log.Fatal(err)
    }
    if len(files) < 2 {
        errlog("No sources yet! Use `new source` to create one.")
        return
    }

    var entries = make([]int,len(files)-1)
    i := 0
    for _, sf := range files {
        if sf.IsDir() {
            n, _ := strconv.Atoi(sf.Name())
            entries[i] = n
            i++
        }
    }

    for i := 0; i < len(entries)-1; i++ {
        for j := 0; j < len(entries)-1; j++ {
            if entries[j] > entries[j+1] {
                tmp := entries[j+1]
                entries[j+1] = entries[j]
                entries[j] = tmp
            }
        }
    }
    if len(entries) > 0 {
        fmt.Println()
    }
    for _, n := range entries {
        entry := "\tSource "
        entry += strconv.Itoa(n)
        entry += " - "
        src := shell.CurP.load_source(n)
        if src.Title != "" {
            entry += src.Title
        } else {
            entry += "No Title"
        }
        fmt.Println(entry)
    }
    if len(entries) > 0 {
        fmt.Println()
    }
}

func (shell *Shell) listn(args []string) {
    f, err := os.Open(filepath.Join(home(), ".noto", shell.CurP.Title, strconv.Itoa(shell.CurS.Id)))
    if err != nil {
        log.Fatal(err)
    }
    files, err := f.Readdir(0)
    if err != nil {
        log.Fatal(err)
    }
    if len(files) < 2 {
        errlog("No notes yet! Use `new note` to create one.")
        return
    }

    var entries = make([]int,len(files)-1)
    i := 0
    for _, sf := range files {
        if sf.Name() != "source.json" {
            n, _ := strconv.Atoi(strings.Split(sf.Name(), ".")[0])
            entries[i] = n
            i++
        }
    } 
    for i := 0; i < len(entries)-1; i++ {
        for j := 0; j < len(entries)-1; j++ {
            if entries[j] > entries[j+1] {
                tmp := entries[j+1]
                entries[j+1] = entries[j]
                entries[j] = tmp
            }
        }
    }

    if len(entries) > 0 {
        fmt.Println()
    }
    for _, n := range entries {
        entry := "\tNote "
        entry += strconv.Itoa(n)
        entry += " - "
        note := shell.CurS.load_note(n)
        if note.Title != "" {
            entry += note.Title
        } else {
            entry += "No Title"
        }
        fmt.Println(entry)
    }
    if len(entries) > 0 {
        fmt.Println()
    }
}

func (shell *Shell) c_new(args []string) {
    if len(args) < 2 {
        errlog("Specify what to create")
        return
    }
    switch args[1] {
        case "project","p":
            if len(args) < 3 {
                errlog("Specify project name")
                return
            }
            var name string
            for i := 2; i < len(args); i++ {
                name += args[i] + " "
            }
            name = strings.TrimRight(name, "\n")
            name = strings.TrimSpace(name)
            p := create_project(name)
            if p != nil {
                shell.CurP = p
                shell.CurS = nil
                shell.CurN = nil
            }
        case "source","s":
            if shell.CurP == nil {
                errlog("Use a project first")
                return
            }
            s := shell.CurP.create_source()
            shell.CurS = s
        case "note","n":
            if shell.CurS == nil {
                errlog("Use a source first")
                return
            }
            n := shell.CurS.create_note()
            shell.CurN = n
        default:
            errlog("Command not recognized")
    }
}

func (shell *Shell) c_use(args []string) {
    if len(args) > 1 {
        switch args[1] {
            case "project", "p":
                if len(args) < 3 {
                    errlog("Specify a project to use")
                    return
                }
                var name string
                for i := 2; i < len(args); i++ {
                    name += args[i] + " "
                }
                name = strings.TrimRight(name, "\n")
                name = strings.TrimSpace(name)
                p := load_project(name)
                if p != nil {
                    shell.CurP = p
                    shell.CurS = nil
                    shell.CurN = nil
                }
            case "source", "s":
                if shell.CurP == nil {
                    errlog("Use a project first")
                    return
                }
                if len(args) < 3 {
                    errlog("Specify a source to use")
                    return
                }
                id, err := strconv.Atoi(args[2])
                if err != nil {
                    errlog("Please specify source id")
                    return
                }
                s := shell.CurP.load_source(id)
                if s != nil {
                    shell.CurS = s
                    shell.CurN = nil
                }
            case "note","n":
                if shell.CurS == nil {
                    errlog("Use a source first")
                    return
                }
                if len(args) < 3 {
                    errlog("Specify a note to use")
                    return
                }
                id, err := strconv.Atoi(args[2])
                if err != nil {
                    errlog("Please specify note id")
                    return
                }
                n := shell.CurS.load_note(id)
                if n != nil {
                    shell.CurN = n
                }
            default:
                errlog("Command not recognized")
        }
    } else {
        errlog("Specify what to use")
    }

}

func (shell *Shell) c_edit(args []string) {
    if shell.CurP == nil {
        errlog("Use a project first")
        return
    }

    if len(args) > 1 {
        switch args[1] {
            case "source","s":
                shell.edits(args)
            case "note","n":
                shell.editn(args)
            default:
                errlog("Command not recognized")
        }
    } else {
        if shell.CurN == nil {
            shell.edits(args)
        } else {
            shell.editn(args)
        }
    } 
    
}

func (shell *Shell) edits(args []string) {

}
func (shell *Shell) editn(args []string) {

}

func (shell *Shell) c_cite(args []string) {

}

func (shell *Shell) c_back(args []string) {
    if shell.CurN != nil {
        shell.CurN = nil
    } else if shell.CurS != nil {
        shell.CurS = nil
    } else {
        shell.CurP = nil
    }
}

func (shell *Shell) c_info(args []string) {
    if len(args) > 1 {
        switch args[1] {
            case "project", "p":
                if len(args) < 3 {
                    errlog("Specify a project to describe")
                    return
                }
                var name string
                for i := 2; i < len(args); i++ {
                    name += args[i] + " "
                }
                name = strings.TrimRight(name, "\n")
                name = strings.TrimSpace(name)
                p := load_project(name)
                if p != nil {
                    p.describe()
                }
            case "source", "s":
                if shell.CurP == nil {
                    errlog("Use a project first")
                    return
                }
                if len(args) < 3 {
                    errlog("Specify a source to describe")
                    return
                }
                id, err := strconv.Atoi(args[2])
                if err != nil {
                    errlog("Please specify source id")
                    return
                }
                s := shell.CurP.load_source(id)
                if s != nil {
                    s.describe()
                }
            case "note","n":
                if shell.CurS == nil {
                    errlog("Use a source first")
                    return
                }
                if len(args) < 3 {
                    errlog("Specify a note to describe")
                    return
                }
                id, err := strconv.Atoi(args[2])
                if err != nil {
                    errlog("Please specify note id")
                    return
                }
                n := shell.CurS.load_note(id)
                if n != nil {
                    n.describe()
                }
            default:
                errlog("Command not recognized")
        }
    } else {
        if shell.CurP == nil {
            errlog("Select project to describe")
        } else if shell.CurS == nil {
            shell.CurP.describe()
        } else if shell.CurN == nil {
            shell.CurS.describe()
        } else {
            shell.CurN.describe()
        }
    }
}


