// script to load in HTML -- all listener functions are called here 
// import { loadSznOptions, selHvr, setupExclusiveCheckboxes, clearCheckBoxes } from "./ui.js"
import * as ui from "./ui.js"
import { makeScoringLeaders } from "./lg_ldg_scorers.js"
import { makeRGTopScorers } from "./rg_ldg_scorers.js"
import { randPlayerBtn, searchPlayer, holdPlayerBtn, clearSearch} from "./player_search.js"

export const base = "https://jdeko.me/bball";

document.addEventListener('DOMContentLoaded', async () => {
    await ui.loadSznOptions();
    await ui.selHvr();
    await randPlayerBtn();
    await searchPlayer();
    await clearSearch();
    await holdPlayerBtn();
    await makeRGTopScorers();
    await ui.setupExclusiveCheckboxes();
    await ui.clearCheckBoxes('post');
    await ui.clearCheckBoxes('reg');
    await ui.lgRadioBtns();
    await loadScoringLeaders();
    document.getElementById('all_lgs').checked = 1;
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