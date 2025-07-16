# endpoint documentation for /bball

## /player
- ### ?player=
    - accepts a player name or player id
    - pass "random" to get a random player
- ### &season=
    - accepts a season id
    - season ids that start with 2 are regular season, 4 is playoffs
    - each player's min/max reg. season/playoff season ids are stored - passing a season > their max szn will return their actual max season, same with min season
    - 99999: career reg. season
    - 99998: career playoffs
    - 99997: career reg. season/playoffs
- ### &team=
    - accepts a team id
    - pass 0 to bypass team query