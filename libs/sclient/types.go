package sclient

// IStorage : client of storage service
type IStorage interface {
	UploadToImage(fname string, group string) (string, error)
	UploadToProject(fname string, projectId int64) (string, error)
	UploadToUser(fname string, userId int64) (string, error)
}
