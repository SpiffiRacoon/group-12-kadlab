build:
    docker build -t kadlab .
    docker-compose up --build -d
up:
    docker-compose up --build -d

down:
    docker-compose down