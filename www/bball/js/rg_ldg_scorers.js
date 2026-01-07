import { tblCaption } from "./table.js"
import { playerBtnListener } from "./player_search.js"
import { table5f } from "./dynamic_table.js"; 
export async function makeRGTopScorersTbl(data, numPl) {
    await table5f(data, 'tstbl', `Scoring Leaders | ${data.recent_games[0].game_date} Games`, 
        ["rank", "name | team", "matchup", "wl | score", "points"], numPl, rgTopScorerRowNew
    );
}

/* 
creates a table with the top scorer from each game from the recent night. each 
player's name is a button that will search the player when clicked
*/
export async function buildRGTopScorersTbl(data, elName) {
    // const tblcont = document.getElementById(elName);
    const tbl = document.getElementById('tstbl');
    const thead = document.createElement('thead');
    const nameH = document.createElement('td');
    const ptsH = document.createElement('td');
    const teamH = document.createElement('td');

    // set table column headers
    nameH.textContent = 'name | team';
    ptsH.textContent = 'points';
    teamH.textContent = 'matchup | win-loss';
    
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
    let tlen = 5;
    let numR = (data.top_scorers.length > tlen) ? tlen : data.top_scorers.length;
    for (let i = 0; i < numR; i++) {
        await rgTopScorerRow(tbl, data.top_scorers[i], data);
    }
}

export async function rgTopScorerRowNew(tbl, data, i) {
    let r = document.createElement('tr');

    let rank = document.createElement('td');
    let pName = document.createElement('td');
    let matchup = document.createElement('td');
    let wlscore = document.createElement('td');
    let pts = document.createElement('td');
    let btn = document.createElement('button');

    let player = data.top_scorers[i];
    let game = data.recent_games.find(g => g.player_id === player.player_id);
    

    rank.textContent = i + 1;

    btn.textContent = `${player.player} | ${game.team}`;
    btn.type = 'button';
    btn.addEventListener('click', async () => {
        await playerBtnListener(player.player);
    }); 
    pName.appendChild(btn);

    matchup.textContent = game.matchup;
    wlscore.textContent = `${game.wl} | ${game.points}-${game.opp_points}`
    pts.textContent = player.points;

    r.appendChild(rank);
    r.appendChild(pName);
    r.appendChild(matchup);
    r.appendChild(wlscore);
    r.appendChild(pts);

    tbl.appendChild(r);
}

// called within the loop in recGmsTopScorersTbl - creates the rows for the table
// with the data from each player
// run .find on recent_games.player_id & scorer.player_id to access team/game data
export async function rgTopScorerRow(tbl, scorer, data) {
    let game = data.recent_games.find(g => g.player_id === scorer.player_id);
    let r = document.createElement('tr');

    let pName = document.createElement('td');
    let pTeam = document.createElement('td');
    let pts = document.createElement('td');
    
    let btn = document.createElement('button');
    btn.textContent = `${scorer.player} | ${game.team}`;
    console.log(scorer);
    console.log(game);
    btn.type = 'button';
    btn.addEventListener('click', async () => {
        await playerBtnListener(scorer.player);
    }); 
    pName.appendChild(btn);
    
    pTeam.textContent = `${game.matchup} | ${game.wl} \
    | ${game.points}-${game.opp_points}`;
    
    pts.textContent = scorer.points;

    r.appendChild(pName);
    r.appendChild(pTeam);
    r.appendChild(pts);

    tbl.appendChild(r);
}