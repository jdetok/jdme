import { base, scrollIntoBySize, MSG, foldedLog, MSG_BOLD, RED_BOLD, foldedErr, logResp } from "../global.js";
import { setPHold, getInputVals } from "./inputs.js";
import { tblColHdrs, tblRowColHdrs } from "./tbls_resp.js";
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
export async function buildOnLoadDash(rgData) {
    try {
        await searchPlayer('onload', null, rgData);
    }
    catch (e) {
        throw new Error(`error building onload dash: ${e}`);
    }
}
export async function searchPlayer(pst = 'submit', playerOverride = null, rgData) {
    const searchElId = 'pSearch';
    const input = document.getElementById(searchElId);
    if (!input)
        throw new Error(`couldn't get element at Id ${searchElId}`);
    let player = '';
    const { lg, season, team } = await getInputVals();
    switch (pst) {
        case 'onload':
            if (!rgData)
                throw new Error(`%cmust pass recent game data for onload mode`);
            player = rgData.top_scorers[0].player;
            foldedLog(`%cfetched onload player ${player} | league: ${lg} | season ${season}`, MSG);
            break;
        case 'submit':
            player = input.value.trim();
            break;
        case 'random':
            foldedLog(`%csearching random player | league: ${lg} | season ${season}`, MSG);
            player = String(pst);
            break;
        case 'button':
            foldedLog(`%cfetching player ${player} from button | league: ${lg} | season ${season}`, MSG);
            if (!playerOverride)
                throw new Error(`cplayerOverride must be passed if called with pst='button'`);
            player = playerOverride;
            break;
        default:
            foldedLog(`%coption passed as pst "${pst} not valid`, RED_BOLD);
            return;
    }
    if (player === '') {
        const pHoldElId = 'pHold';
        const pHoldEl = document.getElementById(pHoldElId);
        player = pHoldEl.value;
    }
    foldedLog(`%csearching for player ${player} | season ${season} | team ${team} | league ${lg}`, MSG_BOLD);
    const recent_data = rgData ?? null;
    // build response player dash section
    try {
        const data = await fetchPlayer(base, player, season, team, lg);
        if (!data)
            throw new Error(`failed to get data for player ${player} | season ${season} | team ${team}`);
        await setPHold(data.player[0].player_meta.player);
        await buildPlayerDash(data.player[0], recent_data);
    }
    catch (e) {
        throw new Error(`error fetching data or building dash: ${e}`);
    }
    if (pst !== 'onload')
        scrollIntoBySize(1350, 1250, "player_title");
}
export async function fetchPlayer(base, player, season, team, lg, errEl = 'sErr') {
    const errmsg = document.getElementById(errEl);
    if (!errmsg)
        throw new Error(`%ccould not find error string element at ${errEl}`);
    if (errmsg.style.display === "block") {
        errmsg.style.display = 'none';
    }
    const s = encodeURIComponent(season);
    const p = encodeURIComponent(player).toLowerCase();
    const url = `${base}/v2/players?player=${p}&season=${s}&team=${team}&league=${lg}`;
    let r;
    try {
        r = await fetch(url);
        if (!r.ok) {
            errmsg.textContent = `can't find ${player}`;
            errmsg.style.display = "block";
            console.error(`an error occured attempting to fetch ${player}\n`);
            throw new Error(`HTTP Error (${r.status}) attempting to fetch ${player}`);
        }
    }
    catch (e) {
        throw new Error(`fetch player error: ${e}`);
    }
    // foldedLog(`%c${await bytes_in_resp(r)} bytes received from ${url}}`, MSG)
    await logResp(url, r);
    return await r.json();
}
// accept player dash data, build tables/fetch images and display on screen
export async function buildPlayerDash(data, ts, el = PLAYER_DASH_ELS) {
    try {
        await appendImg(data.player_meta.headshot_url, el.img.player);
        await appendImg(data.player_meta.team_logo_url, el.img.team);
        await respPlayerTitle(data.player_meta, el.title, ts);
        await respPlayerInfo(data, el.season);
        // box stat tables
        await tblColHdrs(data.totals.box_stats, data.player_meta.cap_box_tot, el.tables.total_boxstats);
        await tblColHdrs(data.per_game.box_stats, data.player_meta.cap_box_avg, el.tables.avg_boxstats);
        // shooting stats tables
        await tblRowColHdrs(data.totals.shooting, data.player_meta.cap_shtg_tot, 'shot type', el.tables.shooting);
        await tblRowColHdrs(data.per_game.shooting, data.player_meta.cap_shtg_avg, 'shot type', el.tables.avg_shooting);
    }
    catch (e) {
        foldedErr(`error building ${ts ? 'top scorer' : ''} player dash for ${data.player_meta.player}: ${e}`);
        return;
    }
    foldedLog(`%cbuilt ${ts ? 'top scorer' : ''} player dash for ${data.player_meta.player}`, MSG);
}
// ts is always nothing, except when buildPlayerDash is called on page load with recent games data
// in that case, ts exists and should be the object returned from /games/recent
async function respPlayerTitle(data, elId, ts) {
    const rTitle = document.getElementById(elId);
    if (!rTitle)
        throw new Error(`couldnt' get response title element at ${elId}`);
    if (ts) {
        rTitle.innerHTML = `
        Top Scorer from ${ts.recent_games[0].game_date}<br>${data.caption}
         | ${ts.top_scorers[0].points} pts | 
         ${ts.top_scorers[0].assists} ast |
         ${ts.top_scorers[0].rebounds} reb`;
    }
    else {
        rTitle.textContent = data.caption;
    }
}
async function respPlayerInfo(data, elId) {
    const cont = document.getElementById(elId);
    if (!cont)
        throw new Error(`couldnt' get response title element at ${elId}`);
    cont.textContent = '';
    const d = document.createElement('div');
    const s = document.createElement('h2');
    const u = document.createElement('h3');
    s.textContent = data.player_meta.season;
    u.textContent = `${data.playtime.games_played} Games Played | 
        ${data.playtime.minutes} Minutes | 
        ${data.playtime.minutes_pg} Minutes/Game`;
    d.appendChild(s);
    d.appendChild(u);
    cont.appendChild(d);
}
export async function appendImg(url, elId) {
    const pEl = document.getElementById(elId);
    if (!pEl)
        throw new Error(`couldnt' get response title element at ${elId}`);
    const img = document.createElement('img');
    pEl.textContent = '';
    img.src = url;
    img.alt = "image not found";
    pEl.append(img);
}
//# sourceMappingURL=player.js.map