# jdeko.me master directory
### visit jdeko.me: https://jdeko.me
## SUBMODULES:
** NOTE: the /sh directory will contain references to direcotries that are not found in this remote repo. these scripts run locally on the server in a directory with symbolic links to the subdirectories referenced here as submodules. the corresponding names are seen below: 
- ## go-api-jdeko.me
	production go api & frontend
	- **referenced  in /sh as `prod-api` 
- ## dev-jdeko.me
	- development go api & frontend
	- **referenced  in /sh as `dev-api` 
- ## docker-jdeko.me
	- apache web server/go api configs for docker containers
	- **referenced  in /sh as `dkr-env` 
- ## bball-etl-cli
	- go app to fetch nba/wnba data & insert into postgres
	- **referenced  in /sh as `etl-cli` 
