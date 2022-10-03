docker build -f logger-service.dockerfile -t tfosorcim/logger-service:1.0.0 .

docker push tfosorcim/logger-service:1.0.0
docker swarm join --token SWMTKN-1-3hlfd1nv6h1dnmkyh3kjjuhqu58obzco79auns82awqfnhb0f8-cpfpswplkdrezcycbv9zqxija 192.168.65.3:2377
## docker swarm
docker swarm join-token manager
docker swarm join-token worker
docker stack deploy -c swarm.yml myapp
docker service ls
