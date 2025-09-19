import { makePlayerDash } from "./resp.js";
import { showHideHvr } from "./hover.js";
import { base } from "./listen.js";

// get player from search bar and make player dash
export async function searchPlayer() {
    const frm = document.getElementById('ui');
    frm.addEventListener('submit', async (event) => {
        event.preventDefault();

        const input = document.getElementById('pSearch');
        let player = input.value.trim();
        if (player === '') {
            player = document.getElementById('pHold').value;
        }
        const season = await handleSeasonBoxes();
        console.log(`searching for season ${season}`)
        await makePlayerDash(base, player, season, '0');
        input.value = ''; // clear input box after searching
    }) 
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
        searchB.value = player;
        // searchB.focus();
        await makePlayerDash(base, player, 88888, 0, 0);
        searchB.value = '';
        let res = document.getElementById("ui");
        if (res) {
            res.scrollIntoView({behavior: "smooth", block: "start"});
        }
    }
}


// read pHold invisible val to add on-screen player's name to search bar
export async function holdPlayerBtn() {
    const btn = document.getElementById('holdP');
    const hlp = document.getElementById('hlpHld');
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        let player = document.getElementById('pHold').value;
        document.getElementById('pSearch').value = player;
    })
    await showHideHvr(
        hlp, 
        'hvrmsg',
        `fill the input box with the current player's name`
    )
}

/*
get a random player from the API and makePlayerDash
*/
export async function randPlayerBtn() {
    const btn = document.getElementById('randP');
    const hlp = document.getElementById('hlpRnd');
    btn.addEventListener('click', async (event) => {        
        event.preventDefault();
        const season = await handleSeasonBoxes();
        console.log(season);
        await makePlayerDash(base, 'random', season, 0);
    })
    await showHideHvr(
        hlp, 
        'hvrmsg',
        `get the stats for a random player in the selected season. if no season 
        is specified, the current/most recent season will be used. if the 
        random player did not play in the selected season, their most 
        recent (or first, whichever is closer) season will be used`
    )
}

// season select hover messages
export async function selHvr() {
    const rs = document.getElementById('hlpRs');
    const ps = document.getElementById('hlpPs');
    await showHideHvr(rs, 'selhvr',
        `search for a specific regular-season. if the player being searched didn't
        play in the selected season, their first or most recent season, whichever
        is closer to the selected, will be used`);
    await showHideHvr(ps, 'selhvr',
        `search for a specific post-season. if the player being searched didn't
        play in the selected season, their first or most recent season, whichever
        is closer to the selected, will be used`);
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

// append options to select
async function makeOption(slct, txt, val) {
    let opt = document.createElement('option');
    opt.textContent = txt;
    opt.value = val;
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

// if checkbox is checked, return the value
export async function checkBoxes(box, sel) {
    const b = document.getElementById(box);
    const s = document.getElementById(sel);
    if (b.checked) {
        return s.value
    }
}