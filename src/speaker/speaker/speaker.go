package speaker

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type Speaker struct {
	folder string
	mp3s   sync.Map
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
			fmt.Println(name)
			s.mp3s.Store(name, data)
		}
	}
}
func (s *Speaker) Speak(text string) []byte {
	name := hash(text[1:])
	name = name + ".mp3"
	tmp, ok := s.mp3s.Load(name)
	if ok {
		mp3, ok1 := tmp.([]byte)
		if ok1 {
			//fmt.Println("have a mp3.")
			//fmt.Println(reflect.TypeOf(tmp).String())
			return mp3
		}
		mp3 = nil
		tmp = nil
		//fmt.Println("wait a mp3.")
		for {
			mp3, ok := s.mp3s.Load(name)
			if ok {
				tmp, ok := mp3.([]byte)
				if ok {
					mp3 = tmp
					tmp = nil
					break
				}
			}
		}
		return mp3

	}
	tmp = nil
	s.mp3s.Store(name, "wait")
	fmt.Println("download a mp3.")
	url := fmt.Sprintf("http://translate.google.com/translate_tts?ie=UTF-8&total=1&idx=0&textlen=32&client=tw-ob&q=%s&tl=%s", url.QueryEscape(text), "en")
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	mp3, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	s.mp3s.Store(name, mp3)
	output, err := os.Create(s.folder + name)
	if err != nil {
		log.Fatal(err)
	}
	_, err = output.Write(mp3)
	output = nil
	if err != nil {
		log.Fatal(err)
	}
	response.Body.Close()
	output.Close()
	//fmt.Println(reflect.TypeOf(mp3).String())
	return mp3
}
