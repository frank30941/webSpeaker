package speaker

//package main

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"sync"
)

func in_array(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

type Speaker struct {
	folder string
	mp3s   sync.Map
}

var instance *Speaker
var once sync.Once
var ch chan []byte

func hash(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return strconv.FormatUint(h.Sum64(), 10)
}
func GetInstance() *Speaker {
	once.Do(func() {
		path, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		instance = &Speaker{
			folder: path + "/audio/",
		}
		instance.loadFiles()
	})
	return instance
}

func (s *Speaker) loadFiles() {
	files, err := ioutil.ReadDir(s.folder)
	if err != nil {
		log.Fatal(err)
	}
	if len(files) != 0 {
		for _, file := range files {
			name := file.Name()
			data, err := ioutil.ReadFile(s.folder + name)
			if err != nil {
				log.Fatal(err)
			}
			s.mp3s.Store(name, data)
		}
	}
}

func (s *Speaker) waitFile(name string) []byte {
	mp3, ok := s.mp3s.Load(name)
	if ok {
		fmt.Println("got a mp3 when downloaded.")
		tmp, _ := mp3.([]byte)
		return tmp
	}
	return s.waitFile(name)
}
func (s *Speaker) Speak(text string) []byte {
	fmt.Printf("%d", os.Getpid())
	name := hash(text[1:])
	name = name + ".mp3"
	tmp, ok := s.mp3s.Load(name)
	if ok {
		mp3, ok1 := tmp.([]byte)
		if ok1 {
			fmt.Println("have a mp3.")
			return mp3
		} else {
			fmt.Println("wait a mp3.")
			return s.waitFile(name)
		}
	}
	s.mp3s.Store(name, "wait")
	var mp3 []byte
	fmt.Println("download a mp3.")
	url := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), "en")
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	mp3, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	s.mp3s.Store(name, mp3)
	output, err := os.Create(s.folder + name)
	if err != nil {
		log.Fatal(err)
	}
	_, err = output.Write(mp3)
	defer output.Close()
	if err != nil {
		log.Fatal(err)
	}
	return mp3
}

/*func main() {
	speaker := GetInstance()
	speaker.Speak("The body of Savannahr.")
}*/
