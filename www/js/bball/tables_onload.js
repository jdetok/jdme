import { playerBtnListener } from "./ui.js"
import { table5f } from "./dynamic_table.js"; 
import { base, bytes_in_resp, FUSC_BOLD, GRN_BOLD } from "./util.js";


// RECENT GAMES TOP SCORERS TABLE FUNCS
export async function makeRGTopScorersTbl(data, numPl) {
    await table5f(data, 'tstbl', `Scoring Leaders | ${data.recent_games[0].game_date} Games`, 
        ["rank", "name | team", "matchup", "wl | score", "points"], numPl, rgTopScorerRow
    );
}

export async function rgTopScorerRow(tbl, data, i) {
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

// LEAGUE TOP SCORERS FUNCS
export async function makeScoringLeaders(numPl) {
    const url = `${base}/league/scoring-leaders?num=${numPl}`;
    const r = await fetch(url);
    if (!r.ok) {
        console.error(` error calling ${url}`);
    }
    console.trace(`%c ${await bytes_in_resp(r)} bytes received from ${url}}`, FUSC_BOLD)
    const data = await r.json();
    console.log(data);
    // await buildLeadingScorersTbl(data, 'top_lg_players', numPl);
    await table5f(data, 'nba_tstbl', 
        `Scoring Leaders | NBA/WNBA Top ${numPl}`, 
        ["rank", `nba | ${data.nba[0].season}`, "points", 
         `wnba | ${(data.wnba[0].season).substring(0, 4)}`, "points"],
        numPl, lgTopScorerRow,
    )
}

// add a row to the league top scorers table. called within a loop
// adds nba player with button, their points, wnba player with button, their points
export async function lgTopScorerRow(tbl, data, i) {
    let r = document.createElement('tr');

    let rank = document.createElement('td');
    let pName = document.createElement('td');
    let pts = document.createElement('td');
    let wpName = document.createElement('td');
    let wpts = document.createElement('td');
    let btn = document.createElement('button');
    let wbtn = document.createElement('button');

    rank.textContent = i + 1;
    btn.textContent = `${data.nba[i].player} | ${data.nba[i].team}`;
    btn.type = 'button';
    btn.addEventListener('click', async () => {
        await playerBtnListener(data.nba[i].player);
    }); 
    pName.appendChild(btn);

    wbtn.textContent = `${data.wnba[i].player} | ${data.wnba[i].team}`;
    wbtn.type = 'button';
    wbtn.addEventListener('click', async () => {
        await playerBtnListener(data.wnba[i].player);
    }); 
    wpName.appendChild(wbtn);

    pts.textContent = data.nba[i].points;
    wpts.textContent = data.wnba[i].points;

    r.appendChild(rank);
    r.appendChild(pName);
    r.appendChild(pts);
    r.appendChild(wpName);
    r.appendChild(wpts);

    tbl.appendChild(r);
}

// CURRENT TEAM RECORDS FUNCS
export async function getTeamRecords() {
    const url = `${base}/teamrecs`;
    try {
        const r = await fetch(`${base}/teamrecs`);
        if (!r.ok) {
            console.error(`failed to get data from ${url}`);
        }
        console.trace(`%c request status ${r.status}: ${await bytes_in_resp(r)} bytes received from ${url}}`, GRN_BOLD)
        const data = await r.json();
        return data;
    } catch(err) {
        console.error(`error fetching from ${url}: ${err}`);
    }
}

export async function makeTeamRecsTable(numRecs) {
    const data = await getTeamRecords();
    await table5f(data, 'trtbl', `NBA/WNBA Regular Season Team Records`,
        ["rank", `nba | ${data.nba_team_records[0].season}`, "record", 
            `wnba | ${data.wnba_team_records[0].season}`, "record"],
        numRecs, teamRecsRow)
}
// add a row to the league top scorers table. called within a loop
// adds nba player with button, their points, wnba player with button, their points
export async function teamRecsRow(tbl, data, i) {
    const nba = data.nba_team_records[i];
    const wnba = data.wnba_team_records[i];

    let r = document.createElement('tr');
    let rank = document.createElement('td');
    let tName = document.createElement('td');
    let rec = document.createElement('td');
    let wtName = document.createElement('td');
    let wrec = document.createElement('td');

    rank.textContent = i + 1;
    tName.textContent = nba.team_long;
    wtName.textContent = wnba.team_long;
    rec.textContent = `${nba.wins}-${nba.losses}`;
    wrec.textContent = `${wnba.wins}-${wnba.losses}`;

    r.appendChild(rank);
    r.appendChild(tName);
    r.appendChild(rec);
    r.appendChild(wtName);
    r.appendChild(wrec);

    tbl.appendChild(r);
}