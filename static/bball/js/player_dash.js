import * as table from "./table.js"

export async function buildPlayerDash(data, ts) {
    await appendImg(data.player_meta.headshot_url, 'pl_img');
    await appendImg(data.player_meta.team_logo_url, 'tm_img');
    await respPlayerTitle(data.player_meta, 'player_title', ts);
    await respPlayerInfo(data, 'player_szn');

    // box stat tables
    await table.basicTable(data.totals.box_stats, data.player_meta.cap_box_tot, 'box');
    await table.basicTable(data.per_game.box_stats, data.player_meta.cap_box_avg, 'avg-box');

    // shooting stats tables
    await table.rowHdrTable(data.totals.shooting, data.player_meta.cap_shtg_tot, 
        'shot type', 'shooting');
    await table.rowHdrTable(data.per_game.shooting, data.player_meta.cap_shtg_avg, 
        'shot type', 'avg-shooting');
}

// ts indicates 'top scorer' - used when called on page refresh to get recent game
export async function getPlayerStats(base, player, season, team, lg) { // add season & team
    const errmsg = document.getElementById('sErr');
    if (errmsg.style.display === "block") {
        errmsg.style.display = 'none';
    }

    // encode passed args to be ready for query string
    const s = encodeURIComponent(season)
    const p = encodeURIComponent(player).toLowerCase();

    const req = `${base}/player?player=${p}&season=${s}&team=${team}&league=${lg}`;
    // attempt to fetch from /player endpoint with encoded params
    try {
        const r = await fetch(req);
        if (!r.ok) {
            throw new Error(`HTTP Error: ${r.status}`);
        }
        
        // get json, set first object in player array as data var
        const js = await r.json()
        if (js) {
            if (js.error_string) {
                errmsg.textContent = js.error_string;
                errmsg.style.display = "block";
                return;
            } else {
                return js;
            }
        }
    } catch(err) {
        errmsg.textContent = `can't find ${player}` 
        console.log(`an error occured attempting to fetch ${player}\n${err}`);
        errmsg.style.display = "block";
        throw new Error(`HTTP Error: ${r.status}`);
    }

    // const data = js.player[0];

    
// handle empty player response

    // build and display player dash
    // await buildPlayerDash(data, ts);

    // set invisible player hold value as current player
    // document.getElementById('pHold').value = data.player_meta.player;
}



// set the player name/team name as result title
// RESULT TITLE - LIKE `LeBron James - Los Angeles Lakers`
async function respPlayerTitle(data, elName, ts) {
    const rTitle = document.getElementById(elName);
    // special title if caller specifies it's a top scorer
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