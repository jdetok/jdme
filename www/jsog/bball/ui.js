import { base, checkBoxEls, foldedLog, MSG, scrollIntoBySize } from "./util.js";
import { getPlayerStatsV2, buildPlayerDash } from "./player_dash.js";
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
        // check season boxes & get appropriate season id, search with random as player
        const season = await checkBoxGroupValue(
            {box: 'post', slct: 'ps_slct'}, 
            {box: 'reg', slct: 'rs_slct'}, 
            88888);
        const team = await checkBoxGroupValue(
            {box: 'nbaTm', slct: 'tm_slct'}, 
            {box: 'wnbaTm', slct: 'wTm_slct'}, 
            0);
        await foldedLog(`%c searching random player | league: ${lg} | season ${season}`, MSG);
        let js = await getPlayerStatsV2(base, 'random', season, team, lg);
        if (js) {
            await buildPlayerDash(js.player[0], 0);
            await setPHold(js.player[0].player_meta.player);
            await scrollIntoBySize(1350, 900, "player_title");
        }
    });
}

// adds a button listener to each individual player button in the leading scorers
// tables. have to create a button, do btn.AddEventListener, and call this function
// within that listener. will insert the player's name in the search bar and call getP
export async function playerBtnListener(player) {
    let searchB = document.getElementById('pSearch');
    await clearCheckBoxes(checkBoxEls);
    if (searchB) {
        searchB.value = player;
        const season = await checkBoxGroupValue(
            {box: 'post', slct: 'ps_slct'}, 
            {box: 'reg', slct: 'rs_slct'}, 
            88888);

        const team = await checkBoxGroupValue(
            {box: 'nbaTm', slct: 'tm_slct'}, 
            {box: 'wnbaTm', slct: 'wTm_slct'}, 0);


        const lg = await lgRadioBtns();

        // search & clear player search bar
        let js = await getPlayerStatsV2(base, player, season, team, lg);
        if (js) {
            await setPHold(js.player[0].player_meta.player);
            await buildPlayerDash(js.player[0], 0);
        }
        await scrollIntoBySize(1350, 900, "player_title");
    }
}

// clear search box
export async function clearSearch() {
    const btn = document.getElementById('clearS');
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        let pSearch = document.getElementById('pSearch');
        pSearch.value = '';
        pSearch.focus();
    });
}


export async function clearSearchBar() {
    let pSearch = document.getElementById('pSearch');
    pSearch.value = '';
}

export async function lgRadioBtns() {
    const selected = document.querySelector('input[name="leagues"]:checked');
    if (selected) {
        // console.log(selected.value); // "all", "nba", or "wnba"
        return selected.value;
    } else {
        return 'all';
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

// return season or 88888 if no sseason box checked
export async function handleSeasonBoxes() {
    const p = await checkBoxes('post', 'ps_slct')
    const r = await checkBoxes('reg', 'rs_slct')

    if (p) {
        return p;
    }
    if (r) {
        return r;
    } // 2{current year} i.e. 22025 during 2025
    return 88888;
    // return `2${new Date().getFullYear()}`;
}

// return season or 88888 if no sseason box checked
// needs to be like lgrp.box and lgrp.select
export async function checkBoxGroupValue(lgrp, rgrp, dflt) {
    const l = await checkBoxes(lgrp.box, lgrp.slct);
    const r = await checkBoxes(rgrp.box, rgrp.slct);

    if (l) return l;
    if (r) return r;
    
    // 88888 for season, 0 for team
    return dflt;
    // return `2${new Date().getFullYear()}`;
}

export async function clearCheckBoxes(boxes) {
    for (let i = 0; i < boxes.length; i++) {
        let b = document.getElementById(boxes[i]);
        b.checked = 0;
    }
}

// if checkbox is checked, return the value
export async function checkBoxes(box, sel) {
    const b = document.getElementById(box);
    const s = document.getElementById(sel);
    if (b.checked) {
        return s.value
    }
}
// make post + reg checkboxes exclusive (but allow neither checked)
export async function setupExclusiveCheckboxes(leftbox, rightbox) {
    let lbox = document.getElementById(leftbox);
    let rbox = document.getElementById(rightbox);
    function handleCheck(e) {
        if (e.target.checked) {
            if (e.target === lbox) rbox.checked = false;
            if (e.target === rbox) lbox.checked = false;
        }
    }
    lbox.addEventListener("change", handleCheck);
    rbox.addEventListener("change", handleCheck);
}

// append options to select
async function makeOption(slct, txt, val) {
    let opt = document.createElement('option');
    opt.textContent = txt;
    opt.value = val;
    opt.style.width = '100%';
    slct.appendChild(opt);
}

// call seaons endpoint for the opts
export async function loadSznOptions() {
    const r = await fetch(base + '/seasons');
    if (!r.ok) { 
        throw new Error(`HTTP Error: ${s.status}`);
    } 
    const data = await r.json();
    await buildSznSelects(data);
}

// accept seasons in data object and make an option for each
async function buildSznSelects(data) {
    const rs = document.getElementById('rs_slct');
    const ps = document.getElementById('ps_slct');
    for (let s of data) {
        if (s.season_id.substring(0, 1) === '4') {
            await makeOption(ps, s.season, s.season_id);
        } else if (s.season_id.substring(0, 1) === '2') {
            await makeOption(rs, s.season, s.season_id);
        }
    }
}

// call seaons endpoint for the opts
export async function loadTeamOptions() {
    const r = await fetch(base + '/teams');
    if (!r.ok) { 
        throw new Error(`HTTP Error: ${s.status}`);
    } 
    const data = await r.json();
    await buildTeamSelects(data);
}

// accept seasons in data object and make an option for each
async function buildTeamSelects(data) {
    const nba = document.getElementById('tm_slct');
    const wnba = document.getElementById('wTm_slct');
    for (let t of data) {
        let txt = `${t.team} | ${t.team_long}`
        if (t.league === 'NBA') {
            await makeOption(nba, txt, t.team_id);
        } else if (t.league === 'WNBA') {
            await makeOption(wnba, txt, t.team_id);
        }
    }
}