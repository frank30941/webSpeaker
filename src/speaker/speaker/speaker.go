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
	Q_tmps []string
	mp3s   map[string][]byte
	sync.RWMutex
}

var instance *Speaker
var once sync.Once

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
			Q_tmps: []string{},
			mp3s:   make(map[string][]byte),
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
			s.mp3s[name] = data
		}
	}
}

func (s *Speaker) readMap(key string) ([]byte, bool) {
    s.RLock()
    value, ok:= s.mp3s[key]
	s.RUnlock()
	if ok {
		return value, true
	}
	return nil, false
}

func (s *Speaker) writeMap(key string, value []byte) {
    s.Lock()
    s.mp3s[key] = value
    s.Unlock()
}
func (s *Speaker) waitFile(name string) []byte {
	mp3, ok := s.readMap(name)
	if ok {
		fmt.Println("got a mp3 when downloaded.")
		return mp3
	}
	return s.waitFile(name)
}
func (s *Speaker) Speak(text string) []byte {
	fmt.Printf("%d", os.Getpid())
	name := hash(text[1:])
	name = name + ".mp3"
	if len(s.mp3s) != 0 {
		mp3, ok := s.readMap(name)
		if ok {
			fmt.Println("have a mp3.")
			return mp3
		}
	}
	var mp3 []byte
	Qok, _ := in_array(name, s.Q_tmps)
	if Qok {
		fmt.Println("wait a mp3.")
		return s.waitFile(name)
	} else {
		fmt.Println("download a mp3.")
		s.Q_tmps = append(s.Q_tmps, name)
		url := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), "en")
		response, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()
		mp3, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		s.writeMap(name, mp3)
		Qok, index := in_array(name, s.Q_tmps)
		if Qok {
			s.Q_tmps = RemoveIndex(s.Q_tmps, index)
		}
		output, err := os.Create(s.folder + name)
		if err != nil {
			log.Fatal(err)
		}
		_, err = output.Write(mp3)
		defer output.Close()
		if err != nil {
			log.Fatal(err)
		}
		mp3, ok := s.readMap(name)
		if ok {
			fmt.Println("downloaded a mp3.")
			return s.mp3s[name]
		}
		return mp3
	}
	return mp3
}

/*func main() {
	speaker := GetInstance()
	speaker.Speak("The body of Savannahr.")
}*/
