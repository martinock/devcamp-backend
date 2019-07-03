# Simple Golang Apps with PostgreSQL

Prerequisites:

1. Docker

How to run:

1. Run `sudo chmod +x wait-for-it.sh`
2. Run `sudo docker-compose up`
3. Populate docker postgres with data. dbname: postgres, table: users, user: postgres, password: postgres
NB: users contain two column: id (int) and name (varchar)
4. Open your browser and type `localhost:3000`. Route `/` is just a welcome page, next route is `/user/:user_id`
