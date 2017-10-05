
## commands

    # scale child processes up/down
    sudo docker service scale swarm-demo-child=1

    # start dashboard
    sudo docker run -d -v /var/run/docker.sock:/var/run/docker.sock \
       -p 8080:8080 charypar/swarm-dashboard

## digital ocean cluster

    cd scripts
    
    # create cluster
    ./do-create-cluster.sh
    doctl compute droplet list
    
    # edit Supfile and set ip addresses
    
    # install docker on all hosts
    sup do install-docker
    
    # init swarm on manager node
    sup do-mgr init-swarm-manager
    
    # edit Supfile and set SWARM_TOKEN
    
    # have other nodes join swarm
    sup do join-swarm
    
    # deploy services
    sup do-mgr create-services
    
    # list services
    sup do-mgr swarm-info
    
    # create load balancer
    doctl compute load-balancer create \
       --name swarm-demo \
       --droplet-ids 62715141,62715497,62715760 \
       --region sfo2 \
       --forwarding-rules 'entry_protocol:tcp,entry_port:80,target_protocol:tcp,target_port:80' \
       --health-check 'protocol:http,port:80,path:/,check_interval_seconds:10'
    
    # list load balancer
    doctl compute load-balancer list
    
    # should now be able to go to load balancer's ip in a browser and see demo

## deploying with stacks

    cd services
    sudo docker stack deploy --compose-file demo-stack.yml swarm-demo
    
    sudo docker stack deploy --compose-file visualizer-stack.yml visualizer


## visualizer

To see what's running on each node in the swarm:

    sup managers start-visualizer

## benchmark

    # wrk
    wrk -c 10 -t 2 -d 30s  http://localhost:9000/fragment/child
