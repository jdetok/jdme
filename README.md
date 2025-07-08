# Full stack repo for my personl website 
## main site: https://jdeko.me
## NBA/WNBA stats site: https://jdeko.me/bball
### written & maintained by Justin DeKock (jdeko17@gmail.com)

# API (Go)
http requests to jdeko.me are handled by the API I wrote primarily with Go' standard library. jdeko.me/bball is served by a Mariadb database I built to store NBA/WNBA statistcs. 

# Environment
every process that serves jdeko.me is containerized with docker. the entire environment is configured in a docker compose file. the compose file for this repo also builds and runs an apache web server for proxying. the mariadb server is in a separate compose, as I use it for other projects. there's an additional docker compose that runs my python script for fetching NBA/WNBA data and inserting it into the databae. 
