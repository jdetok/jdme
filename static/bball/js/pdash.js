import * as table from "./table.js"

export async function getP(base, player, season, team) { // add season & team
    const err = document.getElementById('sErr');
    const loading = document.getElementById('loading');
    if (err.style.display === "block") {
        err.style.display = 'none';
    }

    const s = encodeURIComponent(season)
    const p = encodeURIComponent(player).toLowerCase();
    loading.textContent = `loading player dash for ${player}...`;
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
    await buildPDash(data, 'player');
    document.getElementById('pHold').value = data.player_meta.player;
}

export async function buildPDash(data, pElName) {
    const loading = document.getElementById('loading');
    loading.textContent = '';
    await appendImg(data.player_meta.headshot_url, 'pl_img');
    await appendImg(data.player_meta.team_logo_url, 'tm_img');
    await playerResTitle(data.player_meta, 'player_title');
    await info(data, 'player_szn');

    // box stat tables
    await table.basicTable(data.totals.box_stats, data.player_meta.cap_box_tot, 'box');
    // await table.basicTable(data.totals.box_stats, 'Total Box Stats', 'box');
    await table.basicTable(data.per_game.box_stats, data.player_meta.cap_box_avg, 'avg-box');

    // shooting stats tables
    await table.rowHdrTable(data.totals.shooting, data.player_meta.cap_shtg_tot, 
        'shot_type', 'shooting');
    await table.rowHdrTable(data.per_game.shooting, data.player_meta.cap_shtg_avg, 
        'shot_type', 'avg-shooting');
}

async function appendImg(url, pElName) {
    console.log("in image thing")
    console.log(pElName);
    const pEl = document.getElementById(pElName);
    const img = document.createElement('img');
    pEl.textContent = ''; // clear child element
    img.src = url;
    img.alt = "image not found";
    pEl.append(img);
}
async function appendImgOld(url, pElName, cElName) {
    const pEl = document.getElementById(pElName);
    const cEl = document.getElementById(cElName);
    const img = document.createElement('img');
    cEl.textContent = ''; // clear child element
    img.src = url;
    img.alt = "image not found";
    cEl.appendChild(img);
    pEl.append(cEl);
}

async function info(data, elName) {
    const cont = document.getElementById(elName);
    cont.textContent = '';
    const d = document.createElement('div');
    // const t = document.createElement('h1');
    const s = document.createElement('h2');
    const u = document.createElement('h3');
    s.textContent = data.player_meta.season;
    u.textContent = `${data.playtime.games_played} Games Played | 
        ${data.playtime.minutes} Minutes Played | 
        ${data.playtime.minutes_pg} Minutes Per Game`;
    d.append(s);
    d.append(u);
    cont.append(d);
}

async function playerResTitle(data, elName) {
    const cont = document.getElementById(elName);
    cont.textContent = '';
    const d = document.createElement('div');
    const t = document.createElement('h1');
    const s = document.createElement('h2');
    t.textContent = data.caption;
    d.append(t);
    cont.append(d);
}