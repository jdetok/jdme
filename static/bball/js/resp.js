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
    
    await tsTable(data, 'top_players', 0);

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

// replacement to tsTable, should work for both recent game and league top scorer
export async function buildTopScorersTable(data, outerEl, tblEl, caption, exclude_first) {
    const tblcont = document.getElementById(outerEl);
    const tbl = document.getElementById(tblEl);
    table.tblCaption(tbl, caption);

}


export async function tsTable(data, elName, exclude_first) {
    const tblcont = document.getElementById(elName);
    const tbl = document.getElementById('tstbl');

    // table captions
    const caption = `Top Scorers from ${data.recent_games[0].game_date}`;
    table.tblCaption(tbl, caption);

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
    
    // 
    let scorers = data.top_scorers;
    if (exclude_first) {
        scorers = scorers.slice(1);
    }
    
    for (let i = 0; i < scorers.length; i++) {
        await table.topPlayerRow(tbl, scorers[i], data, 'recent');
    }

    tblcont.appendChild(tbl);
}

export async function playerLinkSearch(player, data) {
    let searchB = document.getElementById('pSearch');
    if (searchB) {
        searchB.value = player;
        searchB.focus();
        await getP(base, player, 88888, 0, data);
    }
}