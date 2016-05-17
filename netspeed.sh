#!/bin/bash
set -e
avgtime=0

#resource="https://static.chachazhan.com/linux-4.5.4.tar.xz"
resource="https://static.chachazhan.com/linux-4.5.4.tar.xz"
echo "resource:"${resource}" 84.3M"
echo "node : 10.10.196.242"
echo "test time:"$(date)
for ((i=1;i<=30;i++));
do
        timea=$(date +%s)
        #echo ${timea}
        curl -o /dev/null ${resource} &> /dev/null
        #sleep 1
        timeb=$(date +%s)
        #echo ${timeb}
        echo "num "${i}" times for testing"
        diff=$(($timeb-$timea))
        echo "time using :" ${diff}s
        #avgtime=$(((${avgtime}+${diff}))/${i})
        #echo ${avgtime}
done
