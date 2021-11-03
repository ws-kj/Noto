package main

import (
    "os"
    "path/filepath"
    "log"
    "strconv"
    "strings"
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

var (
    cur_proj     *Project
    cur_src     *Source
    cur_note     *Note
    pstr   string
)
/*
func Cprintln(c Color, s string) {
    Shprintln(string(c), s, string(Reset))
}
*/
func errlog(s string) {
//    Cprintln(Red, "\n\t"+s+"\n")
    Shprintln("\n\n\t" + s)
}

func build_prompt() string {
    if cur_proj == nil {
        return "Noto ?>";
    }
    r := cur_proj.Title;
    if cur_src == nil {
        return r + " ?>";
    }
    r += " | "
    if cur_src.Title == "" {
        r += "Source "
        r += strconv.Itoa(cur_src.Id)
    } else {
        r += cur_src.Title
    }
    if cur_note == nil {
        r += " ?>"
        return r
    }
    r += " | "
    if cur_note.Title == "" {
        r += "Note "
        r += strconv.Itoa(cur_note.Id)
    } else {
        r += cur_note.Title
    }
    r += " ?>"
    return r
}

/*
func prompt() {
    p := build()
    fmt.Print(p)
    reader := bufio.NewReader(os.Stdin)
    in, _ := reader.ReadString('\n')
    process(in)
    prompt()
}
*/
func process_command(input string) {
    args := strings.Fields(input)

    if len(args) == 0 {
        return
    }

    switch args[0] {
        case "list", "l":
            c_list(args)
        case "new", "n":
            c_new(args)
        case "use", "u":
            c_use(args)
        case "edit", "e":
            c_edit(args)
        case "cite", "c":
            c_cite(args)
        case "info", "i":
            c_info(args)
        case "back", "b", "..":
            c_back(args)
        default:
            errlog("Command not recognized")
    }
}

func c_list(args []string) {
    if len(args) > 1 {
        switch args[1] {
            case "projects","project","p":
                listp(args)
            case "sources","source","s":
                if cur_proj != nil {
                    lists(args)
                } else {
                    errlog("Use a project first")
                }
            case "notes","note","n":
                if cur_src != nil {
                    listn(args)
                } else {
                    errlog("Use a source first")
                }
            default:
                errlog("Command not recognized")
        }
    } else {
        if cur_proj == nil {
            listp(args)
        } else if cur_src == nil {
            lists(args)
        } else {
            listn(args)
        }
    }
}

//DRY!!
func listp(args []string) {
    f, err := os.Open(filepath.Join(home(), ".noto"))
    if err != nil {
        log.Fatal(err)
    }
    files, err := f.Readdir(0)
    if err != nil {
        log.Fatal(err)
    }
    if len(files) == 0 {
        errlog("No projects yet! Use `new project ___` to create one.")
    } else {
        Shprintln("\n")
    }
    for _, proj := range files {
        if proj.IsDir() {
            Shprintln("\t" + proj.Name())
        }
    }
}
func lists(args []string) {
    f, err := os.Open(filepath.Join(home(), ".noto", cur_proj.Title))
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
        Shprintln("\n")
    }
    for _, n := range entries {
        entry := "\tSource "
        entry += strconv.Itoa(n)
        entry += " - "
        src := cur_proj.load_source(n)
        if src.Title != "" {
            entry += src.Title
        } else {
            entry += "No Title"
        }
        Shprintln(entry)
    }
}

func listn(args []string) {
    f, err := os.Open(filepath.Join(home(), ".noto", cur_proj.Title, strconv.Itoa(cur_src.Id)))
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
    } else {
        Shprintln("\n")
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

    for _, n := range entries {
        entry := "\tNote "
        entry += strconv.Itoa(n)
        entry += " - "
        note := cur_src.load_note(n)
        if note.Title != "" {
            entry += note.Title
        } else {
            entry += "No Title"
        }
        Shprintln("\t" + entry)
    }
}

func c_new(args []string) {
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
                cur_proj = p
                cur_src = nil
                cur_note = nil
            }
        case "source","s":
            if cur_proj == nil {
                errlog("Use a project first")
                return
            }
            s := cur_proj.create_source()
            cur_src = s
        case "note","n":
            if cur_src == nil {
                errlog("Use a source first")
                return
            }
            n := cur_src.create_note()
            cur_note = n
        default:
            errlog("Command not recognized")
    }
}

func c_use(args []string) {
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
                    cur_proj = p
                    cur_src = nil
                    cur_note = nil
                }
            case "source", "s":
                if cur_proj == nil {
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
                s := cur_proj.load_source(id)
                if s != nil {
                    cur_src = s
                    cur_note = nil
                }
            case "note","n":
                if cur_src == nil {
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
                n := cur_src.load_note(id)
                if n != nil {
                    cur_note = n
                }
            default:
                errlog("Command not recognized")
        }
    } else {
        errlog("Specify what to use")
    }

}

func c_edit(args []string) {
    if cur_proj == nil {
        errlog("Use a project first")
        return
    }

    if len(args) > 1 {
        switch args[1] {
            case "source","s":
                edits(args)
            case "note","n":
                editn(args)
            default:
                errlog("Command not recognized")
        }
    } else {
        if cur_note == nil {
            edits(args)
        } else {
            editn(args)
        }
    } 
    
}

func edits(args []string) {

}
func editn(args []string) {

}

func c_cite(args []string) {

}

func c_back(args []string) {
    if cur_note != nil {
        cur_note = nil
    } else if cur_src != nil {
        cur_src = nil
    } else {
        cur_proj = nil
    }
}

func c_info(args []string) {
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
                if cur_proj == nil {
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
                s := cur_proj.load_source(id)
                if s != nil {
                    s.describe()
                }
            case "note","n":
                if cur_src == nil {
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
                n := cur_src.load_note(id)
                if n != nil {
                    n.describe()
                }
            default:
                errlog("Command not recognized")
        }
    } else {
        if cur_proj == nil {
            errlog("Select project to describe")
        } else if cur_src == nil {
            cur_proj.describe()
        } else if cur_note == nil {
            cur_src.describe()
        } else {
            cur_note.describe()
        }
    }
}


