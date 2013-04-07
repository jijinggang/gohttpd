// gohttpd project main.go
package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var ROOT string = ""

func main() {
	root := "."
	port := "80"
	flag.Parse()
	if flag.NArg() == 2 {
		root = flag.Arg(0)
		port = flag.Arg(1)
	} else {
		fmt.Println("Usage: gohttpd.exe root_dir port (defaut: gohttpd . 80)")
	}
	fmt.Printf("START gohttpd (DIR: %s  PORT: %s )\n", root, port)
	StatStart()
	start(root, port)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	//if path == "" {
	//	path = "index.html"
	//}
	path = ROOT + path
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(w, "404")
		return
	}
	fi, err := f.Stat()
	if err != nil {
		fmt.Fprintf(w, "404")
		return
	}
	//判断是否目录
	if fi.IsDir() {
		stat := ("1" == r.FormValue("stat"))
		writeFilelist(w, f, stat)
	} else {
		writeFile(w, f, fi.Size(), fi.Name())
	}
	f.Close()
}

func writeFile(w http.ResponseWriter, f *os.File, fileSize int64, fileName string) {
	const BUFSIZE = 512 * 1024
	if fileSize > BUFSIZE {
		fileSize = BUFSIZE
	}
	buf := make([]byte, fileSize)
	for {
		rlen, err := f.Read(buf)
		if err != nil {
			break
		}
		if fileSize < BUFSIZE { //filter only on small file
			rIndex := strings.LastIndex(fileName, ".")
			if rIndex >= 0 {
				filter := doFilter(strings.ToLower(fileName[rIndex+1:]))
				if filter != nil {
					output := filter.filter(buf[0:rlen])
					w.Write(output.Bytes())
					StatAdd(fileName)
					return
				}
			}
		}
		w.Write(buf[0:rlen])
	}
	StatAdd(fileName)
}

type FileInfo struct {
	Name  string
	Size  int64
	Count int64 //download count
}

var tmpl, _ = template.New("filelist").Parse(TMPL_FILELIST)
var tmpl2, _ = template.New("filelist").Parse(TMPL_FILELIST_STAT)

func writeFilelist(w http.ResponseWriter, f *os.File, stat bool) {
	files, err := f.Readdir(0)
	if err != nil {
		fmt.Fprintf(w, "404")
		return
	}
	fileInfos := []*FileInfo{}
	for _, file := range files {
		fileName := file.Name()
		fileSize := file.Size()
		if file.IsDir() {
			fileName += "/"
		}
		fileInfos = append(fileInfos, &FileInfo{Name: fileName, Size: fileSize})
	}
	if stat {
		StatGet(fileInfos)
		err = tmpl2.Execute(w, fileInfos)
	} else {
		err = tmpl.Execute(w, fileInfos)
	}
	checkErr(err)
}

func start(root, port string) {
	root = strings.Replace(root, "\\", "/", -1)
	root = strings.TrimRight(root, "/") + "/"
	//http.Handle("/", http.FileServer(http.Dir(root)))
	ROOT = root
	http.HandleFunc("/", Handler)
	s := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    12 * time.Hour,
		WriteTimeout:   12 * time.Hour,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func checkErr(err error) bool {
	if err != nil {
		fmt.Println("error: ", err.Error())
		return true
	}
	return false
}

const TMPL_FILELIST = `<html>
<head></head>
<body>
<table border="0" cellspacing="8">
	{{with .}}
	{{range .}}  
	<tr>
		<td><a href="{{.Name}}">{{.Name}}</a></td>
		<td align="right">{{.Size}}B</td>
	</tr>
	{{end}} 
	{{end}}
</body>
</html>`
const TMPL_FILELIST_STAT = `<html>
<head></head>
<body>
<table border="0" cellspacing="8">
	{{with .}}
	{{range .}}  
	<tr>
		<td><a href="{{.Name}}">{{.Name}}</a></td>
		<td align="right">{{.Size}}B</td>
		<td align="right">{{.Count}}</td>
	</tr>
	{{end}} 
	{{end}}
</body>
</html>`
