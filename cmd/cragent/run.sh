#! /bin/bash


export INTERVALTIME=3
export DOCKERENDPOINT="unix:///var/run/docker.sock"
export INFLUXDBURL="http://127.0.0.1:8086"
export INFLUXDBNAME="monitor"


./cragent

