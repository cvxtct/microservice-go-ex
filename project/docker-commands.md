docker build -f logger-service.dockerfile -t tfosorcim/logger-service:1.0.0 .

docker push tfosorcim/logger-service:1.0.0
docker swarm join --token SWMTKN-1-3hlfd1nv6h1dnmkyh3kjjuhqu58obzco79auns82awqfnhb0f8-cpfpswplkdrezcycbv9zqxija 192.168.65.3:2377
## docker swarm
docker swarm init
docker swarm join-token manager
docker swarm join-token worker
docker stack deploy -c swarm.yml myapp
docker service ls
 docker service scale myapp_listener-service=3
update the swarm: docker service scale myapp_logger-serice=2
update the swarm: docker service update --image tfosorcim/logger-service:1.0.1 myapp_logger-service
stop docker swarm: docker swarm leave --force
