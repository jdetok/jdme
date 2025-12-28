# TODO: Winter Break

- ## ETL work
    - less golib dependency
    - run for specific date
    - add contexts and timeouts
    - log errors somewhere
    - if requests fail, store, move on, and try again
    - maybe look at mlb code i was doing

- ## DB work
    - explore new schema design
    - standardize column names across tables
        - some tables use plr_id, some player_id etc

- ## serve static files with nginx
    - rename static to /www and mount to volume
    - no longer serve static pages from go
    - static error pages

- ## clean up environment variables
    - each container has .env variables loaded when the container is build, don't need to load vars explicitly in the code
    - no more repetive vars
    - for postgres, use PGPASSWORD

- ## better postgres admin
    - different users for each process

- ## automated log compression
    - after n days compress logs, after n days delete the compressed logs

- ## serve logs as dir with nginx

- ## create an image for serving /wiki
    - server via nginx
    - container should run hugo to build the site