#!/bin/bash
CONTAINERNAME=natsstreaming_gonatseventsourcing

mkdir -pv `pwd`/environments/dev/datastore  
docker stop $CONTAINERNAME || true && docker rm $CONTAINERNAME
docker run --name $CONTAINERNAME -p 4222:4222 -p 8222 -v `pwd`/environments/dev/datastore:/datastore -d nats-streaming -store file -dir datastore --cluster_id gonatseventsourcing_cluster --max_msgs 0 --max_bytes 0 -max_channels 0 -m 8222

#cd workspace
#fresh
