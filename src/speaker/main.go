package main

import (
	"encoding/binary"
	"fmt"
	"net/http"
	"os"
	"speaker/speaker"
	"strconv"
)

func main() {
	speaker.GetInstance()
	fmt.Printf("%d", os.Getpid())
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)

}
func handler(w http.ResponseWriter, r *http.Request) {
	tmp := speaker.GetInstance()
	mp3 := tmp.Speak(r.URL.Path)
	w.Header().Set("Content-type", "audio/mpeg")
	w.Header().Set("Content-Disposition", "inline;filename=output.mp3")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", strconv.Itoa(binary.Size(mp3)))
	w.Write(mp3)
}
