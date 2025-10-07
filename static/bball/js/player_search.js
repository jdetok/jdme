import { base } from "./listen.js"
import { checkBoxGroupValue, lgRadioBtns } from "./ui.js";
import { getPlayerStats, buildPlayerDash } from "./player_dash.js";

/*
get the top scorer from each game from the most recent night where games occured
(usually dated yesterday, but when no games occur it'll get the most recent day
where games did occur). called on page load, it creates a table with all these
scorers and immediately grabs and loads the player dash for the top overall 
scorer. use season id 88888 in getP to get most recent season
*/
export async function getRecentGamesData() {
    const r = await fetch(`${base}/games/recent`);
    if (!r.ok) {
        throw new Error(`${r.status}: error calling /games/recent`);
    }
    const data = await r.json();
    if (data) {
        return data
    }
}

// called in listen.js
export async function buildLoadDash(recent_game_data) {
    const top_scorer = recent_game_data.top_scorers[0].player_id;
    const lg = await lgRadioBtns();
    let js = await getPlayerStats(base, top_scorer, 88888, 0, lg);
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
        console.log(`league in searchPlayer: ${lg}`);
        
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

        // build response player dash section
        let js = await getPlayerStats(base, player, season, team, lg);
        if (js) {
            await setPHold(js.player[0].player_meta.player);
            await buildPlayerDash(js.player[0], 0);
        }

        // TODO: maybe fill error string visible to user at this point?
        // clear player search box
        // input.value = ''; // clear input box after searching
    }) 
}

export async function setPHold(player) {
    document.getElementById('pHold').value = player;
}

// get a random player from the API and getPlayerStats
export async function randPlayerBtn() {
    // listen for random player button press
    const btn = document.getElementById('randP');
    btn.addEventListener('click', async (event) => {        
        event.preventDefault();

        const lg = await lgRadioBtns();
        console.log(`league in randPlayerBtn: ${lg}`);
        // check season boxes & get appropriate season id, search with random as player
        const season = await checkBoxGroupValue(
            {box: 'post', slct: 'ps_slct'}, 
            {box: 'reg', slct: 'rs_slct'}, 
            88888);
        console.log(`searching random player for season ${season}`);
        let js = await getPlayerStats(base, 'random', season, 0, lg);
        if (js) {
            await buildPlayerDash(js.player[0], 0);
        }
    });
}

/* 
adds a button listener to each individual player button in the leading scorers
tables. have to create a button, do btn.AddEventListener, and call this function
within that listener. will insert the player's name in the search bar and call 
getP
*/
export async function playerBtnListener(player) {
    let searchB = document.getElementById('pSearch');
    if (searchB) {
        // searchB.value = player;
        // searchB.value = '';
        const season = await checkBoxGroupValue(
            {box: 'post', slct: 'ps_slct'}, 
            {box: 'reg', slct: 'rs_slct'}, 
            88888);

        const team = await checkBoxGroupValue(
            {box: 'nbaTm', slct: 'tm_slct'}, 
            {box: 'wnbaTm', slct: 'wTm_slct'}, 0);

        const lg = await lgRadioBtns();

        console.log(`player: ${player} | tm: ${team} | szn: ${season}`);

        // search & clear player search bar
        let js = await getPlayerStats(base, player, season, team, lg);
        if (js) {
            console.log(js);
            await setPHold(js.player[0].player_meta.player);
            await buildPlayerDash(js.player[0], 0);
        }
        // if screen is small scroll into it
        if (window.innerWidth <= 700) {
            let res = document.getElementById("ui");
            if (res) {
                res.scrollIntoView({behavior: "smooth", block: "start"});
            }
        }
    }
}

// read pHold invisible val to add on-screen player's name to search bar
export async function holdPlayerBtn() {
    // listen for hold player button press
    const btn = document.getElementById('holdP');
    btn.addEventListener('click', async (event) => {
        event.preventDefault();

        // get player name held in pHold value, fill player search bar with it
        let player = document.getElementById('pHold').value;
        document.getElementById('pSearch').value = player;
    })
}

// clear search box
export async function clearSearch() {
    const btn = document.getElementById('clearS');
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        let pSearch = document.getElementById('pSearch');
        pSearch.value = '';
        pSearch.focus();
    })
}