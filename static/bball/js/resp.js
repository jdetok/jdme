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

export async function getRecGames() {
    const r = await fetch(`${base}/games/recent`);
    if (!r.ok) {
        throw new Error(`${r.status}: error calling /games/recent`);
    }
    const data = await r.json();
    const player = data.top_scorers[0].player_id;
    await getP(base, player, 88888, 0, data);
}

async function respPlayerTitle(data, elName, ts) {
    const rTitle = document.getElementById(elName);
    if (ts) {
        rTitle.textContent = `${data.caption} - Top Scorer from 
            ${ts.recent_games[0].game_date} - ${ts.top_scorers[0].points} points`;    
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