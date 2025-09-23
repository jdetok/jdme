import { base } from "./listen.js"
import { showHideHvr } from "./hover.js";
import { handleSeasonBoxes } from "./ui.js";
import { makePlayerDash } from "./player_dash.js";

// get player from search bar and make player dash
export async function searchPlayer() {
    // listen for form submission
    const frm = document.getElementById('ui');
    frm.addEventListener('submit', async (event) => {
        event.preventDefault();

        // get value of player search box
        const input = document.getElementById('pSearch');
        let player = input.value.trim();

        // if search pressed without anything in search box, searches current player
        if (player === '') {
            player = document.getElementById('pHold').value;
        }

        // check if season box is checked, return sel val if so, 88888 if not
        // 88888 gets the most recent season from the api
        const season = await handleSeasonBoxes();
        console.log(`searching for season ${season}`)

        // build response player dash section
        await makePlayerDash(base, player, season, 0, 0);

        // clear player search box
        input.value = ''; // clear input box after searching
    }) 
}

// get a random player from the API and makePlayerDash
export async function randPlayerBtn() {
    // listen for random player button press
    const btn = document.getElementById('randP');
    btn.addEventListener('click', async (event) => {        
        event.preventDefault();

        // check season boxes & get appropriate season id, search with random as player
        const season = await handleSeasonBoxes();
        console.log(`searching random player for season ${season}`);
        await makePlayerDash(base, 'random', season, 0, 0);
    })

    // hover message for help ?
    const hlp = document.getElementById('hlpRnd');
    await showHideHvr(
        hlp, 
        'hvrmsg',
        `get the stats for a random player in the selected season. if no season 
        is specified, the current/most recent season will be used. if the 
        random player did not play in the selected season, their most 
        recent (or first, whichever is closer) season will be used`
    )
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
        const season = await handleSeasonBoxes();

        // search & clear player search bar
        await makePlayerDash(base, player, season, 0, 0);
        searchB.value = '';

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

    // help button hover val
    const hlp = document.getElementById('hlpHld');
    await showHideHvr(
        hlp, 
        'hvrmsg',
        `fill the input box with the current player's name`
    )
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