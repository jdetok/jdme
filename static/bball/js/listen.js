// script to load in HTML -- all listener functions are called here 
import { loadSznOptions, selHvr, setupExclusiveCheckboxes } from "./ui.js"
import { makeScoringLeaders } from "./lg_ldg_scorers.js"
import { makeRGTopScorers } from "./rg_ldg_scorers.js"
import { randPlayerBtn, searchPlayer, holdPlayerBtn, clearSearch} from "./player_search.js"

export const base = "http://localhost:8080/bball";

document.addEventListener('DOMContentLoaded', async () => {
    await loadSznOptions();
    await selHvr();
    await randPlayerBtn();
    await searchPlayer();
    await clearSearch();
    await holdPlayerBtn();
    await makeRGTopScorers();
    await setupExclusiveCheckboxes();

    let numPl;
    if (window.innerWidth <= 700) {
        numPl = 3;
    } else {
        numPl = 10;
    }
    await makeScoringLeaders(numPl);

    const seeMoreBtn = document.getElementById("seemore");
    if (seeMoreBtn) {
        seeMoreBtn.addEventListener("click", seeMoreLeaders);
    }
});

export async function seeMoreLeaders() {
    await makeScoringLeaders(10);
}