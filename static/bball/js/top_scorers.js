import * as table from "./table.js"
import { base } from "./listen.js"
import * as resp from "./resp.js"

export async function playerLinkSearch(player) {
    let searchB = document.getElementById('pSearch');
    if (searchB) {
        searchB.value = player;
        searchB.focus();
        await resp.getP(base, player, 88888, 0, 0);
        searchB.value = '';
    }
}

export async function getRecentTopScorer() {
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
    
    await tsTable(data, 'top_players',
        `Top Scorers from ${data.recent_games[0].game_date}`, 
        0, 'recent');

    await resp.getP(base, top_scorer, 88888, 0, data);
}

// build top x players table
export async function getLeagueTop5() {
    const r = await fetch(`${base}/players/top5`);
    if (!r.ok) {
        throw new Error(`${r.status}: error calling /players/top5`)
    }
    const data = await r.json();
    console.log(data);
    await tsLgTable(data, 'top_lg_players', 'Top 5 NBA/WNBA Scorers')
}

export async function tsLgTable(data, elName, caption) {
    const tblcont = document.getElementById(elName);

    const tbl = document.getElementById('nba_tstbl');
    const wnbatbl = document.getElementById('wnba_tstbl');

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
    
    for (let i = 0; i < 10; i++) {
        let r = document.createElement('tr');

        let pName = document.createElement('td');
        let pts = document.createElement('td');
        let wpName = document.createElement('td');
        let wpts = document.createElement('td');

        let btn = document.createElement('button');
        btn.textContent = `${data.nba[i].player} | ${data.nba[i].team}`;
        btn.type = 'button';
        btn.addEventListener('click', async () => {
            await playerLinkSearch(data.nba[i].player);
        }); 
        console.log(data.nba[i].player);
        console.log(data);
        let wbtn = document.createElement('button');
        wbtn.textContent = `${data.wnba[i].player} | ${data.wnba[i].team}`;
        wbtn.type = 'button';
        console.log(data.wnba[i].player);
        wbtn.addEventListener('click', async () => {
            await playerLinkSearch(data.wnba[i].player);
        }); 

        pName.appendChild(btn);
        wpName.appendChild(wbtn);

        pts.textContent = data.nba[i].points;
        wpts.textContent = data.wnba[i].points;
        r.appendChild(pName);
        r.appendChild(pts);
        r.appendChild(wpName);
        r.appendChild(wpts);

        tbl.appendChild(r);
    }

    tblcont.appendChild(tbl);
}


export async function tsTable(data, elName, caption) {
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
        await topPlayerRow(tbl, scorers[i], data, 'recent');
    }

    tblcont.appendChild(tbl);
}

// call in top scorer table loop for each player
export async function topPlayerRow(tbl, scorer, data, mode) {
    let game = data.recent_games.find(g => g.player_id === scorer.player_id);
    let r = document.createElement('tr');

    let pName = document.createElement('td');
    let pTeam = document.createElement('td');
    let pts = document.createElement('td');

    let btn = document.createElement('button');
    btn.textContent = scorer.player;
    btn.type = 'button';
    btn.addEventListener('click', async () => {
        await playerLinkSearch(scorer.player);
    }); 

    pName.appendChild(btn);
    if (mode == 'recent') {
        pTeam.textContent = `${game ? game.team_name : ""} | \
        ${game ? game.matchup : ""} | ${game ? game.wl : ""}`;
    } else if (mode == 'league') {
        pTeam.textContent = scorer.team;
    }
    pts.textContent = scorer.points;
    r.appendChild(pName);
    r.appendChild(pTeam);
    r.appendChild(pts);

    tbl.appendChild(r);
}
export async function recGamesTbl(tbl) {
    // table headers
    const thead = document.createElement('thead');
    const nameH = document.createElement('td');
    const ptsH = document.createElement('td');
    const teamH = document.createElement('td');

    nameH.textContent = 'name';
    ptsH.textContent = 'points';
    teamH.textContent = 'team | matchup | win-loss';

    thead.appendChild(nameH);
    thead.appendChild(teamH);
    thead.appendChild(ptsH);

    tbl.appendChild(thead);
}
