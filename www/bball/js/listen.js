// script to load in HTML -- all listener functions are called here 
// import { loadSznOptions, selHvr, setupExclusiveCheckboxes, clearCheckBoxes } from "./ui.js"
import * as ui from "./ui.js"
import { makeScoringLeaders } from "./lg_ldg_scorers.js"
import { buildRGTopScorersTbl } from "./rg_ldg_scorers.js";
import { getTeamRecords, buildTeamRecsTbl } from "./teamrecs.js";
import { randPlayerBtn, searchPlayer, holdPlayerBtn, clearSearch, buildLoadDash,
    getRecentGamesData, clearSearchBar } from "./player_search.js"

export const base = "https://dev.jdeko.me/bball";
export const checkBoxes = ['post', 'reg', 'nbaTm', 'wnbaTm'];

let NUMPL = window.innerWidth <= 700 ? 5 : 10;

// onload content
document.addEventListener('DOMContentLoaded', async () => {
    await buildOnLoadElements();
    await randPlayerBtn();
    await searchPlayer();
    await clearSearch();
    await holdPlayerBtn();
});

// mobile button: jump to current player
const jump_search_btn = document.getElementById("jumptoresp");
jump_search_btn.addEventListener("click", async() => {
    const res = document.getElementById("player_title");
    if (res) {
        res.scrollIntoView({behavior: "smooth", block: "start"});
    }
});

// mobile button: jump to player search section
const jump_plr_btn = document.getElementById("jumptosearch");
jump_plr_btn.addEventListener("click", async() => {
    const res = document.getElementById("ui");
    if (res) {
        res.scrollIntoView({behavior: "smooth", block: "start"});
    }
});

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

const mq = window.matchMedia("(max-width: 700px)");

async function numPlByScreenWidth(e) {
    const numPl = e.matches ? 5 : 10;
    if (numPl !== NUMPL) {
        NUMPL = numPl;
        await makeScoringLeaders(numPl);
    }
}

// initial run
numPlByScreenWidth(mq);

// breakpoint changes only
mq.addEventListener("change", numPlByScreenWidth);

// // change number of scoring leaders when window is resized
// window.addEventListener("resize", async () => {
//     const numPl = window.innerWidth <= 700 ? 5 : 10;
    
//     if (numPl !== NUMPL) {
//         NUMPL = numPl;
//     }
//     await makeScoringLeaders(numPl);
// });

// all elements to build on load
export async function buildOnLoadElements() {
    // empty search bar on load
    await clearSearchBar();
    
    // load team records table
    let trjs = await getTeamRecords();
    await buildTeamRecsTbl(trjs, 'team_recs_tbl');

    // scoring leaders (number of players table based on screen width)
    await makeScoringLeaders(window.innerWidth <= 700 ? 5 : 10);

    // get recent games data, build player dash
    let js = await getRecentGamesData();
    await buildLoadDash(js);
    await buildRGTopScorersTbl(js, 'top_players');

    // setup season/team checkboxes
    await ui.setupExclusiveCheckboxes('post', 'reg');
    await ui.setupExclusiveCheckboxes('nbaTm', 'wnbaTm');

    // get seasons/teams from api & load options for the selects
    await ui.loadSznOptions();
    await ui.loadTeamOptions();

    // setup radio buttons for league selection
    await ui.lgRadioBtns();

    // DEFAULT VALUES: clear all checkboxes, select "Both" lg radio button
    await ui.clearCheckBoxes(checkBoxes);
    document.getElementById('all_lgs').checked = 1;
}

// place border around elements
async function border(selector, border) {
  const els = document.querySelectorAll(selector);
  if (!els.length) return;

  els.forEach(el => { 
    el.style.border = border;
  });
}
