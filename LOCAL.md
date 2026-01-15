# USEFUL COMMANDS
- ### delete all top level files in z_log (ignore /dbg and /app)
    `
    rm $(find ./z_log -maxdepth 1  -type f)
    `
- tree command , -I flag to ignore files/dirs
    - **FIRST, LOAD THIS ENV VAR:**
    - `
export EXCTR='*.sql|z_*|*.md|sql|wiki|*.png|node*'
    `
        - `
tree -I $EXCTR
        `
    - save to file in z_dev dir:
        - `
tree -I $EXCTR > z_dev/tree.txt
        `
    - output to mac clipboard:
        - `
tree -I $EXCTR | pbcopy
        `

# PROD SERVER ADMIN
- ## MERGING IN NEW CODE
    - ### any changes or fixes must be developed in a separate branch and tested locally
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
    1. after confirming no issues with main running on prod server, sync local branch
        1. locally on branch mac, pull in changes from main
            - `
            git pull origin main
            `
        1. push those changs to mac remote branch
            -  `git push`
        1. pull mac remote to local
            - `git pull`
            - pushing to origin mac converts all prod urls to localhost:8080 for local dev


# local development
- ### developoment should be done in a non-main branch, no local work in mian
    - ssh into pi for work on main branch
- ### checkout repo: 
    -  `
    git clone https://github.com/jdetok/jdme.git
    `
    - `
    git checkout mac
    `
    - `
    git pull
    `