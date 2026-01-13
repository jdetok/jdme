import { makeLgTopScorersTbl, makeRgTopScorersTbl, makeTeamRecordsTbl } from "./tbl.js";
import { base } from "../global.js";
const WINDOWSIZE = 700;
let exBtnsInitComplete = false;
// counter class for number of rows displayed per table
class rowNum {
    val;
    min;
    constructor(val, min = 2) {
        this.val = val;
        this.min = min;
    }
    ;
    get value() {
        return this.val;
    }
    increase = (by = 5) => {
        this.val += by;
        return this.val;
    };
    decrease = (by = 5) => {
        this.val = Math.max(this.min, this.val - by);
        return this.val;
    };
    reset = (to) => {
        this.val = Math.max(this.min, to);
        return this.val;
    };
}
;
// tracks row counters for both top scorer tables
export class rowsState {
    lgRowNum;
    rgRowNum;
    startRows;
    constructor(winSize = WINDOWSIZE) {
        this.startRows = window.innerWidth <= winSize ? 5 : 10;
        this.lgRowNum = new rowNum(this.startRows);
        this.rgRowNum = new rowNum(this.startRows);
        this.listenForResize();
    }
    ;
    resetRows(winSize = WINDOWSIZE) {
        const rows = window.innerWidth <= winSize ? 5 : 10;
        this.lgRowNum.reset(rows);
        this.rgRowNum.reset(rows);
        return rows;
    }
    listenForResize() {
        const mq = window.matchMedia(`(max-width: ${WINDOWSIZE}px)`);
        mq.addEventListener('change', async () => {
            const newRows = this.resetRows();
            await Promise.all([
                makeTeamRecordsTbl(newRows),
                makeLgTopScorersTbl(this.lgRowNum.value),
                makeRgTopScorersTbl(this.rgRowNum.value)
            ]);
        });
    }
}
;
export async function expandedListBtns(rs, btns = [
    { elId: "seemoreplayers", rows: rs.lgRowNum, pm: '+', build: makeLgTopScorersTbl },
    { elId: "seelessplayers", rows: rs.lgRowNum, pm: '-', build: makeLgTopScorersTbl },
    { elId: "seemoreRGplayers", rows: rs.rgRowNum, pm: '+', build: makeRgTopScorersTbl },
    { elId: "seelessRGplayers", rows: rs.rgRowNum, pm: '-', build: makeRgTopScorersTbl },
]) {
    if (exBtnsInitComplete)
        return;
    exBtnsInitComplete = true;
    for (let btnObj of btns) {
        const btn = document.getElementById(btnObj.elId);
        if (!btn)
            continue;
        btn.addEventListener('click', async () => {
            const newNum = btnObj.pm === '+' ? btnObj.rows.increase() : btnObj.rows.decrease();
            await btnObj.build(newNum);
        });
    }
}
export async function checkBoxGroupValue(lgrp, rgrp, dflt) {
    const l = await checkBoxes(lgrp.box, lgrp.slct);
    const r = await checkBoxes(rgrp.box, rgrp.slct);
    if (l)
        return l;
    if (r)
        return r;
    return String(dflt);
}
export async function checkBoxes(box, sel) {
    const b = document.getElementById(box);
    const s = document.getElementById(sel);
    if (!b || !s) {
        throw new Error(`couldn't get element with id ${box} or ${sel}`);
    }
    if (b.checked) {
        return s.value;
    }
}
export async function clearCheckBoxes(boxes) {
    for (let i = 0; i < boxes.length; i++) {
        let b = document.getElementById(boxes[i]);
        if (!b)
            throw new Error(`couldn't find input element with id ${boxes[i]}`);
        b.checked = false;
    }
}
// make post + reg checkboxes exclusive (but allow neither checked)
export async function setupExclusiveCheckboxes(leftbox, rightbox) {
    let lbox = document.getElementById(leftbox);
    let rbox = document.getElementById(rightbox);
    if (!lbox || !rbox)
        throw new Error(`couldn't get ${lbox} or ${rbox}`);
    function handleCheck(e) {
        const target = e.target;
        if (!target)
            throw new Error(`can't find ${e.type}`);
        if (target.checked) {
            if (e.target === lbox)
                rbox.checked = false;
            if (e.target === rbox)
                lbox.checked = false;
        }
    }
    lbox.addEventListener("change", handleCheck);
    rbox.addEventListener("change", handleCheck);
}
export async function clearSearchBtn() {
    const btn = document.getElementById('clearS');
    if (!btn)
        return;
    btn.addEventListener('click', (event) => {
        event.preventDefault();
        clearSearch(true);
    });
}
// BUTTONS SECTION
export async function holdPlayerBtn(elId = 'holdP') {
    const btn = document.getElementById(elId);
    if (!btn)
        throw new Error(`couldn't get button element at ${elId}`);
    btn.addEventListener('click', async (event) => {
        event.preventDefault();
        // get player name held in pHold value, fill player search bar with it
        const holdElId = 'pHold';
        const hold = document.getElementById(holdElId);
        if (!hold)
            throw new Error(`couldn't get input element at ${holdElId}`);
        let player = hold.value;
        if (player === '') {
            console.error(`%chold button pressed, empty string in ${holdElId}`);
            return;
        }
        const searchElId = 'pSearch';
        let search = document.getElementById(searchElId);
        search.value = player;
    });
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
export function clearSearch(focus = false) {
    const elId = 'pSearch';
    const pSearch = document.getElementById(elId);
    if (!pSearch)
        throw new Error(`couldn't get search bar element at ${elId}`);
    pSearch.value = '';
    if (focus)
        pSearch.focus();
}
export async function setPHold(player) {
    const ph = document.getElementById('pHold');
    ph.value = player;
}
export async function lgRadioBtns() {
    const selected = document.querySelector('input[name="leagues"]:checked');
    if (selected) {
        return selected.value;
    }
    else {
        return 'all';
    }
}
export async function getInputVals() {
    const lg = await lgRadioBtns();
    const season = await checkBoxGroupValue({ box: 'post', slct: 'ps_slct' }, { box: 'reg', slct: 'rs_slct' }, 88888);
    const team = await checkBoxGroupValue({ box: 'nbaTm', slct: 'tm_slct' }, { box: 'wnbaTm', slct: 'wTm_slct' }, 0);
    return { lg, season, team };
}
// SELECTORS
// append options to select
async function makeOption(slct, txt, val) {
    let opt = document.createElement('option');
    opt.textContent = txt;
    opt.value = val;
    opt.style.width = '100%';
    slct.appendChild(opt);
}
// call seaons endpoint for the opts
export async function loadSznOptions() {
    const url = base + '/seasons';
    const r = await fetch(url);
    if (!r.ok)
        throw new Error(`HTTP Error from ${url}`);
    const data = await r.json();
    await buildSznSelects(data);
}
// accept seasons in data object and make an option for each
async function buildSznSelects(data) {
    const rs = document.getElementById('rs_slct');
    const ps = document.getElementById('ps_slct');
    for (let s of data) {
        if (s.season_id.substring(0, 1) === '4') {
            await makeOption(ps, s.season, s.season_id);
        }
        else if (s.season_id.substring(0, 1) === '2') {
            await makeOption(rs, s.season, s.season_id);
        }
    }
}
// call seaons endpoint for the opts
export async function loadTeamOptions() {
    const url = base + '/teams';
    const r = await fetch(url);
    if (!r.ok)
        throw new Error(`HTTP Error from ${url}`);
    const data = await r.json();
    await buildTeamSelects(data);
}
// accept seasons in data object and make an option for each
async function buildTeamSelects(data) {
    const nba = document.getElementById('tm_slct');
    const wnba = document.getElementById('wTm_slct');
    for (let t of data) {
        let txt = `${t.team} | ${t.team_long}`;
        if (t.league === 'NBA') {
            await makeOption(nba, txt, t.team_id);
        }
        else if (t.league === 'WNBA') {
            await makeOption(wnba, txt, t.team_id);
        }
    }
}
//# sourceMappingURL=inputs.js.map