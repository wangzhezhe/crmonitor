package parse

import (
	"log"
	"testing"
)

var data = `version: '2'
services:
  web:
    build: .
    depends_on:
      - db
      - redis
  redis:
    image: redis
  db:
    image: postgres`

func TestComposecheck(t *testing.T) {
	//t.Skip()
	m, err := Composecheck([]byte(data))
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("return value:%+v", m)
	}

	log.Printf("services are %+v:", m["services"])

}

func TestNormalization(t *testing.T) {
	name := "abc_def/_ghi"
	value := normalization(name)
	log.Println(value)
}

func TestFromconfigtoapp(t_ *testing.T) {
	projectname := "test"
	m, err := Composecheck([]byte(data))
	if err != nil {
		log.Println(err)
		return
	}
	crmproject, err := Fromconfigtoapp(projectname, m)
	if err != nil {
		log.Println(err)
	}
	log.Printf("the crm project\n %+v", crmproject)

}
