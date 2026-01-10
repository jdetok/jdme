import { clearSearchBtn, clearSearch, lgRadioBtns, setPHold } from "./btns.js";
import { base, mediaQueryBreak, checkBoxEls, scrollIntoBySize, MSG, foldedLog, MSG_BOLD } from "../global.js";
import { checkBoxGroupValue, clearCheckBoxes, setupExclusiveCheckboxes } from "./checkbox.js";
import { fetchPlayer, buildPlayerDash, buildLoadDash } from "./player.js";
import { makeLgTopScorersTbl, makeRgTopScorersTbl, makeTeamRecordsTbl } from "./tbl.js";
import { loadSznOptions, loadTeamOptions } from "./selectors.js";

let NUMROWS = window.innerWidth <= 700 ? 5 : 10;

export async function LoadContent(): Promise<void> {
    document.addEventListener('DOMContentLoaded', async () => {
        foldedLog(`%cloading content for page {${window.innerWidth}px x ${window.innerHeight}px}...`, MSG_BOLD)
        clearSearch();

        await buildOnLoadElements();
        await listenForWindowSize(mediaQueryBreak);

        await clearSearchBtn();
        await searchPlayer();
        await randPlayerBtn();
        await holdPlayerBtn();
        
    });
}

export async function buildOnLoadElements() {
    foldedLog(`%c building page load elements for page width ${window.innerWidth}`, MSG);
    const rows_on_load = window.innerWidth <= 700 ? 5 : 10;
    await makeLgTopScorersTbl(rows_on_load);
    await makeRgTopScorersTbl(rows_on_load);
    await makeTeamRecordsTbl(rows_on_load);
    await setup_jump_btns();
    await expandListBtns();

    await setupExclusiveCheckboxes('post', 'reg');
    await setupExclusiveCheckboxes('nbaTm', 'wnbaTm');

    // get seasons/teams from api & load options for the selects
    foldedLog(`%c loading team/season selectors data...`, MSG)
    await loadSznOptions();
    await loadTeamOptions();

    await lgRadioBtns();

    await buildLoadDash();
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
        let js = await fetchPlayer(base, player, season, team, lg);
        if (js) {
            await setPHold(js.player[0].player_meta.player);
            await buildPlayerDash(js.player[0], 0);
        }
        scrollIntoBySize(1350, 900, "player_title");
    }
}

export async function holdPlayerBtn() {
    // listen for hold player button press
    const elId = 'holdP';
    const btn = document.getElementById(elId);
    if (!btn) throw new Error(`couldn't get button element at ${elId}`);
    btn.addEventListener('click', async (event) => {
        event.preventDefault();

        // get player name held in pHold value, fill player search bar with it
        const holdElId = 'pHold';
        const hold = document.getElementById(holdElId) as HTMLInputElement;
        if (!hold) throw new Error(`couldn't get input element at ${holdElId}`);
        let player = hold.value;

        const searchElId = 'pSearch';
        let search = document.getElementById(searchElId) as HTMLInputElement;
        search.value = player;
    })
}


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

async function expandListBtns(btns = [
    {elId: "seemoreplayers", run: increase, build: makeLgTopScorersTbl}, 
    {elId: "seelessplayers", run: decrease, build: makeLgTopScorersTbl},
    {elId: "seemoreRGplayers", run: increase, build: makeRgTopScorersTbl}, 
    {elId: "seelessRGplayers", run: decrease, build: makeRgTopScorersTbl},
]) {
    for (let btnObj of btns) {
        const btn = document.getElementById(btnObj.elId);
        if (btn) {
            btn.addEventListener('click', async() => {
                let newNum = await btnObj.run(NUMROWS, 5);
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
    const numRows = e.matches ? 5 : 10;
    if (numRows !== NUMROWS) {
        NUMROWS = numRows;
    }
    await makeLgTopScorersTbl(numRows);
    await makeRgTopScorersTbl(numRows);
    await makeTeamRecordsTbl(numRows);
}

async function listenForWindowSize(width) {
    const mq = window.matchMedia(`(max-width: ${width}px)`);
    mq.addEventListener("change", async (e) => await numRowsByScreenWidth(e));
}

export async function searchPlayer() {
    // listen for form submission
    const elId = 'ui';
    const frm = document.getElementById(elId) as HTMLFormElement;
    if (!frm) throw new Error(`couldn't get element at Id ${elId}`);
    frm.addEventListener('submit', async (event) => {
        event.preventDefault();

        // get value of player search box
        const searchElId = 'pSearch';
        const input = document.getElementById(searchElId) as HTMLInputElement;
        if (!input) throw new Error(`couldn't get element at Id ${searchElId}`);
        
        const lg = await lgRadioBtns();
        
        let player = input.value.trim();
        
        // if search pressed without anything in search box, searches current player
        if (player === '') {
            const pHoldElId = 'pHold';
            const pHoldEl =  document.getElementById(pHoldElId) as HTMLInputElement;
            player = pHoldEl.value;
        }

        // check if season box is checked, return sel val if so, 88888 if not
        // 88888 gets the most recent season from the api
        // const season = await handleSeasonBoxes();
        const season = await checkBoxGroupValue(
            {box: 'post', slct: 'ps_slct'}, 
            {box: 'reg', slct: 'rs_slct'}, 
            88888);

        const team = await checkBoxGroupValue(
            {box: 'nbaTm', slct: 'tm_slct'}, 
            {box: 'wnbaTm', slct: 'wTm_slct'}, 
            0);

        foldedLog(`%csearching for player ${player} | season ${season} | team ${team} | league ${lg}`, MSG_BOLD);
        // console.trace();
        // build response player dash section
        let js = await fetchPlayer(base, player, season, team, lg);
        if (js) {
            await setPHold(js.player[0].player_meta.player);
            await buildPlayerDash(js.player[0], 0);
        }
    }) 
}

// get a random player from the API and getPlayerStats
export async function randPlayerBtn() {
    // listen for random player button press
    const elId = 'randP';
    const btn = document.getElementById(elId) as HTMLInputElement;
    if (!btn) throw new Error(`couldn't get button element at id ${elId}`);
    btn.addEventListener('click', async (event) => {        
        event.preventDefault();
        const lg = await lgRadioBtns();
        // check season boxes & get appropriate season id, search with random as player
        const season = await checkBoxGroupValue(
            {box: 'post', slct: 'ps_slct'}, 
            {box: 'reg', slct: 'rs_slct'}, 
            88888);
        const team = await checkBoxGroupValue(
            {box: 'nbaTm', slct: 'tm_slct'}, 
            {box: 'wnbaTm', slct: 'wTm_slct'}, 
            0);
        foldedLog(`%c searching random player | league: ${lg} | season ${season}`, MSG);
        let js = await fetchPlayer(base, 'random', season, team, lg);
        if (js) {
            await buildPlayerDash(js.player[0], 0);
            await setPHold(js.player[0].player_meta.player);
            scrollIntoBySize(1350, 900, "player_title");
        }
    });
}