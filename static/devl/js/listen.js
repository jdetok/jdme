// script to load in HTML -- all listener functions are called here 

import * as home from "./home.js";
import * as buttons from "./buttons.js"

export const base = "https://jdeko.me/bball";
export const dev = "https://jdeko.me/devl/bball";

document.addEventListener('DOMContentLoaded', async () => {
    await home.loadSeasonOpts();
    await home.loadTeamOpts();
    await home.gamesRecent();
    await home.topScorer();
    await home.lgChangeListener();
    await buttons.clear();
    await buttons.randPlayerBtn();
    await buttons.search();
});
