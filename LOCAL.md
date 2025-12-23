# local development
* ### developoment should be done in a non-main branch, no local work in mian
    - ssh into pi for work on main branch
### checkout repo: 
-  `
git clone https://github.com/jdetok/go-api-jdeko.me.git
`
- `
git checkout mac
`
- `
git pull
`

## build docker environment
- ### mongodb logging container
    - this is built first
    - in ./mongo
- ### api container
    - compose file: `./compose.yaml`
        - dev creates service `api_jdme` and network `jdme_net`
            - the network is provided by proxy for prod
        - connects to lognet
    - dockerfile: `./Dockerfile`
- ## need to add PGDB here