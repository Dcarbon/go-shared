package sclient

import (
	"fmt"
	"log"
	"testing"
)

const jwt = "x"
const host = "http://54.151.254.13:4100"

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
	s.UploadToProject("./static/55/069815ae8de635b86cf79.jpg", 55)
	// if nil != err {
	// 	log.Fatalln("UploadFile error: ", err)
	// }
	//fmt.Println(host + path)
}

func TestUploadUser(t *testing.T) {
	path, err := s.UploadToUser("./static/img.png", 1)
	if nil != err {
		log.Fatalln("UploadAudio error: ", err)
	}
	fmt.Println(host + path)
}
