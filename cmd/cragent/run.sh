# the script is used for local env testing
#!/bin/bash
export INTERVALTIME="3"
export DOCKERENDPOINT="unix:///var/run/docker.sock"
export INFLUXDBURL="http://114.115.208.145:8086"
export INFLUXDBNAME_CONTAINER="monitorinfo"
export INFLUXDBNAME_PACKET="packetinfo"

export HOSTIP="10.211.55.5"
export INTERFACE="eth5"

./cragent

