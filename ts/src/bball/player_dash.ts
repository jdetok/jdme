import { MSG, foldedLog, foldedErr, scrollIntoBySize, SBL } from "../global.js";
import { tblColHdrs, tblRowColHdrs } from "./tbls_resp.js";
import { makeTmPlrImageDiv, fillImageDiv } from "./elements.js";
import { TopScorer, PlayerResp, playerMeta, RGData, PlayersResp } from "./resp_types.js";
import { searchPlayer, PlayerSearchType } from "./player_search.js";
import { setPHold } from "./inputs.js";

export async function fetchAndBuildPlayerDash(pst: PlayerSearchType = 'submit',
    playerOverride: string | null = null, rgData?: RGData,
): Promise<void> {
    foldedLog(`%fetching pst=${pst} player dash data`, SBL);
    let data: PlayersResp | null = null;
    try {
        data = await searchPlayer(pst, playerOverride, rgData);
    } catch (e) {
        foldedErr(`error searching player with player search type ${pst}`);
        return;
    }
    if (!data) throw new Error(`error getting data for search type ${pst}`);

    const playerResp = data.player[0];
    const player = playerResp.player_meta.player;
    foldedLog(`%cbuilding pst=${pst} dash for player ${player}`, SBL);
    try {
        await setPHold(player);
        await buildPlayerDash(playerResp, rgData ?? null);
    } catch (e) {
        foldedErr(`error building player with player search type ${pst}`);
        return;
    }
    if (pst !== 'onload') scrollIntoBySize(1350, 1250, 'resp_ttl');
}

// accept player dash data, build tables/fetch images and display on screen
async function buildPlayerDash(data: PlayerResp, ts: TopScorer): Promise<void> {
    try {
        await fillImageDiv(makeTmPlrImageDiv('dash_imgs', {
            tm_url: data.player_meta.team_logo_url,
            tm: data.player_meta.team,
            plr_url: data.player_meta.headshot_url,
            plr: data.player_meta.player,
        }));

        await respPlayerTitle(data.player_meta, 'resp_ttl', ts);
        await respPlayerInfo(data, 'resp_subttl');

        // box stat tables
        await tblColHdrs(data.totals.box_stats, data.player_meta.cap_box_tot, 'box_tot');
        await tblColHdrs(data.per_game.box_stats, data.player_meta.cap_box_avg, 'box_avg');

        // shooting stats tables
        await tblRowColHdrs(data.totals.shooting, data.player_meta.cap_shtg_tot, 'shot type', 'stg_tot');
        await tblRowColHdrs(data.per_game.shooting, data.player_meta.cap_shtg_avg, 'shot type', 'stg_avg');
    } catch (e) {
        foldedErr(`error building ${ts ? 'top scorer' : ''} player dash for ${data.player_meta.player}: ${e}`);
        return;
    }
    foldedLog(`%cbuilt ${ts ? 'top scorer' : ''} player dash for ${data.player_meta.player}`, MSG);
}

// ts is nothing 
// in that case, ts exists and should be the object returned from /games/recent
async function respPlayerTitle(data: playerMeta, elId: string, ts: TopScorer) {
    const rTitle = document.getElementById(elId);

    if (!rTitle) throw new Error(`couldnt' get response title element at ${elId}`);

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

async function respPlayerInfo(data: PlayerResp, elId: string) {
    const cont = document.getElementById(elId);
    if (!cont) throw new Error(`couldnt' get response title element at ${elId}`);
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