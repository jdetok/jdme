import * as table from "./table.js"
import { base } from "./listen.js"
import * as resp from "./resp.js"

/* 
adds a button listener to each individual player button in the leading scorers
tables. have to create a button, do btn.AddEventListener, and call this function
within that listener. will insert the player's name in the search bar and call 
getP
*/
export async function playerBtnListener(player) {
    let searchB = document.getElementById('pSearch');
    if (searchB) {
        searchB.value = player;
        searchB.focus();
        await resp.getP(base, player, 88888, 0, 0);
        searchB.value = '';
    }
}

/*
get the top scorer from each game from the most recent night where games occured
(usually dated yesterday, but when no games occur it'll get the most recent day
where games did occur). called on page load, it creates a table with all these
scorers and immediately grabs and loads the player dash for the top overall 
scorer. use season id 88888 in getP to get most recent season
*/
export async function getRecentTopScorers() {
    const r = await fetch(`${base}/games/recent`);
    if (!r.ok) {
        throw new Error(`${r.status}: error calling /games/recent`);
    }
    const data = await r.json();
    const top_scorer = data.top_scorers[0].player_id;

    const plrs = data.top_scorers.slice(1);
    console.log(plrs);
    console.log("FULL DATA: ");
    console.log(data);
    
    await recGmsTopScorersTbl(data, 'top_players',
        `Top Scorers from ${data.recent_games[0].game_date}`, 
        0, 'recent');

    await resp.getP(base, top_scorer, 88888, 0, data);
}

/* 
creates a table with the top scorer from each game from the recent night. each 
player's name is a button that will search the player when clicked
*/
export async function recGmsTopScorersTbl(data, elName, caption) {
    const tblcont = document.getElementById(elName);
    const tbl = document.getElementById('tstbl');
    table.tblCaption(tbl, caption);

    // table headers
    const thead = document.createElement('thead');

    const nameH = document.createElement('td');
    nameH.textContent = 'name';

    const ptsH = document.createElement('td');
    ptsH.textContent = 'points';

    const teamH = document.createElement('td');
    
    teamH.textContent = 'team | matchup | win-loss';
    

    thead.appendChild(nameH);
    thead.appendChild(teamH);
    thead.appendChild(ptsH);

    tbl.appendChild(thead);
    
    let scorers = data.top_scorers;
    
    for (let i = 0; i < scorers.length; i++) {
        await rgTopScorerRow(tbl, scorers[i], data);
    }
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

// build top x players table
export async function getScoringLeaders() {
    let numPl = 10;
    const r = await fetch(`${base}/league/scoring-leaders?num=${numPl}`);
    if (!r.ok) {
        throw new Error(`${r.status}: error calling /league/scoring-leaders`)
    }
    const data = await r.json();
    console.log(data);
    await leadingScorersTbl(data, 'top_lg_players', numPl)
}

export async function leadingScorersTbl(data, elName, numPl) {
    const tblcont = document.getElementById(elName);
    const tbl = document.getElementById('nba_tstbl');

    const caption = `Top ${numPl} NBA/WNBA Scorers`
    table.tblCaption(tbl, caption);
    
    const thead = document.createElement('thead');

    const nbaH = document.createElement('td');
    nbaH.textContent = `nba | ${data.nba[0].season}`;

    const ptsH = document.createElement('td');
    ptsH.textContent = 'points';

    const wnbaH = document.createElement('td');
    wnbaH.textContent = `wnba | ${data.wnba[0].season}`;

    const wptsH = document.createElement('td');
    wptsH.textContent = 'points';

    thead.appendChild(nbaH);
    thead.appendChild(ptsH);
    thead.appendChild(wnbaH);
    thead.appendChild(wptsH);

    tbl.appendChild(thead);
    
    for (let i = 0; i < numPl; i++) {
        await lgTopScorerRow(tbl, data.nba[i], data.wnba[i]);
    }

    tblcont.appendChild(tbl);
}

/* 
add a row to the league top scorers table. called within a loop
adds nba player with button, their points, wnba player with button, their points
*/
export async function lgTopScorerRow(tbl, nba, wnba) {
    let r = document.createElement('tr');

    let pName = document.createElement('td');
    let pts = document.createElement('td');
    let wpName = document.createElement('td');
    let wpts = document.createElement('td');

    let btn = document.createElement('button');
    btn.textContent = `${nba.player} | ${nba.team}`;
    btn.type = 'button';
    btn.addEventListener('click', async () => {
        await playerBtnListener(nba.player);
    }); 

    let wbtn = document.createElement('button');
    wbtn.textContent = `${wnba.player} | ${wnba.team}`;
    wbtn.type = 'button';

    wbtn.addEventListener('click', async () => {
        await playerBtnListener(wnba.player);
    }); 

    pName.appendChild(btn);
    wpName.appendChild(wbtn);

    pts.textContent = nba.points;
    wpts.textContent = wnba.points;
    r.appendChild(pName);
    r.appendChild(pts);
    r.appendChild(wpName);
    r.appendChild(wpts);

    tbl.appendChild(r);
}