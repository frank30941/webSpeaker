package main

import (
	"fmt"
	"net/http"
	"os"
	"speaker/speaker"
)

func main() {
	speaker.GetInstance()
	fmt.Printf("%d", os.Getpid())
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)

}
func handler(w http.ResponseWriter, r *http.Request) {
	tmp := speaker.GetInstance()
	w.Write(tmp.Speak(r.URL.Path))
}
