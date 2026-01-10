import { clearSearch, lgRadioBtns, setPHold } from "./btns.js";
import { base, checkBoxEls, scrollIntoBySize, MSG, foldedLog, MSG_BOLD, RED_BOLD } from "../global.js";
import { checkBoxGroupValue, clearCheckBoxes } from "./checkbox.js";

let NUMPL = window.innerWidth <= 700 ? 5 : 10;


export async function LoadContent(): Promise<void> {
    document.addEventListener('DOMContentLoaded', async () => {
        foldedLog(`%cloading content for page {${window.innerWidth}px x ${window.innerHeight}px}...`, MSG_BOLD)
        // await buildOnLoadElements();
        // await searchPlayer();
        // await ui.randPlayerBtn();
        clearSearch();
        // await ui.holdPlayerBtn();
    });
}

// adds a button listener to each individual player button in the leading scorers
// tables. have to create a button, do btn.AddEventListener, and call this function
// within that listener. will insert the player's name in the search bar and call getP
export async function playerBtnListener(player) {
    let searchB = document.getElementById('pSearch');
    await clearCheckBoxes(checkBoxEls);
    if (searchB) {
        // searchB.value = player;
        const season = await checkBoxGroupValue(
            {box: 'post', slct: 'ps_slct'}, 
            {box: 'reg', slct: 'rs_slct'}, 
            88888);

        const team = await checkBoxGroupValue(
            {box: 'nbaTm', slct: 'tm_slct'}, 
            {box: 'wnbaTm', slct: 'wTm_slct'}, 0);


        const lg = await lgRadioBtns();

        // search & clear player search bar
        let js = await getPlayerStatsV2(base, player, season, team, lg);
        if (js) {
            await setPHold(js.player[0].player_meta.player);
            await buildPlayerDash(js.player[0], 0);
        }
        scrollIntoBySize(1350, 900, "player_title");
    }
}