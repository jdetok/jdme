// script to load in HTML -- all listener functions are called here 
// import { loadSznOptions, selHvr, setupExclusiveCheckboxes, clearCheckBoxes } from "./ui.js"
import * as ui from "./ui.js"
import { makeScoringLeaders, makeRGTopScorersTbl, makeTeamRecsTable } from "./tables_onload.js"
import { searchPlayer, buildLoadDash, getRecentGamesData } from "./player_search.js"
import { checkBoxEls, MSG, foldedLog, MSG_BOLD, RED_BOLD } from "./util.js";

let NUMPL = window.innerWidth <= 700 ? 5 : 10;

// onload content
document.addEventListener('DOMContentLoaded', async () => {
    await foldedLog('%c loading page...', MSG)
    await buildOnLoadElements();
    await searchPlayer();
    await ui.randPlayerBtn();
    await ui.clearSearch();
    await ui.holdPlayerBtn();
    console.log(`height: ${window.innerHeight}`);
});

export async function setup_jump_btns() {
    const btns = [{el: "jumptoresp", jumpTo: "player_title"}, {el: "jumptosearch", jumpTo: "ui"}]
    for (const btn of btns) {
        const btnEl = document.getElementById(btn.el);
        if (btnEl) {
            btnEl.addEventListener('click', async() => {
                const jmpEl = document.getElementById(btn.jumpTo);
                if (jmpEl) {
                    jmpEl.scrollIntoView({behavior: "smooth", block: "start"});
                }
            })
        }    
    }
}

async function expandListBtns(data, btns = [
    {elId: "seemoreplayers", run: increase, build: makeScoringLeaders}, 
    {elId: "seelessplayers", run: decrease, build: makeScoringLeaders},
    {elId: "seemoreRGplayers", run: increase, build: makeRGTopScorersTbl}, 
    {elId: "seelessRGplayers", run: decrease, build: makeRGTopScorersTbl},
]) {
    for (let btnObj of btns) {
        const btn = document.getElementById(btnObj.elId);
        if (btn) {
            btn.addEventListener('click', async() => {
                let newNum = await btnObj.run(NUMPL, 5);
                console.log(`%c${newNum}`, RED_BOLD)
                if (btnObj.build === makeScoringLeaders) {
                    await makeScoringLeaders(newNum);
                }
                if (btnObj.build === makeRGTopScorersTbl) {
                    await makeRGTopScorersTbl(data, newNum);
                }
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

// mobile button: see more top players
// const seeless_btn = document.getElementById("seelessplayers");
// seeless_btn.addEventListener("click", async() => {
//     NUMPL = window.innerWidth <= 700 ? 5 : 10;;
//     await makeScoringLeaders(NUMPL);
// });


async function numPlByScreenWidth(e, data) {
    const numPl = e.matches ? 5 : 10;
    if (numPl !== NUMPL) {
        NUMPL = numPl;
    }
    await makeScoringLeaders(numPl);
    await makeRGTopScorersTbl(data, numPl);
}

async function listenForWindowSize(data, width) {
    const mq = window.matchMedia(`(max-width: ${width}px)`);
    await numPlByScreenWidth(mq, data);
    mq.addEventListener("change", async (e) => await numPlByScreenWidth(e, data));
}

// all elements to build on load
export async function buildOnLoadElements() {
    await foldedLog(`%c building page load elements for page width ${window.innerWidth}`, MSG);
    const rows_on_load = window.innerWidth <= 700 ? 5 : 10;

    await ui.clearSearchBar();
    await setup_jump_btns();
    await makeTeamRecsTable(rows_on_load);
    await makeScoringLeaders(rows_on_load);
    
    try {
        let js = await getRecentGamesData();
        await foldedLog(`%c fetched games data for ${js.recent_games[0].game_date}`, MSG);
        // await makeRGTopScorersTbl(js, rows_on_load);
        await listenForWindowSize(js, 850);
        await expandListBtns(js);
        await buildLoadDash(js);
    } catch(err) {
        await foldedLog(`%cerror fetching recent games data: ${err}`, RED_BOLD);
    }
    
    

    // setup season/team checkboxes
    await ui.setupExclusiveCheckboxes('post', 'reg');
    await ui.setupExclusiveCheckboxes('nbaTm', 'wnbaTm');

    // get seasons/teams from api & load options for the selects
    await foldedLog(`%c loading team/season selectors data...`, MSG)
    await ui.loadSznOptions();
    await ui.loadTeamOptions();

    // setup radio buttons for league selection
    await ui.lgRadioBtns();

    // DEFAULT VALUES: clear all checkboxes, select "Both" lg radio button
    await ui.clearCheckBoxes(checkBoxEls);
    document.getElementById('all_lgs').checked = 1;
    await foldedLog(`%cpage load compelete`, MSG_BOLD);
}