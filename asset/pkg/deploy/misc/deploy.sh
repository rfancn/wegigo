#!/bin/bash
echo "===1: Begin kube installation"

# following tags make sure kube are installed on all nodes
ansible-playbook -i hosts kube_setup.yaml --tags "ensure_kube"
echo "===20:Install kubernetes on all nodes==="

# run setup master role
ansible-playbook -i hosts kube_setup.yaml --tags "setup_master"
echo "===40:Setup kubernetes master node==="

# get master info: ip/port/token
MASTER_INFO=$(ansible-playbook -i hosts kube_setup.yaml --tags "get_master_info"| grep "msg" | head -1)
MASTER_CONNECTION=$(echo $MASTER_INFO | cut -f2 -d "|")
MASTER_TOKEN=$(echo $MASTER_INFO | cut -f3 -d "|")
MASTER_IP=$(echo $MASTER_CONNECTION | cut -f1 -d ":")
MASTER_PORT=$(echo $MASTER_CONNECTION | cut -f2 -d ":")

echo $MASTER_TOKEN
echo $MASTER_IP
echo $MASTER_PORT

ansible-playbook -i hosts kube_setup.yaml --tags "setup_node" --extra-vars "kube_token=$MASTER_TOKEN kube_master_ip=$MASTER_IP kube_master_port=$MASTER_PORT"
echo "===60:Setup kubernetes worker node==="
