package svc

import (
	"context"
	"time"

	"github.com/Dcarbon/arch-proto/pb"
	"github.com/Dcarbon/go-shared/gutils"
)

// type ProjectStatus int

// const (
// 	ProjectStatusReject   ProjectStatus = -1
// 	ProjectStatusRegister ProjectStatus = 1
// 	ProjectStatusActived  ProjectStatus = 20
// )

type projectService struct {
	client pb.ProjectServiceClient
}

// type Project struct {
// 	Id           int64          `json:"id"`                                                    //
// 	LocationName string         `json:"locationName,omitempty"    `                            //
// 	Descs        []*ProjectDesc `json:"descs,omitempty"           gorm:"foreignKey:ProjectId"` //

// } //@name Project

// type ProjectDesc struct {
// 	Id        int64  `gorm:"primaryKey"`
// 	ProjectId int64  `gorm:"index:idx_project_desc_lang,unique,priority:1"` //
// 	Language  string `languague"`                                           //
// 	Name      string `name`
// }

func NewRequestProjectService(host string) (IProject, error) {
	cc, err := gutils.GetCCTimeout(host, 5*time.Second)
	if nil != err {
		return nil, err
	}

	var client = &projectService{
		client: pb.NewProjectServiceClient(cc),
	}
	return client, nil
}

type IProject interface {
	GetById(ctx context.Context, in *pb.RPGetById) (*pb.Project, error)
}

func (sv *projectService) GetById(ctx context.Context, req *pb.RPGetById) (*pb.Project, error) {
	project, err := sv.client.GetById(ctx, &pb.RPGetById{ProjectId: req.ProjectId})
	if err != nil {
		return nil, err
	}
	return project, nil
}
