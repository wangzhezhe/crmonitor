package api

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"io/ioutil"
	"time"

	"github.com/crmonitor/pkg/crmaster/manager"
	"github.com/crmonitor/pkg/util/parse"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
)

var CRMasterregisterpath = "composedir"

func Getengine() *gin.Engine {
	r := gin.Default()
	// Apply the middleware to the router (works with groups too)
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	return r

}

func Loadcragentapi(e *gin.Engine) *gin.Engine {
	e.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "started",
		})
	})
	return e
}

func getCRMastermanager() (*manager.CRMastermanager, error) {
	crmanager, err := manager.GetCRMastermanager(manager.Defaultetcdurl)
	if err != nil {
		log.Println(err)
	}
	log.Printf("the crmanager %+v", crmanager)
	return crmanager, nil

}

func writesocket(event string) {

	//ClientMap["hack"].Write(fmt.Sprintf("container %s: %s", value.ID, value.Status))
	//ClientMap["hack"] = client
}

func socket(c *gin.Context) {
	//var ClientMap map[string]*gateway.Client

}

func doregister(c *gin.Context) {

	file, header, err := c.Request.FormFile("composefile")
	projectname := c.Request.FormValue("projectname")

	//check project name
	if err != nil {
		c.JSON(500, gin.H{"message": "do not receive file"})
		return
	}
	if projectname == "" {
		c.JSON(500, gin.H{"message": "do not get the projectname"})
		return
	}
	log.Println("get project name", projectname)
	filename := header.Filename
	log.Println("the header.Filename", header.Filename)
	//create the dir with projectname
	err = os.MkdirAll("./"+CRMasterregisterpath+"/"+projectname, 0777)
	if err != nil {
		c.JSON(500, gin.H{"message": "failed to create the dir"})
		return
	}
	filepath := "./" + CRMasterregisterpath + "/" + projectname + "/" + filename
	out, err := os.Create(filepath)
	if err != nil {
		c.JSON(500, gin.H{"message": "failed to create file"})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(500, gin.H{"message": "failed to copy the file"})
		return
	}

	log.Println("uploaded ok")
	rawdata, err := ioutil.ReadFile(filepath)
	if err != nil {
		c.JSON(500, gin.H{"message": "failed to read file"})
		return
	}

	if err != nil {
		c.JSON(500, gin.H{"message": "failed to red json file"})
		return
	}

	crmanager, err := manager.GetCRMastermanager(manager.Defaultetcdurl)
	if err != nil {
		log.Println("failed to ge the manager", err)
		c.JSON(500, gin.H{"message": "failed to get the manager"})
		return

	}

	//change raw data to project instance
	//projectname string, composemap map[interface{}]interface{}) (*crtype.CRMProject
	rawdatamap, err := parse.Composecheck(rawdata)
	if err != nil {

	}
	projectinstance, err := parse.Fromconfigtoapp(projectname, rawdatamap)
	if err != nil {
		c.JSON(500, gin.H{"message": "failed to change the config into the appinstance"})
	}
	rawdata_instance, _ := json.Marshal(projectinstance)
	err = crmanager.Registerproject(projectname, string(rawdata_instance))
	if err != nil {
		c.JSON(500, gin.H{"message": "failed to register the app"})
		return
	}

	c.JSON(200, gin.H{"message": "uploded ok"})
	return
}

func getproject(c *gin.Context) {
	crmanager, err := manager.GetCRMastermanager(manager.Defaultetcdurl)
	if err != nil {
		log.Println("failed to ge the manager", err)
		c.JSON(500, gin.H{"message": "failed to get the manager"})
		return

	}
	applist, err := crmanager.Getproject()
	if err != nil {
		log.Println("error , failed to get the project list ", err)
		c.JSON(500, gin.H{"message": "failed to get the project list"})
		return
	}
	c.JSON(200, applist)
}

func getimages(c *gin.Context) {

}

func getcontainersinproject(c *gin.Context) {

}

func Loadcrmasterapi(e *gin.Engine) *gin.Engine {
	e.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "started",
		})
	})
	e.POST("crmonitor/register", doregister)

	e.GET("crmonitor/socket", socket)

	e.GET("crmonitor/project", getproject)

	e.GET("crmonitor/images", getimages)

	e.GET("crmonitor/project/images/containers", getcontainersinproject)

	return e
}
