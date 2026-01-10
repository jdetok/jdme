import { base, bytes_in_resp, foldedLog, MSG_BOLD, YLW_BOLD } from "./util.js"
import { checkBoxGroupValue, lgRadioBtns, setPHold } from "./ui.js";
import { buildPlayerDash, getPlayerStatsV2 } from "./player_dash.js";

// get the top scorer from each game from the most recent night where games occured
// (usually dated yesterday, but when no games occur it'll get the most recent day
// where games did occur). called on page load, it creates a table with all these
// scorers and immediately grabs and loads the player dash for the top overall 
// scorer. use season id 88888 in getP to get most recent season
export async function getRecentGamesData() {
    const url = `${base}/games/recent`;
    const r = await fetch(url);
    if (!r.ok) {
        console.error(`%cerror fetching ${url}`, YLW_BOLD);
    }
    await foldedLog(`%c ${await bytes_in_resp(r)} bytes received from ${url}}`, MSG_BOLD);

    const data = await r.json();
    if (data) {
        return data
    }
}

// called in listen.js
export async function buildLoadDash(recent_game_data) {
    const top_scorer = recent_game_data.top_scorers[0].player_id;
    const lg = await lgRadioBtns();
    let js = await getPlayerStatsV2(base, top_scorer, 88888, 0, lg);
    await buildPlayerDash(js.player[0], recent_game_data);
}

// get player from search bar and make player dash
export async function searchPlayer() {
    // listen for form submission
    const frm = document.getElementById('ui');
    frm.addEventListener('submit', async (event) => {
        event.preventDefault();

        // get value of player search box
        const input = document.getElementById('pSearch');
        let player = input.value.trim();

        const lg = await lgRadioBtns();
        
        // if search pressed without anything in search box, searches current player
        if (player === '') {
            player = document.getElementById('pHold').value;
        }

        // check if season box is checked, return sel val if so, 88888 if not
        // 88888 gets the most recent season from the api
        // const season = await handleSeasonBoxes();
        const season = await checkBoxGroupValue(
            {box: 'post', slct: 'ps_slct'}, 
            {box: 'reg', slct: 'rs_slct'}, 
            88888);

        const team = await checkBoxGroupValue(
            {box: 'nbaTm', slct: 'tm_slct'}, 
            {box: 'wnbaTm', slct: 'wTm_slct'}, 
            0);

        await foldedLog(`%csearching for player ${player} | season ${season} | team ${team} | league ${lg}`, MSG_BOLD);
        // console.trace();
        // build response player dash section
        let js = await getPlayerStatsV2(base, player, season, team, lg);
        if (js) {
            await setPHold(js.player[0].player_meta.player);
            await buildPlayerDash(js.player[0], 0);
        }
    }) 
}


