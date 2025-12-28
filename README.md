# jdeko.me source code
- ## written & maintained by Justin DeKock 
    jdeko17@gmail.com | [jdetok on github](https://github.com/jdetok)
    ### visit the site: 
    - [jdeko.me home page](https://jdeko.me)
    - [nba/wnba stats page](https://jdeko.me/bball)
    - [main about page](https://jdeko.me/about)
    - [nba/wnba stats about page](https://jdeko.me/bball/about)

# Summary
[jdeko.me](https://jdeko.me) is a personal website designed & rewritten primarily to gain real-world experience with backend development, systems design/architecture, and database design/management.

# Environment (Docker)
every component of jdeko.me is containerized via Docker and orchestrated via a single [docker compose file](./jdme-dkr/compose.yaml). this project builds several docker images, each of which is defined in a Dockerfile within the [jdme-dkr directory](./jdme-dkr/) 

# API (Go)
http requests to jdeko.me are handled by the API I wrote primarily with Go's standard library. jdeko.me/bball is served by a Postgres database I built to store NBA/WNBA statistcs. 

# Storage (Postgres, MongoDB)
a custom-designed postgres database stores the nba/wnba data returned by the API. new data is fetched nightly after the conclusion of all games using the [bball-etl-go](https://github.com/jdetok/bball-etl-go) package I wrote to fetch NBA/WNBA stats from stats.nba.com. the database was built using the same package to fetch data from every NBA/WNBA game since 1970. go's concurrency features enable sourcing and inserting 50+ seasons of data (10+ GET requests to nba.com per season, ~ 1.6 million rows inserted on build as of 8/6/2025) in ~15 minutes.

# Frontend (HTML/CSS/JavaScript)
the site's frontend is written using only pure HTML/CSS/Javascript. all static files are served from /static by the handlers in /api/static.go. /static/bball are the current production files that serve jdeko.me/bball
![Alt Text](https://jdeko.me/img/jdekome_ex_080825.png "main example")
![Alt Text](https://jdeko.me/img/bball_ex_121025.png "/bball example")

# ** CHANGELOG / REVISION HISTORY
- ## original deployment: early July 2025
- ## early august 2025
    - redesigned original MariaDB database as a Postgres DB
- ## early september 2025
    - 9/10/2025 added top scorers interactive box
    - 9/12/2025 created documentation with go doc
    - moved api to its own package, isolated main
- ## october-early december 2025
    - redesign in-mmemory storage system using maps rather than struct slices
    - add mongodb server container for logging http requests
- ## late december 2025
    - ## started CI/CD pipeline
        - using github actions
        - build go code
        - python script to make changes to files on push/pull based on branch
    - ## refactor environment
        - all containers in a single compose.yaml file
        - switched from apache to nginx for reverse proxy server
        - db containers/stats fetch containers now included in same compose
        - db backups