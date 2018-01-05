#!/bin/bash

echo "===1:Begin kube installation==="

### Install kube ###
ansible-playbook -i hosts kube_setup.yaml --tags "ensure_kube"
echo "===20:Install kubernetes on all nodes==="

### Setup kube master ###
# run setup master role
ansible-playbook -i hosts kube_setup.yaml --tags "setup_master"
echo "===40:Setup kubernetes master nodes==="

### Setup kube node ###
# run setup worker node role
# get master info: ip/port/token
MASTER_INFO=$(ansible-playbook -i hosts kube_setup.yaml --tags "get_master_info"| grep "msg" | head -1)
MASTER_CONNECTION=$(echo $MASTER_INFO | cut -f2 -d "|")
MASTER_TOKEN=$(echo $MASTER_INFO | cut -f3 -d "|")
MASTER_IP=$(echo $MASTER_CONNECTION | cut -f1 -d ":")
MASTER_PORT=$(echo $MASTER_CONNECTION | cut -f2 -d ":")

ansible-playbook -i hosts kube_setup.yaml --tags "setup_node" --extra-vars "kube_token=$MASTER_TOKEN kube_master_ip=$MASTER_IP kube_master_port=$MASTER_PORT"

# wait for all nodes get ready
ansible-playbook -i hosts kube_setup.yaml --tags "wait_for_nodes_ready"
echo "===50:Setup kubernetes worker nodes==="

### Deploy rabbitmq ###
ansible-playbook -i hosts kube_setup.yaml --tags "deploy_rabbitmq"

# wait for rabbitmq get ready
ansible-playbook -i hosts kube_setup.yaml --tags "wait_for_rabbitmq_ready"
echo "===60:Deploy RabbitMQ MessageBroker==="
