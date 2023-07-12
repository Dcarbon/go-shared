package main

import (
	"log"

	"github.com/Dcarbon/go-shared/libs/sclient"
	"github.com/Dcarbon/go-shared/libs/utils"
)

const jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE5OTkxNTIzMTksImlkIjoxLCJyb2xlIjoic3VwZXItYWRtaW4ifQ.8wFklyWSVI1iy35MH7kSoxesxWYTqYdQg9c5cGBXVGE"
const host = "http://localhost:4005"

func main() {
	client, _ := sclient.NewStorage(host, jwt)
	path, err := client.UploadToProject("../static/gg.png", 1)
	utils.PanicError("", err)
	log.Println("Path: ", path)
}
