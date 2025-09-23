import { tblCaption } from "./table.js"
import { base } from "./listen.js"
import {makePlayerDash} from "./player_dash.js"
import { playerBtnListener } from "./player_search.js"

/*
get the top scorer from each game from the most recent night where games occured
(usually dated yesterday, but when no games occur it'll get the most recent day
where games did occur). called on page load, it creates a table with all these
scorers and immediately grabs and loads the player dash for the top overall 
scorer. use season id 88888 in getP to get most recent season
*/
export async function makeRGTopScorers() {
    const r = await fetch(`${base}/games/recent`);
    if (!r.ok) {
        throw new Error(`${r.status}: error calling /games/recent`);
    }
    const data = await r.json();

    // capture overall leading scorer's ID, pass get the player and display 
    const top_scorer = data.top_scorers[0].player_id;
    
    // build the recent games top scorers 
    await buildRGTopScorersTbl(data, 'top_players');

    await makePlayerDash(base, top_scorer, 88888, 0, data);
}


/* 
creates a table with the top scorer from each game from the recent night. each 
player's name is a button that will search the player when clicked
*/
export async function buildRGTopScorersTbl(data, elName) {
    const tblcont = document.getElementById(elName);
    const tbl = document.getElementById('tstbl');
    const thead = document.createElement('thead');
    const nameH = document.createElement('td');
    const ptsH = document.createElement('td');
    const teamH = document.createElement('td');

    // set table column headers
    nameH.textContent = 'name';
    ptsH.textContent = 'points';
    teamH.textContent = 'team | matchup | win-loss';
    
    // set table caption
    const caption = `Scoring Leaders | ${data.recent_games[0].game_date} Games`
    tblCaption(tbl, caption);

    // append headers to head, head to table
    thead.appendChild(nameH);
    thead.appendChild(teamH);
    thead.appendChild(ptsH);
    tbl.appendChild(thead);
    
    /* loop through number of top scorers (should be number of games from that day)
    rgTopScorerRow adds a row to the table with player's stats/team/name etc*/
    for (let i = 0; i < data.top_scorers.length; i++) {
        await rgTopScorerRow(tbl, data.top_scorers[i], data);
    }
    // append the table element to the container (passed elName)
    tblcont.appendChild(tbl);
}

/*
called within the loop in recGmsTopScorersTbl - creates the rows for the table
with the data from each player
*/
export async function rgTopScorerRow(tbl, scorer, data) {
    let game = data.recent_games.find(g => g.player_id === scorer.player_id);
    let r = document.createElement('tr');

    let pName = document.createElement('td');
    let pTeam = document.createElement('td');
    let pts = document.createElement('td');

    let btn = document.createElement('button');
    btn.textContent = scorer.player;
    btn.type = 'button';
    btn.addEventListener('click', async () => {
        await playerBtnListener(scorer.player);
    }); 
    pName.appendChild(btn);
    
    pTeam.textContent = `${game ? game.team_name : ""} | \
    ${game ? game.matchup : ""} | ${game ? game.wl : ""}`;
    
    pts.textContent = scorer.points;

    r.appendChild(pName);
    r.appendChild(pTeam);
    r.appendChild(pts);

    tbl.appendChild(r);
}