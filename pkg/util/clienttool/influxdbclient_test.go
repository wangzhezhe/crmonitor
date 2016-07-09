package clienttool

import (
	"testing"
)

func TestAddStats(t *testing.T) {
	//influxserver string, username string, password string, dbname string
	//"http://127.0.0.1:8086", "wangzhe", "123456", "testb"
	client, err := GetinfluxClient("http://127.0.0.1:8086", "wangzhe", "123456", "testa")
	if err != nil {
		t.Error(err)
	}

	measurement := "testameasure"
	tags := map[string]string{"taga": "taga-valuea", "tagb": "tagb-valueb"}
	fields := map[string]interface{}{
		"valuea": 2,
		"valueb": 2,
		"valuec": 4,
	}

	dataa := &InfluxData{Measurement: measurement, Fields: fields, Tags: tags}
	datab := &InfluxData{Measurement: measurement + "datab", Fields: fields, Tags: tags}
	var influxList []*InfluxData
	influxList = append(influxList, dataa)
	influxList = append(influxList, datab)
	err = client.AddStats(influxList)
	if err != nil {
		t.Error(err)
	}

}
