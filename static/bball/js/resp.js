import * as table from "./table.js"
import { base } from "./listen.js"
import { search } from "./ui.js";

export async function buildPDash(data, ts) {
    await appendImg(data.player_meta.headshot_url, 'pl_img');
    await appendImg(data.player_meta.team_logo_url, 'tm_img');
    await respPlayerTitle(data.player_meta, 'player_title', ts);
    await respPlayerInfo(data, 'player_szn');

    

    // box stat tables
    await table.basicTable(data.totals.box_stats, data.player_meta.cap_box_tot, 'box');
    await table.basicTable(data.per_game.box_stats, data.player_meta.cap_box_avg, 'avg-box');

    // shooting stats tables
    await table.rowHdrTable(data.totals.shooting, data.player_meta.cap_shtg_tot, 
        'shot_type', 'shooting');
    await table.rowHdrTable(data.per_game.shooting, data.player_meta.cap_shtg_avg, 
        'shot_type', 'avg-shooting');
}

// ts indicates 'top scorer' - used when called on page refresh to get recent game
export async function getP(base, player, season, team, ts) { // add season & team
    const err = document.getElementById('sErr');
    if (err.style.display === "block") {
        err.style.display = 'none';
    }

    const s = encodeURIComponent(season)
    const p = encodeURIComponent(player).toLowerCase();
    const r = await fetch(base + `/player?player=${p}&season=${s}&team=${team}`);
    if (!r.ok) {
        throw new Error(`HTTP Error: ${r.status}`);
    }
    
    const js = await r.json();
    const data = js.player[0];

    if (data.player_meta.player_id === 0) {
        if (player != '') {
            err.textContent = `'${player}' not found...`
            err.style.display = "block"    
        }
        throw new Error(`Player not found error`);
    } 
    await buildPDash(data, ts);
    document.getElementById('pHold').value = data.player_meta.player;
}

// RESULT TITLE - LIKE `LeBron James - Los Angeles Lakers`
async function respPlayerTitle(data, elName, ts) {
    const rTitle = document.getElementById(elName);
    if (ts) {
        rTitle.innerHTML = `
        Top Scorer from ${ts.recent_games[0].game_date}<br>${data.caption}
         | ${ts.top_scorers[0].points} points`;    
    } else {
        rTitle.textContent = data.caption;
    }
}

async function respPlayerInfo(data, elName) {
    const cont = document.getElementById(elName);
    cont.textContent = '';
    const d = document.createElement('div');
    const s = document.createElement('h2');
    const u = document.createElement('h3');
    s.textContent = data.player_meta.season;
    u.textContent = `${data.playtime.games_played} Games Played | 
        ${data.playtime.minutes} Minutes | 
        ${data.playtime.minutes_pg} Minutes/Game`;
    d.append(s);
    d.append(u);
    cont.append(d);
}

async function appendImg(url, pElName) {
    const pEl = document.getElementById(pElName);
    const img = document.createElement('img');
    pEl.textContent = ''; // clear child element
    img.src = url;
    img.alt = "image not found";
    pEl.append(img);
}


export async function getRecentTopScorer() {
    const r = await fetch(`${base}/games/recent`);
    if (!r.ok) {
        throw new Error(`${r.status}: error calling /games/recent`);
    }
    const data = await r.json();
    const top_scorer = data.top_scorers[0].player_id;
    /* TODO: BUILD A TABLE OF THE OTHER TOP PLAYERS, LINK THE PLAYERS' NAMES TO
    CALL getP FOR THAT PLAYER
    */

    const plrs = data.top_scorers.slice(1);
    console.log(plrs);
    console.log("FULL DATA: ");
    console.log(data);
    
    await tsTable(data, 'top_players',
        `Top Scorers from ${data.recent_games[0].game_date}`, 
        0, 'recent');

    await getP(base, top_scorer, 88888, 0, data);
    await getLeagueTop5();
}


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
    
    for (let i = 0; i < 5; i++) {
        let r = document.createElement('tr');

        let pName = document.createElement('td');
        let pts = document.createElement('td');
        let wpName = document.createElement('td');
        let wpts = document.createElement('td');

        let btn = document.createElement('button');
        btn.textContent = `${data.nba[i].player} | ${data.nba[i].team}`;
        btn.type = 'button';
        btn.addEventListener('click', async () => {
            await playerLinkSearch(data.nba[i].player, data);
        }); 

        let wbtn = document.createElement('button');
        wbtn.textContent = `${data.wnba[i].player} | ${data.wnba[i].team}`;
        wbtn.type = 'button';
        wbtn.addEventListener('click', async () => {
            await playerLinkSearch(data.wnba[i].player, data);
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

export async function playerLinkSearch(player, data) {
    let searchB = document.getElementById('pSearch');
    if (searchB) {
        searchB.value = player;
        searchB.focus();
        await getP(base, player, 88888, 0, data);
        searchB.value = '';

    }
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
        await playerLinkSearch(scorer.player, data);
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

