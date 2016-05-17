package parse

import (
	"errors"
	"fmt"
	"strings"

	"regexp"

	"github.com/crmonitor/pkg/crtype"
	"gopkg.in/yaml.v2"
)

// parse the compose yaml file
// support for version '2'

//normalization of project_name
//    def normalize_name(name):
//        return re.sub(r'[^a-z0-9]', '', name.lower())
//the details in config/main.py get_project_name

func normalization(name string) string {
	lowname := strings.ToLower(name)
	reg := regexp.MustCompile("[^a-z0-9]")
	normalizename := reg.ReplaceAll([]byte(lowname), []byte(""))
	return string(normalizename)
}

//a dir is created and there is a docker-compose file in it before excuting this function
func Composecheck(rawdata []byte) (map[interface{}]interface{}, error) {
	//transfer the yaml file into the config struct
	//only the first layer
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal(rawdata, &m)
	if err != nil {
		return nil, err
	}
	if m["version"] != "2" {
		return nil, errors.New("parse failed , only support v2 version of compose file")
	}
	fmt.Printf("value %+v", m)
	return m, nil
}

//transfer config map into Crmapp
//if there is not a specific image name in service
//use normalized(projectdir)_servicename to replace
func Fromconfigtoapp(projectname string, composemap map[interface{}]interface{}) (*crtype.CRMProject, error) {

	Crmapp := &crtype.CRMProject{Projectname: normalization(projectname)}
	Servicemap := composemap["services"]
	//get service info
	for key, value := range Servicemap.(map[interface{}]interface{}) {
		valuemap := value.(map[interface{}]interface{})
		imagename, ok := valuemap["image"]
		imagenamestr := ""
		if !ok {
			//create a new image name by projectdir : projectdir_servicename
			imagenamestr = Crmapp.Projectname + "_" + key.(string)
		} else {
			imagenamestr = imagename.(string)
		}

		//need to add latest automatically?
		if !strings.Contains(imagenamestr, ":") {
			imagenamestr = imagenamestr + ":latest"
		}
		Layer := crtype.Layer{Name: key.(string), Imagename: imagenamestr}
		Crmapp.Layers = append(Crmapp.Layers, Layer)
	}

	//get volume info

	//get network info

	return Crmapp, nil
}
