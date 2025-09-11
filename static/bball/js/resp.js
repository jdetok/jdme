import * as table from "./table.js"
import { base } from "./listen.js"

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

export async function otherTopPlayersTable(data, elName) {
    const tblcont = document.getElementById(elName);
    const tbl = document.getElementById('tstbl');

    const lbl = document.createElement('caption');
    lbl.textContent = `Other Top Scorers from ${data.recent_games[0].game_date}`;
    tbl.appendChild(lbl);

    const thead = document.createElement('thead');
    
    const nameH = document.createElement('td');
    const lgH = document.createElement('td');
    const ptsH = document.createElement('td');
    const gmH = document.createElement('td');
    const wlH = document.createElement('td');
    const teamH = document.createElement('td');

    nameH.textContent = 'name';
    lgH.textContent = 'league';
    ptsH.textContent = 'points';
    gmH.textContent = 'game';
    wlH.textContent = 'win/loss';
    teamH.textContent = 'team';

    thead.appendChild(nameH);
    thead.appendChild(lgH);
    thead.appendChild(teamH);
    thead.appendChild(gmH);
    thead.appendChild(wlH);
    thead.appendChild(ptsH);

    tbl.appendChild(thead);
    
    const scorers = data.top_scorers.slice(1);

    for (let i = 0; i < scorers.length; i++) {
        let scorer = scorers[i];
        let game = data.recent_games.find(g => g.player_id === scorer.player_id);

        let r = document.createElement('tr');

        let pName = document.createElement('td');
        let pTeam = document.createElement('td');
        let pLg = document.createElement('td');
        let pGm = document.createElement('td');
        let pWl = document.createElement('td');
        let pts = document.createElement('td');

        pName.textContent = scorer.player;
        pTeam.textContent = game ? game.team_name : "";
        pLg.textContent = scorer.league;
        pGm.textContent = game ? game.matchup : "";
        pWl.textContent = game ? game.wl : "";
        pts.textContent = scorer.points;

        r.appendChild(pName);
        r.appendChild(pTeam);
        r.appendChild(pLg);
        r.appendChild(pGm);
        r.appendChild(pWl);
        r.appendChild(pts);

        tbl.appendChild(r);
    }

    tblcont.appendChild(tbl);
    
    // await table.basicTable(data, "Other Top Scorers", elName)
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
    
    await otherTopPlayersTable(data, 'top_players');

    await getP(base, top_scorer, 88888, 0, data);

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