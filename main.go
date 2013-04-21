// gohttpd project main.go
package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var ROOT string = ""

func main() {
	root := "."
	//	root = "e:/www"
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
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//run=1 -> give *.md markdown html && give dir statics data
	run := ("1" == r.FormValue("run"))
	if fi.IsDir() {
		writeFilelist(w, f, run)
	} else if run && fi.Size() < BUFSIZE {
		fileName := fi.Name()
		if index := strings.LastIndex(fileName, "."); index >= 0 {
			buf, err := ioutil.ReadAll(f)
			if err != nil {
				return
			}

			filter := doFilter(strings.ToLower(fileName[index+1:]))
			if filter != nil {
				output := filter.filter(buf)
				w.Write(output.Bytes())
			} else {
				w.Write(buf)
			}
			StatAdd(fileName)
			return

		}
	} else {
		code := http.StatusOK
		fileSize := fi.Size()
		codeRange, sendSize, err := doHeaderRange(w, r, fileSize, f)
		if codeRange >= 0 {
			code = codeRange
		}
		w.WriteHeader(code)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		if r.Method != "HEAD" {
			writeFile(w, f, sendSize, fi.Name())
		}
	}
}

const BUFSIZE int64 = 512 * 1024

func writeFile(w http.ResponseWriter, f *os.File, sendSize int64, fileName string) {

	bufSize := BUFSIZE
	if sendSize < BUFSIZE {
		bufSize = sendSize
	}
	buf := make([]byte, bufSize)
	for {
		//		curbuf := buf
		if sendSize < bufSize {
			buf = buf[0:sendSize]
		}
		rlen, err := f.Read(buf)
		if err != nil {
			break
		}
		w.Write(buf[0:rlen])
		sendSize -= int64(rlen)
		if sendSize <= 0 {
			break
		}
	}
	StatAdd(fileName)
}

type FileInfo struct {
	Name  string
	Url   string
	Size  int64
	Count int64 //download count
}

var tmpl, _ = template.New("filelist").Parse(TMPL_FILELIST)
var tmpl2, _ = template.New("filelist").Parse(TMPL_FILELIST_STAT)

func writeFilelist(w http.ResponseWriter, f *os.File, run bool) {
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
		url := fileName
		if run {
			url += "?run=1"
		}

		fileInfos = append(fileInfos, &FileInfo{Name: fileName, Url: url, Size: fileSize})
	}
	if run {
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
		<td><a href="{{.Url}}">{{.Name}}</a></td>
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
		<td><a href="{{.Url}}">{{.Name}}</a></td>
		<td align="right">{{.Size}}B</td>
		<td align="right">{{.Count}}</td>
	</tr>
	{{end}} 
	{{end}}
</body>
</html>`
