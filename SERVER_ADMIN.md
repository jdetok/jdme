# notes for local server administration
# MERGING IN NEW CODE
- ### any changes or fixes must be developed in a separate branch and tested locally
1. #### from local machine
    1. after testing locally, push to branch origin
    1. start a pull request to the main branch
    1. make sure all actions pass
1. SSH into pi
1. cd into local repo
    1. stop go app's docker container
        - `
    docker compose down -v
    `
    1. merge in changes
        - `
        git pull --no-rebase
        `
    1. check over repo, make sure replacements in github actions worked as expected
        1. lingering localhost URLs in /static .html .css .js 
        1. incompatible CPU option for Go compiler in Dockerfile
            - must be `arm64` for raspberry pi, `amd64` when running locally on mac silicon
    1. rebuild docker container
        - `
        docker compose up --build
        `
# USEFUL COMMANDS
- ## delete all top level files in z_log (ignore /dbg and /app)
    `
    rm $(find ./z_log -maxdepth 1  -type f)
    `

