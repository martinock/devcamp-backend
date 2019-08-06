# Simple Golang Apps with PostgreSQL

Prerequisites:

1. Docker

How to run:

1. Run `sudo chmod +x wait-for-it.sh`
2. Run `sudo docker-compose up`
3. Populate docker postgres with data. `cat devcamp.pgsql | docker exec -i <your-db-container> psql -U postgres`
4. Open your browser and type `localhost:3000`. Route `/` is just a welcome page, next route is `/user/:user_id`
