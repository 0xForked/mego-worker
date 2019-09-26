## MEGO 
SMS (Short Messaging Service) Worker that handle incoming queue from Message Broker and send it with gammu - part of next project (Microservices Arch Notify-Service)

### Installation Instructions
1. Run git clone https://github.com/aasumitro/mego-worker.git mego_worker
2. From the projects root run

    ``
    cp .env.example .env
    ``
 3. Configure your ``.env`` file
 4. Run with docker-compose 
 
    ``docker-compose up``
    
    as daemon
    
      ``docker-compose up -d``
      
    build with docker compose
    
     ``docker-compose up --build``
     
5. Run with docker 
    
    build as image
    
    ``docker build -t {new_image_name} .`` e.g ``docker build -t mego-worker .``
    
    show all image
    
    ``docker images -a`` or ``docker image ls``
    
    run image as container
    
    ``docker run -d --name {new_container_name} {image_name}`` e.g ``docker run -d --name mego_app mego-worker``
    
    show all container 
    
    ``docker ps -a``
    
    show running container
    
    `` docker ps``
    
    docker container mgmt
    
    ``docker container {start|stop|restart} {container_name}``
    
    docker stats 
    
    ``docker stats --all {container_name}``
    
    open container bash
    
    ``docker exec -it {container_name} /bin/bash``