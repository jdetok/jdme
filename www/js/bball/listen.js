// script to load in HTML -- all listener functions are called here 
// import { loadSznOptions, selHvr, setupExclusiveCheckboxes, clearCheckBoxes } from "./ui.js"
import * as ui from "./ui.js"
import { makeScoringLeaders, makeRGTopScorersTbl, makeTeamRecsTable } from "./tables_onload.js"
import { searchPlayer, buildLoadDash, getRecentGamesData } from "./player_search.js"
import { checkBoxEls, AQUA_BOLD, AQUA } from "./util.js";

let NUMPL = window.innerWidth <= 700 ? 5 : 10;

// onload content
document.addEventListener('DOMContentLoaded', async () => {
    console.log('%c loading page...', 'color: green; font-weight: bold;')
    await buildOnLoadElements();
    await searchPlayer();
    await ui.randPlayerBtn();
    await ui.clearSearch();
    await ui.holdPlayerBtn();
});

export async function setup_jump_btns() {
    const btns = [{el: "jumptoresp", jumpTo: "player_title"}, {el: "jumptosearch", jumpTo: "ui"}]
    for (const btn of btns) {
        const btnEl = document.getElementById(btn.el);
        if (btnEl) {
            btnEl.addEventListener('click', async() => {
                const jmpEl = document.getElementById(btn.jumpTo);
                if (jmpEl) {
                    jmpEl.scrollIntoView({behavior: "smooth", block: "start"});
                }
            })
        }    
    }
}

// mobile button: see more top players
const seemore_btn = document.getElementById("seemoreplayers");
seemore_btn.addEventListener("click", async() => {
    if (!NUMPL) {
        NUMPL = window.innerWidth <= 700 ? 5 : 10;
    };
    NUMPL += 5;
    await makeScoringLeaders(NUMPL);
});

// mobile button: see more top players
const seeless_btn = document.getElementById("seelessplayers");
seeless_btn.addEventListener("click", async() => {
    NUMPL = window.innerWidth <= 700 ? 5 : 10;;
    await makeScoringLeaders(NUMPL);
});


async function numPlByScreenWidth(e) {
    const numPl = e.matches ? 5 : 10;
    if (numPl !== NUMPL) {
        NUMPL = numPl;
        await makeScoringLeaders(numPl);
    }
}

// initial run
const mq = window.matchMedia("(max-width: 700px)");
numPlByScreenWidth(mq);

// breakpoint changes only
mq.addEventListener("change", numPlByScreenWidth);

// all elements to build on load
export async function buildOnLoadElements() {
    console.trace(`%c building page load elements for page width ${window.innerWidth}`, AQUA_BOLD)
    const rows_on_load = window.innerWidth <= 700 ? 5 : 10
    // empty search bar on load
    await ui.clearSearchBar();

    await setup_jump_btns();
    await makeTeamRecsTable(rows_on_load);

    // scoring leaders (number of players table based on screen width)
    await makeScoringLeaders(rows_on_load);

    // get recent games data, build player dash

    
    let js = await getRecentGamesData();
    console.trace(`%c fetched games data for ${js.recent_games[0].game_date}`, AQUA)
    await buildLoadDash(js);

    await makeRGTopScorersTbl(js, rows_on_load);
    // await buildRGTopScorersTbl(js, 'top_players');

    // setup season/team checkboxes
    await ui.setupExclusiveCheckboxes('post', 'reg');
    await ui.setupExclusiveCheckboxes('nbaTm', 'wnbaTm');

    // get seasons/teams from api & load options for the selects
    console.trace(`%c loading team/season selectors data...`, AQUA)
    await ui.loadSznOptions();
    await ui.loadTeamOptions();

    // setup radio buttons for league selection
    await ui.lgRadioBtns();

    // DEFAULT VALUES: clear all checkboxes, select "Both" lg radio button
    await ui.clearCheckBoxes(checkBoxEls);
    document.getElementById('all_lgs').checked = 1;
}