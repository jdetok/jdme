# root files
- ### .air.toml
    - used for hot reloading - listens to the repo and builds/runs app with any changes
- ### .env
    - app requires a .env file with environment variables used in the project to function
- ### Dockerfile
    - the Dockerfile configures the container the api runs in. compose.yaml is in another repo; it's composed with the apache web server proxy
# directories
- ## api
    - app entrypoint, main.go lives here
    - all http routing/handlers
    - nba.go contains most of the handlers that server /bball
- ## static
    - all HTML/CSS/Javascript files live here
- ## internal
    resources in Go
    - ### mariadb
        - initialzes connection to mariadb server
    - ### store
        - response data structures are built here. uses mariadb package to query db, then scans rows into structs and returns the struct as marshaled JSON
    - ### env, errs
        - resources for standardized error handling/fetching .env vars
- ## scripts
    - sql and sh scripts