package main

import (
    "fmt"
    "time"
    "os"
    "path/filepath"
    "log"
    "encoding/json"
    "strconv"
    "strings"
)

type SourceType string
const (
    Encyclopedia SourceType = "Encyclopedia"
    Documentary  SourceType = "Documentary"
    Journal      SourceType = "Journal"
    Website      SourceType = "Website"
    Article      SourceType = "Article"
    Letter       SourceType = "Letter"
    Book         SourceType = "Book"
)

type Name struct {
    First  string
    Middle string
    Last   string
}

func (name Name) mla() string {
    if name.First != "" && name.Last != "" && name.Middle != "" {
        return name.Last + ", " + name.First + string(name.Middle) + "."
    } else {
        return ""
    }
}

type Project struct {
    Title    string
    Created  time.Time
    Modified time.Time
    Sourcec  int
}

type Source struct {
    Id       int
    Project  string
    Title    string
    Created  time.Time
    Modified time.Time

    Notec    int

    Stype     SourceType
    Url       string
    Author    Name
    Publisher string
    Cdate     string
    Adate     string

    Citation  string
}

type Note struct {
    Id       int
    Source   int
    Project  string
    Title    string
    Created  time.Time
    Modified time.Time

    Content  string
}

func home() string {
    d, err := os.UserHomeDir()
    if err != nil {
        log.Fatal(err)
    }
    return d
}

func create_project(title string) *Project {
    os.Mkdir(filepath.Join(home(), ".noto"), 0755)

    title = strings.Join(strings.Fields(strings.TrimSpace(title)), " ")

    path := filepath.Join(home(), ".noto", title)

    _, err := os.Stat(path); if os.IsNotExist(err) {
        os.Mkdir(path, 0755)
    } else {
        errlog("Project already exists")
        return nil
    }

    project := Project {
        Title:title,
        Created:time.Now(),
        Modified:time.Now(),
        Sourcec:0,
    }

    bytes, err := json.Marshal(project)
    if err != nil {
        log.Fatal(err)
    }

    err = os.WriteFile(filepath.Join(path, "project.json"), bytes, 0644)
    if err != nil {
        log.Fatal(err)
    }

    return &project;
}

func load_project(title string) *Project {
    path := filepath.Join(home(), ".noto", title)
//    Shprintln("Attempting to load " + path)
    _, err := os.Stat(path); if os.IsNotExist(err) {
        errlog("Project does not exist")
        return nil
    }

    conf, err := os.ReadFile(filepath.Join(path, "project.json"))
    if err != nil {
        log.Fatal(err)
    }

    project := Project{}
    confs := string(conf)
    json.Unmarshal([]byte(confs), &project)
    return &project
}

func (project *Project) update() {
    path := filepath.Join(home(), ".noto", project.Title)

    project.Modified = time.Now()

    bytes, err := json.Marshal(project)
    if err != nil {
        log.Fatal(err)
    }

    err = os.WriteFile(filepath.Join(path, "project.json"), bytes, 0644)
    if err != nil {
        log.Fatal(err)
    }
}

func (project *Project) create_source() *Source {
    id := project.Sourcec
    project.Sourcec += 1

    source := Source {
        Id:id,
        Project:project.Title,
        Created:time.Now(),
        Modified:time.Now(),
    }

    bytes, err := json.Marshal(source)
    if err != nil {
        log.Fatal(err)
    }


    path := filepath.Join(home(), ".noto", project.Title, strconv.Itoa(id))
    os.Mkdir(path, 0755)

    err = os.WriteFile(filepath.Join(path, "source.json"), bytes, 0644)
    if err != nil {
        log.Fatal(err)
    }

    project.update()

    return &source
}

func (project *Project) describe() {
    r := "\n"
    r += "\t" + project.Title + "\n\n"
    r += "\tCreated:  " + project.Created.String() + "\n"
    r += "\tModified: " + project.Modified.String() + "\n"
    r += "\tSources:  " + strconv.Itoa(project.Sourcec) + "\n\n"
    fmt.Print(r)
}

func (project *Project) load_source(id int) *Source {
    path :=filepath.Join(home(), ".noto", project.Title, strconv.Itoa(id))

    _, err := os.Stat(path); if os.IsNotExist(err) {
        errlog("Source does not exist")
        return nil
    }

    conf, err := os.ReadFile(filepath.Join(path, "source.json"))
    if err != nil {
        log.Fatal(err)
    }

    source := Source{}
    confs := string(conf)
    json.Unmarshal([]byte(confs), &source)

    return &source
}

func (source *Source) update() {
    path := filepath.Join(home(), ".noto", source.Project, strconv.Itoa(source.Id))

    source.Modified = time.Now()

    bytes, err := json.Marshal(source)
    if err != nil {
        log.Fatal(err)
    }

    err = os.WriteFile(filepath.Join(path, "source.json"), bytes, 0644)
    if err != nil {
        log.Fatal(err)
    }
}

func (source *Source) create_note() *Note {
    id := source.Notec
    source.Notec += 1

    note := Note {
        Id:id,
        Source:source.Id,
        Project:source.Project,
        Created:time.Now(),
        Modified:time.Now(),
    }

    bytes, err := json.Marshal(note)
    if err != nil {
        log.Fatal(err)
    }


    path := filepath.Join(home(), ".noto", source.Project, strconv.Itoa(source.Id))
    os.Mkdir(path, 0755)

    err = os.WriteFile(filepath.Join(path, strconv.Itoa(id) + ".json"), bytes, 0644)
    if err != nil {
        log.Fatal(err)
    }

    source.update()

    return &note
}

func (source *Source) load_note(id int) *Note {
    path := filepath.Join(home(), ".noto", source.Project, strconv.Itoa(source.Id))

    _, err := os.Stat(path); if os.IsNotExist(err) {
        errlog("Source does not exist")
        return nil
    }

    conf, err := os.ReadFile(filepath.Join(path, strconv.Itoa(id) + ".json"))
    if err != nil {
        errlog("Note does not exist")
        return nil
    }

    note := Note{}
    confs := string(conf)
    json.Unmarshal([]byte(confs), &note)

    return &note
}

func (source *Source) describe() {

    t := "No Title"
    if source.Title != "" {
        t = source.Title
    }

    url := ""
    if len(source.Url) < 48 {
        url = source.Url
    } else {
        r := []rune(source.Url)
        url = string(r[0:24])
        url += "... (`url` for full url)"
    }

    r := "\n"
    r += "\tSource " + strconv.Itoa(source.Id) + " - " + t + "\n\n"
    r += "\tAdded:    " + source.Created.String() + "\n"
    r += "\tModified: " + source.Modified.String() + "\n"
    r += "\tNotes:    " + strconv.Itoa(source.Notec) + "\n\n"

    r += "\tType:      " + string(source.Stype) + "\n"
    r += "\tAuthor:    " + source.Author.mla() + "\n"
    r += "\tPublisher: " + source.Publisher + "\n"
    r += "\tURL:       " + url + "\n"
    r += "\tCreated: " + source.Cdate + "\n"
    r += "\tAccessed: " + source.Adate + "\n\n"
    fmt.Print(r)
}

func (note *Note) update() {
    path := filepath.Join(home(), ".noto", note.Project, strconv.Itoa(note.Source))

    note.Modified = time.Now()

    bytes, err := json.Marshal(note)
    if err != nil {
        log.Fatal(err)
    }

    err = os.WriteFile(filepath.Join(path, strconv.Itoa(note.Id) + ".json"), bytes, 0644)
    if err != nil {
        log.Fatal(err)
    }
}

func (note *Note) describe() {
    t := "No Title"
    if note.Title != "" {
        t = note.Title
    }

    s := strconv.Itoa(len(note.Content)) + "B"

    r := "\n"
    r += "\tNote " + strconv.Itoa(note.Id) + " - " + t + "\n\n"
    r += "\tAdded:    " + note.Created.String() + "\n"
    r += "\tModified: " + note.Modified.String() + "\n"
    r += "\tSize:     " + s + "\n\n"
    fmt.Print(r)
}

