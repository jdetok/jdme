# full stack repo for my personal website 
## main site: https://jdeko.me
![Alt Text](https://jdeko.me/img/ex2.png "main example")
## NBA/WNBA stats site: https://jdeko.me/bball
![Alt Text](https://jdeko.me/img/ex.png "/bball example")
### written & maintained by Justin DeKock (jdeko17@gmail.com)

# API (Go)
http requests to jdeko.me are handled by the API I wrote primarily with Go's standard library. jdeko.me/bball is served by a Mariadb database I built to store NBA/WNBA statistcs. 

# Frontend
the site's frontend is written using only pure HTML/CSS/Javascript. all static files are served from /static by the handlers in /api/static.go. /static/bball are the current production files that serve jdeko.me/bball

# Environment
every process that serves jdeko.me is containerized with docker. the entire environment is configured in a docker compose file. the compose file for this repo also builds and runs an apache web server for proxying. the mariadb server is in a separate compose, as I use it for other projects. there's an additional docker compose that runs my python script for fetching NBA/WNBA data and inserting it into the databae. 

# Database (Mariadb Server)
the nba/wnba data returned by the API is stored in a mariadb database running in a docker container. the data is inserted/updated nightly from a python project i wrote. the project utilizes the nba_api python package, which calls stats/nba.com. the python source code is also containerized - it runs via a shell script i wrote to build, compsoe up, run, & compose down. this project's repo can be seen here: https://github.com/jdetok/py-nba-mdb