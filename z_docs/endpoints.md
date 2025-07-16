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

## /games/recent
returns all games from the most recent game day (usually yesterday), along with the top scorer of each game. /player is called on the overall top scorerr from the day and their player dash is build on page load