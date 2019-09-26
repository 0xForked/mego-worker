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
      
    rebuild with docker compose
    
     ``docker-compose up --build``