// script to load in HTML -- all listener functions are called here 
// import { loadSznOptions, selHvr, setupExclusiveCheckboxes, clearCheckBoxes } from "./ui.js"
import * as ui from "./ui.js"
import { makeScoringLeaders } from "./lg_ldg_scorers.js"
import { makeRGTopScorers } from "./rg_ldg_scorers.js"
import { randPlayerBtn, searchPlayer, holdPlayerBtn, clearSearch} from "./player_search.js"

export const base = "https://jdeko.me/bball";

document.addEventListener('DOMContentLoaded', async () => {
    await ui.setupExclusiveCheckboxes('post', 'reg');
    await ui.setupExclusiveCheckboxes('nbaTm', 'wnbaTm');
    await ui.clearCheckBoxes(['post', 'reg', 'nbaTm', 'wnbaTm']);
    await ui.loadSznOptions();
    await ui.loadTeamOptions();
    await ui.lgRadioBtns();
    await ui.selHvr();
    await randPlayerBtn();
    await searchPlayer();
    await clearSearch();
    await holdPlayerBtn();
    await makeRGTopScorers();
    
    await loadScoringLeaders();
    document.getElementById('all_lgs').checked = 1;
    // await border('#topsect > div', '2px solid black');
});

export async function loadScoringLeaders() {
    let numPl;
    if (window.innerWidth <= 700) {
        numPl = 3;
    } else {
        numPl = 10;
    }
    await makeScoringLeaders(numPl);
}

// place border around elements
async function border(selector, border) {
  const els = document.querySelectorAll(selector);
  if (!els.length) return;

  els.forEach(el => { 
    el.style.border = border;
  });
}
