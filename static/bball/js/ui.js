import { showHideHvr } from "./hover.js";
import { base } from "./listen.js";

// return season or 88888 if no sseason box checked
export async function handleSeasonBoxes() {
    const p = await checkBoxes('post', 'ps_slct')
    const r = await checkBoxes('reg', 'rs_slct')

    if (p) {
        return p;
    }
    if (r) {
        return r;
    } // 2{current year} i.e. 22025 during 2025
    return 88888;
    // return `2${new Date().getFullYear()}`;
}

export async function clearCheckBoxes(box) {
    const b = document.getElementById(box);
    b.checked = 0;
}

// if checkbox is checked, return the value
export async function checkBoxes(box, sel) {
    const b = document.getElementById(box);
    const s = document.getElementById(sel);
    if (b.checked) {
        return s.value
    }
}
// make post + reg checkboxes exclusive (but allow neither checked)
export async function setupExclusiveCheckboxes() {
    const post = document.getElementById("post");
    const reg = document.getElementById("reg");

    function handleCheck(e) {
        if (e.target.checked) {
            if (e.target === post) reg.checked = false;
            if (e.target === reg) post.checked = false;
        }
    }

    post.addEventListener("change", handleCheck);
    reg.addEventListener("change", handleCheck);
}


// season select hover messages
export async function selHvr() {
    const rs = document.getElementById('hlpRs');
    const ps = document.getElementById('hlpPs');
    await showHideHvr(rs, 'selhvr',
        `search for a specific regular-season. if the player being searched didn't
        play in the selected season, their first or most recent season, whichever
        is closer to the selected, will be used`);
    await showHideHvr(ps, 'selhvr',
        `search for a specific post-season. if the player being searched didn't
        play in the selected season, their first or most recent season, whichever
        is closer to the selected, will be used`);
}

// append options to select
async function makeOption(slct, txt, val) {
    let opt = document.createElement('option');
    opt.textContent = txt;
    opt.value = val;
    slct.appendChild(opt);
}

// call seaons endpoint for the opts
export async function loadSznOptions() {
    const r = await fetch(base + '/seasons');
    if (!r.ok) { 
        throw new Error(`HTTP Error: ${s.status}`);
    } 
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
        } else if (s.season_id.substring(0, 1) === '2') {
            await makeOption(rs, s.season, s.season_id);
        }
    }
}