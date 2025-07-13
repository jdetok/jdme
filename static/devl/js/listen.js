// script to load in HTML -- all listener functions are called here 

import * as home from "./home.js";
import * as buttons from "./buttons.js"
import * as selectors from "./selectors.js"

export const base = "https://jdeko.me/bball";
export const dev = "https://jdeko.me/devl/bball";

export let crnt = "first";
export async function updateCrnt(new_crnt) {
    console.log(`current pre test: ${crnt}`)
    crnt = new_crnt;
    console.log(`current post test: ${crnt}`);
}

document.addEventListener('DOMContentLoaded', async () => {
    // await home.loadSeasonOpts();
    await selectors.loadSznOptions();
    await selectors.loadAllTeamOpts();
    // await home.loadTeamOpts();
    // await home.gamesRecent();
    // await home.topScorer();
    // await home.lgChangeListener();
    await buttons.clear();
    await buttons.randPlayerBtn();
    await buttons.search();
    await buttons.clearSearch();
    await buttons.holdPlayerBtn();
});
