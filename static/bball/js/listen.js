// script to load in HTML -- all listener functions are called here 
// import { loadSznOptions, selHvr, setupExclusiveCheckboxes, clearCheckBoxes } from "./ui.js"
import * as ui from "./ui.js"
import { makeScoringLeaders } from "./lg_ldg_scorers.js"
import { buildRGTopScorersTbl } from "./rg_ldg_scorers.js";
import { getTeamRecords, buildTeamRecsTbl } from "./teamrecs.js";
import { randPlayerBtn, searchPlayer, holdPlayerBtn, clearSearch, buildLoadDash,
    getRecentGamesData, clearSearchBar } from "./player_search.js"

export const base = "http://localhost:8080/bball";

document.addEventListener('DOMContentLoaded', async () => {
    await buildOnLoadElements();
    await randPlayerBtn();
    await searchPlayer();
    await clearSearch();
    await holdPlayerBtn();
});
// TODO: separate UI setup and table builds
export async function buildOnLoadElements() {
    // empty search bar on load
    await clearSearchBar();
    
    // load team records table
    let trjs = await getTeamRecords();
    await buildTeamRecsTbl(trjs, 'team_recs_tbl');

    // scoring leaders (number of players table based on screen width)
    await makeScoringLeaders(numScoringLeaders());

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
    await ui.clearCheckBoxes(['post', 'reg', 'nbaTm', 'wnbaTm']);
    document.getElementById('all_lgs').checked = 1;
}

export function numScoringLeaders() {
    if (window.innerWidth <= 700) {
        return 3;
    } else if (window.innerWidth <= 1500){
        return 5;
    } else {
        return 10;
    }
}

// place border around elements
async function border(selector, border) {
  const els = document.querySelectorAll(selector);
  if (!els.length) return;

  els.forEach(el => { 
    el.style.border = border;
  });
}
