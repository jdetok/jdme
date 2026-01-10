import { clearSearch, lgRadioBtns, setPHold } from "./btns.js";
import { base, checkBoxEls, scrollIntoBySize, MSG, foldedLog, MSG_BOLD, RED_BOLD } from "../global.js";
import { checkBoxGroupValue, clearCheckBoxes } from "./checkbox.js";
import { fetchPlayer, buildPlayerDash } from "./player.js";
import { makeLgTopScorersTbl, makeRgTopScorersTbl, makeTeamRecordsTbl } from "./tbls_onload.js";
let NUMROWS = window.innerWidth <= 700 ? 5 : 10;
export async function LoadContent() {
    document.addEventListener('DOMContentLoaded', async () => {
        foldedLog(`%cloading content for page {${window.innerWidth}px x ${window.innerHeight}px}...`, MSG_BOLD);
        await buildOnLoadElements();
        // await searchPlayer();
        // await ui.randPlayerBtn();
        clearSearch();
        // await ui.holdPlayerBtn();
        await listenForWindowSize(NUMROWS);
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
        const season = await checkBoxGroupValue({ box: 'post', slct: 'ps_slct' }, { box: 'reg', slct: 'rs_slct' }, 88888);
        const team = await checkBoxGroupValue({ box: 'nbaTm', slct: 'tm_slct' }, { box: 'wnbaTm', slct: 'wTm_slct' }, 0);
        const lg = await lgRadioBtns();
        // search & clear player search bar
        let js = await fetchPlayer(base, player, season, team, lg, 'sErr');
        if (js) {
            await setPHold(js.player[0].player_meta.player);
            await buildPlayerDash(js.player[0], 0);
        }
        scrollIntoBySize(1350, 900, "player_title");
    }
}
export async function setup_jump_btns() {
    const btns = [{ el: "jumptoresp", jumpTo: "player_title" }, { el: "jumptosearch", jumpTo: "ui" }];
    for (const btn of btns) {
        const btnEl = document.getElementById(btn.el);
        if (btnEl) {
            btnEl.addEventListener('click', async () => {
                const jmpEl = document.getElementById(btn.jumpTo);
                if (jmpEl) {
                    jmpEl.scrollIntoView({ behavior: "smooth", block: "start" });
                }
            });
        }
    }
}
async function expandListBtns(btns = [
    { elId: "seemoreplayers", run: increase, build: makeLgTopScorersTbl },
    { elId: "seelessplayers", run: decrease, build: makeLgTopScorersTbl },
    { elId: "seemoreRGplayers", run: increase, build: makeRgTopScorersTbl },
    { elId: "seelessRGplayers", run: decrease, build: makeRgTopScorersTbl },
]) {
    for (let btnObj of btns) {
        const btn = document.getElementById(btnObj.elId);
        if (btn) {
            btn.addEventListener('click', async () => {
                let newNum = await btnObj.run(NUMROWS, 5);
                console.log(`%c${newNum}`, RED_BOLD);
                await btnObj.build(newNum);
            });
        }
    }
}
async function increase(num, by) {
    return num += by;
}
async function decrease(num, by) {
    return num -= by;
}
async function numRowsByScreenWidth(e) {
    const numPl = e.matches ? 5 : 10;
    if (numPl !== NUMROWS) {
        NUMROWS = numPl;
    }
    await makeLgTopScorersTbl(numPl);
    await makeRgTopScorersTbl(numPl);
    await makeTeamRecordsTbl(numPl);
}
async function listenForWindowSize(width) {
    const mq = window.matchMedia(`(max-width: ${width}px)`);
    mq.addEventListener("change", async (e) => await numRowsByScreenWidth(e));
}
export async function buildOnLoadElements() {
    foldedLog(`%c building page load elements for page width ${window.innerWidth}`, MSG);
    const rows_on_load = window.innerWidth <= 700 ? 5 : 10;
    await makeLgTopScorersTbl(rows_on_load);
    await makeRgTopScorersTbl(rows_on_load);
    await makeTeamRecordsTbl(rows_on_load);
    await setup_jump_btns();
}
//# sourceMappingURL=listeners.js.map