# docker instructions for jdeko.me api

## containers in repo
- there are three containers in this repo
    - /proxy contains the apache proxy compose and dockerfile
        - handles reverse proxying
    - /mongo contains compose for MongoDB container where http logs are written
    - API CONTAINER
        - the compose and dockerfile at root compile and run the Go API
        - the proxy and mongo containers must be running to start the main container

## external containers
- the API relies on an external container running the postgres db with NBA stats. this runs on the same docker engine as the production API