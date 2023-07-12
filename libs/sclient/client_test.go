package sclient

import (
	"fmt"
	"log"
	"testing"
)

const jwt = ""
const host = "http://localhost:4005"

var s IStorage

func init() {
	s, _ = NewStorage(host, jwt)
}

func TestUploadImage(t *testing.T) {
	path, err := s.UploadToImage("./static/gg.png", "1")
	if nil != err {
		log.Fatalln("UploadImage error: ", err)
	}
	fmt.Println(host + path)
}

func TestUploadProject(t *testing.T) {
	path, err := s.UploadToProject("./static/gg.png", 1)
	if nil != err {
		log.Fatalln("UploadFile error: ", err)
	}
	fmt.Println(host + path)
}

func TestUploadUser(t *testing.T) {
	path, err := s.UploadToUser("./static/img.png", 1)
	if nil != err {
		log.Fatalln("UploadAudio error: ", err)
	}
	fmt.Println(host + path)
}
