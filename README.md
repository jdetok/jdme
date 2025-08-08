# jdeko.me | full stack repo
## written & maintained by Justin DeKock 
jdeko17@gmail.com | [jdetok on github](https://github.com/jdetok)
### visit jdeko.me: 
- https://jdeko.me
- https://jdeko.me/bball
- https://jdeko.me/about 

# API (Go)
http requests to jdeko.me are handled by the API I wrote primarily with Go's standard library. jdeko.me/bball is served by a Postgres database I built to store NBA/WNBA statistcs. 

# Postgres Database
the nba/wnba data returned by the API is stored in a postgres database running in a docker container. another go program - [nightly-bball-etl](https://github.com/jdetok/nightly-bball-etl.git), which utilizes the [bball-etl-go](https://github.com/jdetok/bball-etl-go) package I wrote to fetch NBA/WNBA stats from stats.nba.com & insert them into the postgres database, runs nightly in a cronjob to update the database with stats from the most recent day. the database is built from [bball-postgres-bld-etl](https://github.com/jdetok/bball-postgres-bld-etl), which also utilizes [bball-etl-go](https://github.com/jdetok/bball-etl-go) to get data from every NBA/WNBA game since 1970. go's concurrency features enable sourcing and inserting 50+ seasons of data (10+ GET requests to nba.com per season, ~ 1.6 million rows inserted on build as of 8/6/2025) in ~15 minutes.

# Frontend
the site's frontend is written using only pure HTML/CSS/Javascript. all static files are served from /static by the handlers in /api/static.go. /static/bball are the current production files that serve jdeko.me/bball
![Alt Text](https://jdeko.me/img/jdekome_ex_080825.png "main example")
![Alt Text](https://jdeko.me/img/bball_ex_080825.png "/bball example")

# other repos that power jdeko.me:
- ### [bball-etl-go](https://github.com/jdetok/bball-etl-go)
    - etl (extract-transform-load) package written in go to fetch NBA/WNBA data & load it into postgres database
    - ### [bball-postgres-bld-etl](https://github.com/jdetok/bball-postgres-bld-etl)
        - uses etl package to build postgres database | fetch & insert with NBA/WNBA data from 1970-current
    - ### [nightly-bball-etl](https://github.com/jdetok/nightly-bball-etl)
        - uses etl package to fetch NBA/WNBA data from previous day, insert into existing postgres database
- ### [docker-jdeko.me](https://github.com/jdetok/docker-jdeko.me)
    - docker compose, Dockerfile, & apache files for environment setup 
- ### [golib](https://github.com/jdetok/golib)
    - go packages i wrote to standardize logging/error handling/etc in all my projects