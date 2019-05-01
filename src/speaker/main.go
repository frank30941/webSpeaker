package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"speaker/speaker"

	_ "net/http/pprof"
)

func main() {
	speaker.GetInstance()
	fmt.Printf("%d \n", os.Getpid())
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)

}
func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "audio/mpeg")
	w.Header().Set("Content-Disposition", "inline;filename=output.mp3")
	w.Header().Set("Transfer-Encoding", "chunked")
	//w.Header().Set("Content-Length", strconv.Itoa(binary.Size(mp3)))
	pipeReader, pipeWriter := io.Pipe()
	go func() {
		tmp := speaker.GetInstance()
		mp3 := tmp.Speak(r.URL.Path)
		io.Copy(pipeWriter, bytes.NewReader(mp3))
		pipeWriter.Close()
	}()
	io.Copy(w, pipeReader)
}
