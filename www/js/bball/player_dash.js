import * as table from "./table.js"
import { bytes_in_resp, foldedLog, MSG } from "./util.js";

// html elements to fill
const PLAYER_DASH_ELS = {
    title: 'player_title', 
    season: 'player_szn',
    img: {
        player: 'pl_img',
        team: 'tm_img',
    },
    tables: {
        total_boxstats: 'box',
        avg_boxstats: 'avg-box',
        shooting: 'shooting',
        avg_shooting: 'avg-shooting',
    },
};

// accept player dash data, build tables/fetch images and display on screen
export async function buildPlayerDash(data, ts, el = PLAYER_DASH_ELS) {
    // console.trace(data);
    // console.trace(`%cts: ${ts ? ts : 'ts var empty'}`, FUSC);
    foldedLog(`%cts: ${ts ? `fetching top scorer from ${ts.recent_games[0].game_date}` : 'no ts var, normal fetch'}`, MSG);
    await appendImg(data.player_meta.headshot_url, el.img.player);
    await appendImg(data.player_meta.team_logo_url, el.img.team);

    await respPlayerTitle(data.player_meta, el.title, ts);
    await respPlayerInfo(data, el.season);

    // box stat tables
    await table.basicTable(data.totals.box_stats, data.player_meta.cap_box_tot, el.tables.total_boxstats);
    await table.basicTable(data.per_game.box_stats, data.player_meta.cap_box_avg, el.tables.avg_boxstats);

    // shooting stats tables
    await table.rowHdrTable(data.totals.shooting, data.player_meta.cap_shtg_tot, 'shot type', el.tables.shooting);
    await table.rowHdrTable(data.per_game.shooting, data.player_meta.cap_shtg_avg, 'shot type', el.tables.avg_shooting);
}

// ts indicates 'top scorer' - used when called on page refresh to get recent game
export async function getPlayerStatsV2(base, player, season, team, lg) { // add season & team
    const errmsg = document.getElementById('sErr');
    if (errmsg.style.display === "block") {
        errmsg.style.display = 'none';
    }

    // encode passed args to be ready for query string
    const s = encodeURIComponent(season)
    const p = encodeURIComponent(player).toLowerCase();

    const req = `${base}/v2/players?player=${p}&season=${s}&team=${team}&league=${lg}`;
    // attempt to fetch from /player endpoint with encoded params
    try {
        const r = await fetch(req);
        if (!r.ok) {
            throw new Error(`HTTP Error (${r.status}) attempting to fetch ${player}`);
        }
        foldedLog(`%c ${await bytes_in_resp(r)} bytes received from ${req}}`, MSG)

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
        errmsg.textContent = `can't find ${player}`;
        errmsg.style.display = "block";
        console.error(`an error occured attempting to fetch ${player}\n${err}`);
        
    }
}

// ts is always nothing, except when buildPlayerDash is called on page load with recent games data
// in that case, ts exists and should be the object returned from /games/recent
async function respPlayerTitle(data, elName, ts) {
    const rTitle = document.getElementById(elName);
    if (ts) {
        rTitle.innerHTML = `
        Top Scorer from ${ts.recent_games[0].game_date}<br>${data.caption}
         | ${ts.top_scorers[0].points} pts | 
         ${ts.top_scorers[0].assists} ast |
         ${ts.top_scorers[0].rebounds} reb`;    
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
    pEl.textContent = '';
    img.src = url;
    img.alt = "image not found";
    pEl.append(img);
}