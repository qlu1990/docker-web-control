package compose

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

var CREATE_CONTAINER_ERROR error = errors.New("创建容器失败")

//Project Project
type Project struct {
	Name     string
	Services map[string]*Service
	Client   *client.Client
}

//NewProject get new project by compose file
func NewProject(name string, compose string) *Project {
	conf := LoadYaml(compose)
	services := getServices(name, conf)
	client, err := client.NewEnvClient()
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	project := &Project{
		Name:     name,
		Services: services,
		Client:   client,
	}

	return project

}

//Service control for service
type Service struct {
	ProjectName      string
	Name             string
	Config           *container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
	DependsOn        []string
}

//GetNewService get new service
func GetNewService(projectName string, name string) *Service {
	service := &Service{}
	service.ProjectName = projectName
	service.Name = name
	service.Config = &container.Config{}
	service.HostConfig = &container.HostConfig{}
	service.NetworkingConfig = &network.NetworkingConfig{}
	return service
}

func (srv *Service) Up(client *client.Client) error {
	ctx := context.Background()
	read, err := client.ImagePull(ctx, srv.Config.Image, *new(types.ImagePullOptions))
	if err != nil {
		return err
	}
	read.Read()
	createBody, err := client.ContainerCreate(ctx, srv.Config, srv.HostConfig, srv.NetworkingConfig, srv.ProjectName+"_"+srv.Name)
	if err != nil {
		return err
	}
	err = client.ContainerStart(ctx, createBody.ID, *new(types.ContainerStartOptions))
	if err != nil {
		return err
	}
	return nil
}
