import * as table from "./table.js"

export async function getP(base, player) { // add season & team
    const err = document.getElementById('sErr');
    if (err.style.display === "block") {
        err.style.display = 'none';
    }

    const p = encodeURIComponent(player).toLowerCase();

    const r = await fetch(base + `/player?player=${p}&season=88888`);
    if (!r.ok) {
        throw new Error(`HTTP Error: ${r.status}`);
    }

    const js = await r.json();
    const data = js.player[0];

    if (data.player_meta.player_id === 0) {
        err.textContent = `'${player}' not found...`
        err.style.display = "block"
        throw new Error(`Player not found error`);
    } 
    await buildPDash(data, 'main');
}

export async function buildPDash(data, pElName) {
    const pEl = document.getElementById(pElName);
    pEl.textContent = '';

    await appendImg(data.player_meta.headshot_url, 'imgs', 'pl_img');
    await appendImg(data.player_meta.team_logo_url, 'imgs', 'tm_img');
    await playerResTitle(data.player_meta, 'player_title');

    // box stat tables
    await table.basicTable(data.totals.box_stats, 'Total Box Stats', 'box');
    await table.basicTable(data.per_game.box_stats, 'Avg Box Stats', 'avg-box');

    // shooting stats tables
    await table.rowHdrTable(data.totals.shooting, 'Total Shooting Stats', 
        'shot_type', 'shooting');
    await table.rowHdrTable(data.per_game.shooting, 'Per Game Shooting Stats', 
        'shot_type', 'avg-shooting');
}

async function appendImg(url, pElName, cElName) {
    const pEl = document.getElementById(pElName);
    const cEl = document.getElementById(cElName);
    const img = document.createElement('img');
    cEl.textContent = ''; // clear child element
    img.src = url;
    img.alt = "image not found";
    cEl.appendChild(img);
    pEl.append(cEl);
}

async function playerResTitle(data, elName) {
    const cont = document.getElementById(elName);
    cont.textContent = '';
    const d = document.createElement('div');
    const t = document.createElement('h1');
    const s = document.createElement('h2');
    t.textContent = data.caption;
    s.textContent = data.season;
    d.append(t);
    d.append(s);
    cont.append(d);
}