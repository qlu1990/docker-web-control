package compose

import (
	"fmt"
	"io/ioutil"

	"github.com/docker/go-connections/nat"
	"gopkg.in/yaml.v2"
)

//LoadYaml load yaml  file
func LoadYaml(fileName string) map[interface{}]interface{} {
	conf := make(map[interface{}]interface{})
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	return conf
}

func getServices(projectName string, conf map[interface{}]interface{}) map[string]*Service {
	services := make(map[string]*Service, 0)
	for k, v := range conf {
		if k.(string) == "services" {
			for srv, srvConf := range v.(map[interface{}]interface{}) {
				service := GetNewService(projectName, srv.(string))
				for param, values := range srvConf.(map[interface{}]interface{}) {
					switch param.(string) {
					case "ports":
						exposedPorts, binding, err := nat.ParsePortSpecs(getStringSlice(values))
						if err != nil {
							panic(fmt.Sprintln(err))
						}
						service.Config.ExposedPorts = exposedPorts
						service.HostConfig.PortBindings = binding
						break
					case "environment":
						service.Config.Env = getStringSliceFromMap(values)
						break
					case "image":
						service.Config.Image = values.(string)
						break
					case "volumes":
						service.HostConfig.Binds = getStringSlice(values)
						break
					case "depends_on":
						service.DependsOn = getStringSlice(values)
						break
					case "cmd":
						service.Config.Cmd = getStringSlice(values)
						break
					case "entrypoint":
						service.Config.Entrypoint = getStringSlice(values)
						break
						// case "cmd":
						// 	service.Cmd = getStringSlice(values)
						// default:
						// 	service.OtherParam[param.(string)] = values.(string)
					}
				}
				services[service.Name] = service
			}
		}
	}
	return services
}

func getStringSlice(s interface{}) []string {
	s1 := s.([]interface{})
	ret := make([]string, 0)
	for _, v := range s1 {
		ret = append(ret, v.(string))
	}
	return ret

}

func getStringSliceFromMap(s interface{}) []string {
	s1 := s.(map[interface{}]interface{})
	ret := make([]string, 0)
	for k, v := range s1 {
		ret = append(ret, k.(string)+"="+v.(string))
	}
	return ret

}

// func getVolumesFromSlice(s interface{}) []string {
// 	var (
// 		volumes     = make(map[string]struct{})
// 		hostVolumes = make([]string, 0)
// 	)
// 	for _, v := range s.([]interface{}) {
// 		vols := strings.Split(v.(string), ":")

// 	}
// }
