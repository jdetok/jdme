// script to load in HTML -- all listener functions are called here 

// import * as home from "./home.js";
import * as buttons from "./buttons.js"
import * as pdash from "./pdash.js"
import * as selectors from "./selectors.js"

export const base = "https://jdeko.me/bball";
// export const dev = "https://jdeko.me/devl/bball";

export async function getRecGames() {
    const r = await fetch(`${base}/games/recent`);
    if (!r.ok) {
        throw new Error(`${r.status}: error calling /games/recent`);
    }
    const data = await r.json();
    const player = data.top_scorers[0].player_id;
    await pdash.getP(base, player, 88888, 0, data);
}

document.addEventListener('DOMContentLoaded', async () => {
    await selectors.loadSznOptions();
    await selectors.loadAllTeamOpts();
    await buttons.randPlayerBtn();
    await buttons.search();
    await buttons.clearSearch();
    await buttons.holdPlayerBtn();
    // await pdash.getP(base, 'random', 88888, 0);

    await getRecGames();
});