package service

import "template/model"

type ReadmeService struct{}

var readme *model.Readme = &model.Readme{
	Name:     "template",
	Desc:     "服务模板",
	Upgrades: make([]*model.ReadmeUpgrade, 0)}

func init() {
	readme.Upgrades = append(readme.Upgrades, &model.ReadmeUpgrade{
		Id:      "1000",
		Version: "1.0.0.0",
		Created: "20250429",
		Notes:   "模板"})
}

func Readme() *model.Readme {
	return readme
}
