package clienttool

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/crmonitor/cmd/cragent/conf"
	influxpackage "github.com/influxdata/influxdb/client/v2"
)

var InfluxClient *InfluxdbStorage

type InfluxData struct {
	Measurement string
	Fields      map[string]interface{}
	Tags        map[string]string
}

type InfluxdbStorage struct {
	Influxclient    influxpackage.Client
	MachineName     string
	Database        string
	RetentionPolicy string
	BufferDuration  time.Duration
	LastWrite       time.Time
	//points          []*influxdb.Point
	Lock         sync.Mutex
	ReadyToFlush func() bool
}

func queryDB(clnt influxpackage.Client, cmd string, dbname string) (res []influxpackage.Result, err error) {
	q := influxpackage.Query{
		Command:  cmd,
		Database: dbname,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func GetinfluxClient(dbname string) (*InfluxdbStorage, error) {
	//influxserver string, username string, password string, dbname string
	influxserver := conf.GlobalConfig.DefaultInfluxURL

	username := ""
	password := ""
	if InfluxClient != nil && InfluxClient.Database == dbname {
		return InfluxClient, nil
	}
	client, err := NewinfluxClient(influxserver, username, password, dbname)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewinfluxClient(influxserver string, username string, password string, dbname string) (*InfluxdbStorage, error) {
	c, err := influxpackage.NewHTTPClient(influxpackage.HTTPConfig{
		Addr:     influxserver,
		Username: username,
		Password: password,
	})

	if err != nil {
		return nil, err
	}

	_, err = queryDB(c, fmt.Sprintf("CREATE DATABASE %s", dbname), dbname)
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	influxc := &InfluxdbStorage{
		Influxclient: c,
		MachineName:  hostname,
		Database:     dbname,
		LastWrite:    time.Now(),
	}

	return influxc, nil
}

func (self *InfluxdbStorage) AddStats(influxDataList []*InfluxData) error {

	influxclient := self.Influxclient

	bp, _ := influxpackage.NewBatchPoints(influxpackage.BatchPointsConfig{
		Database:  self.Database,
		Precision: "us",
	})
	for _, value := range influxDataList {

		fmt.Println("**the measurement**", value.Measurement)
		point, err := influxpackage.NewPoint(value.Measurement, value.Tags, value.Fields, time.Now())
		if err != nil {
			return err
		}
		fmt.Println("the point name:", point.Name())
		bp.AddPoint(point)
	}
	influxclient.Write(bp)

	return nil
}

func (self *InfluxdbStorage) Close() error {
	err := self.Influxclient.Close()
	if err != nil {
		return err
	}
	return nil
}
