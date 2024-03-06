// gohttpd project main.go
package main

import (
	"flag"
	"fmt"
	"mime"
	"path"

	//	"html/template"

	//"io/ioutil"
	"log"
	"net/http"

	//	"os"
	"strings"
	"time"
)

var ROOT string
var (
	root  = flag.String("dir", ".", "http root folder")
	port  = flag.String("p", "80", "Address of the http server")
	key   = flag.String("key", "", "key.pem file， using by https")
	cert  = flag.String("cert", "", "cert.pem file， using by https")
	cache = flag.Bool("cache", true, "if cache is false, tell broswer dont't cache file")
	gz = flag.Bool("gz",true, "enable/disable gzip Content-Type support")
	cors = flag.Bool("cors", true, "enable/disable cors access")
	redirUrl = flag.String("redir", "", "redirect url to index")
)

func main() {
	flag.Parse()
	proto := "http"
	if isHttps() {
		proto = "https"
		if *port == "80" {
			*port = "443"
		}
	}
	fmt.Printf("START %s (DIR: %s  PORT: %s )\n", proto, *root, *port)
	if *cache == false {
		fmt.Printf("no-cache enable!\n")
	}
	initMimeExt()
	start(*root, *port)
}
func initMimeExt() {
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".gz", "gzip")
	mime.AddExtensionType(".gz", "gzip")
}
func isHttps() bool {
	return len(*key) > 0 && len(*cert) > 0
}

func fileHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)

	if *redirUrl != "" && r.URL.Path != "/" {
		if strings.HasPrefix(r.URL.Path, *redirUrl) {
			r.URL.Path="/"
			print("redirect --> index")
		}
	}

	if *cache == false {
		w.Header().Add("Cache-Control", "no-cache")
	}
	url := r.URL.Path
	ext := path.Ext(url)
	if(*gz == true) && (ext == ".gz"){
		ext = path.Ext(url[0:len(url)-len(ext)])
		w.Header().Set("Content-Encoding", "gzip")
	}

	if mimeType := mime.TypeByExtension(ext); len(mimeType) > 0 {
		w.Header().Set("Content-Type", mimeType)
	}
	if(*cors == true){
		enableCors(w,r);
	}

	_handler.ServeHTTP(w, r)
}
func enableCors(w http.ResponseWriter, r *http.Request){
    if origin := r.Header.Get("Origin"); origin != "" {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token")
        w.Header().Set("Access-Control-Expose-Headers", "Authorization")
    }
}

var _handler http.Handler

func start(root, port string) {

	root = strings.Replace(root, "\\", "/", -1)
	root = strings.TrimRight(root, "/") + "/"
	ROOT = root

	//http.Handle("/", http.FileServer(http.Dir(root))) //use fileserver directly
	//http.HandleFunc("/", Handler)
	_handler = http.FileServer(http.Dir(root))
	http.HandleFunc("/", fileHandle)

	s := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    12 * time.Hour,
		WriteTimeout:   12 * time.Hour,
		MaxHeaderBytes: 1 << 20,
	}
	if isHttps() {
		log.Fatal(s.ListenAndServeTLS(*cert, *key))
	} else {
		log.Fatal(s.ListenAndServe())
	}
}

//func Handler(w http.ResponseWriter, r *http.Request) {
//	path := r.URL.Path[1:]
//	//if path == "" {
//	//	path = "index.html"
//	//}
//	f, err := os.Open(ROOT + path)
//	if err != nil {
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//	defer f.Close()
//	fi, err := f.Stat()
//	if err != nil {
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}

//	if fi.IsDir() {
//		writeFilelist(w, f)
//	} else {
//		println(path)
//		if *cache == false {
//			w.Header().Add("Cache-Control", "no-cache")
//		}
//		http.ServeFile(w, r, path)
//	}

//	/*else if run && fi.Size() < BUFSIZE {
//		fileName := fi.Name()
//		if index := strings.LastIndex(fileName, "."); index >= 0 {
//			buf, err := ioutil.ReadAll(f)
//			if err != nil {
//				return
//			}

//			filter := doFilter(strings.ToLower(fileName[index+1:]))
//			if filter != nil {
//				output := filter.filter(buf)
//				w.Write(output.Bytes())
//			} else {
//				w.Write(buf)
//			}
//			StatAdd(fileName)
//			return

//		}
//	} else {
//		code := http.StatusOK
//		fileSize := fi.Size()
//		codeRange, sendSize, err := doHeaderRange(w, r, fileSize, f)
//		if codeRange >= 0 {
//			code = codeRange
//		}
//		w.WriteHeader(code)
//		if err != nil {
//			w.Write([]byte(err.Error()))
//		}
//		if r.Method != "HEAD" {
//			writeFile(w, f, sendSize, fi.Name())
//		}
//	}
//	*/
//}

//const BUFSIZE int64 = 512 * 1024

//func writeFile(w http.ResponseWriter, f *os.File, sendSize int64, fileName string) {

//	bufSize := BUFSIZE
//	if sendSize < BUFSIZE {
//		bufSize = sendSize
//	}
//	buf := make([]byte, bufSize)
//	for {
//		//		curbuf := buf
//		if sendSize < bufSize {
//			buf = buf[0:sendSize]
//		}
//		rlen, err := f.Read(buf)
//		if err != nil {
//			break
//		}
//		w.Write(buf[0:rlen])
//		sendSize -= int64(rlen)
//		if sendSize <= 0 {
//			break
//		}
//	}
//}

//type FileInfo struct {
//	Name string
//	Url  string
//	Size int64
//}

//var tmpl, _ = template.New("filelist").Parse(TMPL_FILELIST)

//func writeFilelist(w http.ResponseWriter, f *os.File) {
//	files, err := f.Readdir(0)
//	if err != nil {
//		fmt.Fprintf(w, "404")
//		return
//	}
//	fileInfos := []*FileInfo{}
//	for _, file := range files {
//		fileName := file.Name()
//		fileSize := file.Size()
//		if file.IsDir() {
//			fileName += "/"
//		}
//		url := fileName
//		fileInfos = append(fileInfos, &FileInfo{Name: fileName, Url: url, Size: fileSize})
//	}

//	err = tmpl.Execute(w, fileInfos)
//	checkErr(err)
//}

//func checkErr(err error) bool {
//	if err != nil {
//		fmt.Println("error: ", err.Error())
//		return true
//	}
//	return false
//}

//const TMPL_FILELIST = `<html>
//<head></head>
//<body>
//<table border="0" cellspacing="8">
//	{{with .}}
//	{{range .}}
//	<tr>
//		<td><a href="{{.Url}}">{{.Name}}</a></td>
//		<td align="right">{{.Size}}B</td>
//	</tr>
//	{{end}}
//	{{end}}
//</body>
//</html>`
