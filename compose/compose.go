package compose

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

var defaultRes string = "docker.io/library/"
var CREATE_CONTAINER_ERROR error = errors.New("创建容器失败")

//Project Project
type Project struct {
	Name     string
	Services map[string]*Service
	Client   *client.Client
	DevNull  io.Writer
}

//NewProject get new project by compose file
func NewProject(name string, compose string) *Project {
	conf := LoadYaml(compose)
	services := getServices(name, conf)
	client, err := client.NewEnvClient()

	if err != nil {
		panic(fmt.Sprintln(err))
	}
	file, err := os.Open(os.DevNull)
	if err != nil {
		panic(err)
	}
	project := &Project{
		Name:     name,
		Services: services,
		Client:   client,
		DevNull:  os.NewFile(file.Fd(), "/dev/null"),
	}

	return project

}

func (pro *Project) Up(services []string) {
	if len(services) == 0 {

	}
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

//Up up service
func (srv *Service) Up(client *client.Client, stdout io.Writer) error {
	ctx := context.Background()
	image := srv.Config.Image
	i := strings.Split(srv.Config.Image, "/")
	fmt.Println(i)
	if len(i) <= 1 {
		image = "docker.io/library/" + srv.Config.Image
	}
	read, err := client.ImagePull(ctx, image, *new(types.ImagePullOptions))
	if err != nil {
		return err
	}
	io.Copy(stdout, read)
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
