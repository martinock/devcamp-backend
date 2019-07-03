Simple Golang Apps with PostgreSQL

Prerequisites:
1. Docker

How to run:
1. Run `sudo chmod +x wait-for-it.sh`
2. Populate docker postgres with data. dbname: postgres, table: users, user: postgres, password: postgres
NB: users contain two column: id (int) and name (varchar)
3. Run `docker-compose -f docker-compose.yml up`
4. Open your browser and type `localhost:3000`. Route "/" is just a welcome page, next route is `/hello/:user_id`
