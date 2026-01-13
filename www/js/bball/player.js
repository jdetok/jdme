import { base, bytes_in_resp, scrollIntoBySize, MSG, foldedLog, MSG_BOLD, RED_BOLD } from "../global.js";
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
// get the top scorer from each game from the most recent night where games occured
// (usually dated yesterday, but when no games occur it'll get the most recent day
// where games did occur). called on page load, it creates a table with all these
// scorers and immediately grabs and loads the player dash for the top overall 
// scorer. use season id 88888 in getP to get most recent season
export async function getRecentGamesData() {
    const url = `${base}/games/recent`;
    const r = await fetch(url);
    if (!r.ok) {
        console.error(`%cerror fetching ${url}`, RED_BOLD);
    }
    foldedLog(`%c ${await bytes_in_resp(r)} bytes received from ${url}}`, MSG_BOLD);
    return await r.json();
}
export async function buildOnLoadDash() {
    await searchPlayer('onload');
}
export async function searchPlayer(pst = 'submit', playerOverride) {
    const searchElId = 'pSearch';
    const input = document.getElementById(searchElId);
    if (!input)
        throw new Error(`couldn't get element at Id ${searchElId}`);
    let player = '';
    const { lg, season, team } = await getInputVals();
    let recent_data;
    switch (pst) {
        case 'onload':
            recent_data = await getRecentGamesData();
            player = recent_data.top_scorers[0].player_id;
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
            if (!playerOverride) {
                console.error(`%cplayerOverride must be passed if called with pst='button'`, RED_BOLD);
                return;
            }
            else {
                player = playerOverride;
            }
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
    // build response player dash section
    let js = await fetchPlayer(base, player, season, team, lg);
    if (js) {
        await setPHold(js.player[0].player_meta.player);
        await buildPlayerDash(js.player[0], recent_data);
        if (pst !== 'onload')
            scrollIntoBySize(1350, 1250, "player_title");
    }
}
export async function fetchPlayer(base, player, season, team, lg) {
    const errEl = 'sErr';
    const errmsg = document.getElementById(errEl); // sErr is elId
    if (!errmsg)
        throw new Error(`%ccould not find error string element at ${errEl}`);
    if (errmsg.style.display === "block") {
        errmsg.style.display = 'none';
    }
    // encode passed args to be ready for query string
    const s = encodeURIComponent(season);
    const p = encodeURIComponent(player).toLowerCase();
    const req = `${base}/v2/players?player=${p}&season=${s}&team=${team}&league=${lg}`;
    // attempt to fetch from /player endpoint with encoded params
    try {
        const r = await fetch(req);
        if (!r.ok) {
            throw new Error(`HTTP Error (${r.status}) attempting to fetch ${player}`);
        }
        foldedLog(`%c ${await bytes_in_resp(r)} bytes received from ${req}}`, MSG);
        const js = await r.json();
        if (js) {
            if (js.error_string) {
                errmsg.textContent = js.error_string;
                errmsg.style.display = "block";
                return;
            }
            else {
                return js;
            }
        }
    }
    catch (err) {
        errmsg.textContent = `can't find ${player}`;
        errmsg.style.display = "block";
        console.error(`an error occured attempting to fetch ${player}\n${err}`);
    }
}
// accept player dash data, build tables/fetch images and display on screen
export async function buildPlayerDash(data, ts, el = PLAYER_DASH_ELS) {
    foldedLog(`%cts: ${ts ? `fetching top scorer from ${ts.recent_games[0].game_date}` : 'no ts var, normal fetch'}`, MSG);
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
    d.append(s);
    d.append(u);
    cont.append(d);
}
async function appendImg(url, elId) {
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