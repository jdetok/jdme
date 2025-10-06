# jdeko.me api request flows
## HIGH LEVEL PLAYER SEARCH SCENARIOS
- on-load search
    - top scorer from most recent night of games is requested when page loads
- normal player search
    - user types a player's name into the search bar & submits
    - checks for checked season and team checkboxes
- random player search
    - user clicks "search random player" button
    - checks for checked season check box, default 88888 season id if not
    - checks league radio buttons
## PLAYER SEARCH / RANDOM PLAYER SEARCH DETAIL
- `searchPlayer()` is called when search button is clicked or ui form is otherwise submitted, `randPlayerBtn()` is called when the random player button is clicked
    - assign current user selections to variables
        - searchPlayer() assigns value entered in search bar to `player` variable
            - if nothing in search bar, filled with current displaying player's name
        - randPlayerBtn() assigns the word `random` to the player variable
        - calls lgRadioBtns() to get the value of the selected league selector and assign to `lg` variable
        - calls checkBoxGroupValue() for both season and team selectors to assign 
        `season` and `team` variables
    - uses these variables to call `getPlayerStats(base, player, season, team, lg)` and attempt to get JSON response from the API's /players endpoint
        - the API handles the request via the `HndlPlayer()` handler
            - parses the query string & assigns the passed player, season, and team as strings in the `PlayerQuery` struct
            - passes this to `ValidatePlayerSzn()`, which validates and returns the player, season, and team IDs as uint64s in the `PQueryIds` struct 
                - if the search is not valid (no name match, player didn't play in requested season/for requested team, etc.) the JSON response will include an error string, which populates the `err` HTML element
                - `RandomPlayerId()` is called from here if the player variable sent to the API is `random`
                    -  this first calls `SlicePlayersSzn()` to create a slice of players only from within the season passed 
                    - a random number less than the length of this slice is then used to return a randomly indexed player ID in this slice
            - if `PQueryIds` is successfully populated, `GetPlayerDash()` is called to query the database, scan the rows to various response structs, which are then marshalled into valid JSON and returned as a slice of bytes
        - the raw JSON produced by the API is finally passed to `app.JSONWriter()`, which sends the response to the client as JSON
    - if successful, passes JSON response to `buildPlayerDash()` which populates response elements on the site

## ON-LOAD FLOW
- several functions are called to load page/response data when jdeko.me/bball is loaded
    - SETUP FUNCTIONS: 
        - each setup function is called within static/bball/js/listen.js in a `DOMContentLoaded` event listener
        - SETUP UI ELEMENTS:
            - setupExclusiveCheckboxes()
                - pass two checkbox elements, makes sure only one can be pressed at a time
                - called twice: once for season checkboxes & once for team checkboxes
            - clearCheckBoxes()
                - unchecks all page checkboxes
            - lgRadioBtns()
                - sets up league radio buttons
            - loadSznOptions()
                - calls /seasons API endpoint & fills both season select elements with an option element for each season
                    - option text is season spelled out, option value is season ID
            - loadTeamOptions()
                - calls /teams API endpoint & fills both team select elements with an option element for each team
                    - option text is team spelled out, option value is team ID
    - LOAD STATS LEADERS TABLES AND RECENT GAMES TOP SCORER DASH
        - `getRecentGamesData()` calls /games/recent endpoint to get scores from the most recent night of games. the "Recent Games Top Scorers" table is populated with each of the top scorers in the response
            - the top overall scorer from these games is used to call getPlayerData() & buildPlayerDash()