package sclient

import (
	"errors"
	"fmt"
	"log"

	"github.com/Dcarbon/go-shared/libs/aidh"
)

// Storage :
type Storage struct {
	host   string
	jwt    string
	caller *aidh.Caller
}

// NewStorage :
func NewStorage(host, jwt string) (IStorage, error) {
	if host == "" {
		return nil, errors.New("missing storage host")
	}

	if jwt == "" {
		return nil, errors.New("missing storage token")
	}

	var st = &Storage{
		host: host,
		jwt:  jwt,
		caller: aidh.NewCaller(map[string]string{
			"Authorization": "Bearer " + jwt,
		}),
	}

	return st, nil
}

// PostImage :
func (s *Storage) UploadToImage(fname string, group string) (string, error) {
	var path = "/static/images/upload"
	return s.Upload(path, fname, group)
}

func (s *Storage) UploadToProject(fname string, projectId int64) (string, error) {
	var path = "/static/projects/upload"
	return s.Upload(path, fname, fmt.Sprintf("%d", projectId))
}

func (s *Storage) UploadToUser(fname string, userId int64) (string, error) {
	var path = "/static/users/upload"
	return s.Upload(path, fname, fmt.Sprintf("%d", userId))
}

func (s *Storage) Upload(path, file, group string,
) (string, error) {
	var URL = s.host + path
	log.Println("URL: ", URL)
	rs := &rsFile{}
	err := s.caller.FormFile(
		URL,
		aidh.FormFields{
			{
				Key:   "group",
				Type:  aidh.FormFieldText,
				Value: group,
			},
			{
				Key:   "file",
				Type:  aidh.FormFieldFile,
				Value: file,
			},
		},
		rs,
	)
	if nil != err {
		return "", err
	}
	return rs.File, nil
}

func (s *Storage) GetHost() string {
	return s.host
}

type rsFile struct {
	File string `json:"file"`
}
