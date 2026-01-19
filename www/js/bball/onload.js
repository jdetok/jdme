// import { RED_BOLD, foldedLog } from "../global.js";
import { NBA_WNBA_LOGO_IMGS, fillImageDiv, normalizeImgHeights } from "./elements.js";
import { base, fetchJSON, foldedLog, SBL, MSG, BIGWINDOW, LARGEROWS } from "../global.js";
import { Tbl } from "./tbl.js";
import { fetchAndBuildPlayerDash } from "./player_dash.js";
import { clearSearch, lgRadioBtns, loadSznOptions, loadTeamOptions, } from "./inputs.js";
import { listenForInput, setup_jump_btns, setupExclusiveCheckboxes } from "./listeners.js";
let exBtnsInitComplete = false;
export async function initUIElements(rs) {
    try {
        clearSearch();
        await lgRadioBtns();
        await setup_jump_btns();
        await makeExpandTblBtns(rs);
        await setupExclusiveCheckboxes('post', 'reg');
        await setupExclusiveCheckboxes('nbaTm', 'wnbaTm');
        await loadSznOptions();
        await loadTeamOptions();
    }
    catch (e) {
        throw new Error(`error occured building UI elements: ${e}`);
    }
}
export async function initEventListeners() {
    try {
        return await listenForInput();
    }
    catch (e) {
        throw new Error(`error occured setting up event listeners: ${e}`);
    }
}
export async function buildOnLoadContent(rs, rg) {
    try {
        await makeLogoImgs();
        await makeLgTopScorersTbl(rs.lgRowNum.value);
        await makeRgTopScorersTbl(rs.rgRowNum.value, rg);
        await makeTeamRecordsTbl(rs.startRows);
        await buildOnLoadDash(rg);
    }
    catch (e) {
        throw new Error(`error occured setting up event listeners: ${e}`);
    }
}
async function makeLogoImgs() {
    const imgs = await fillImageDiv(NBA_WNBA_LOGO_IMGS);
    await normalizeImgHeights(imgs);
}
export async function rebuildContent(tr_rows, lg_rows, rg_rows, rgData) {
    return await Promise.all([
        makeLogoImgs(),
        makeTeamRecordsTbl(tr_rows),
        makeLgTopScorersTbl(lg_rows),
        makeRgTopScorersTbl(rg_rows, rgData),
    ]);
}
async function buildOnLoadDash(rgData) {
    try {
        await fetchAndBuildPlayerDash('onload', null, rgData);
    }
    catch (e) {
        throw new Error(`error building onload dash: ${e}`);
    }
}
export async function makeExpandTblBtns(rs, tblBtns = [
    { elId: "seemorelessLGbtns", rows: rs.lgRowNum, build: makeLgTopScorersTbl },
    { elId: "seemorelessRGbtns", rows: rs.rgRowNum, build: makeRgTopScorersTbl },
    { elId: "seemorelessTRbtns", rows: rs.trRowNum, build: makeTeamRecordsTbl },
]) {
    if (exBtnsInitComplete)
        return;
    exBtnsInitComplete = true;
    for (let etb of tblBtns) {
        const d = document.getElementById(etb.elId);
        if (!d)
            continue;
        let to_append = [];
        for (const obj of [
            { op: 'all', lbl: 'see all' },
            { op: '+', lbl: 'see more' },
            { op: '-', lbl: 'see less' },
            { op: 'rst', lbl: 'reset' },
            { op: 'min', lbl: 'minimize' }
        ]) {
            let newNum;
            const btn = document.createElement('button');
            btn.textContent = obj.lbl;
            btn.addEventListener('click', async () => {
                switch (obj.op) {
                    case 'all':
                        newNum = etb.rows.max();
                        break;
                    case 'min':
                        newNum = etb.rows.min();
                        break;
                    case '+':
                        newNum = etb.rows.increase();
                        break;
                    case '-':
                        newNum = etb.rows.decrease();
                        break;
                    case 'rst':
                        if (etb.elId === 'seemorelessLGbtns' && window.innerWidth >= BIGWINDOW) {
                            newNum = etb.rows.reset(LARGEROWS);
                        }
                        else {
                            newNum = etb.rows.reset();
                        }
                        break;
                    default:
                        throw new Error(`invalid case: ${obj.op} | ${obj.lbl}`);
                }
                await etb.build(newNum);
            });
            to_append.push(btn);
        }
        for (const b of to_append) {
            d.appendChild(b);
        }
    }
}
;
const getRGRow = (d, i) => {
    const player = d.top_scorers[i];
    const game = d.recent_games.find((g) => g.player_id === player.player_id);
    return { player, game };
};
export async function getRGData() {
    return await fetchJSON(`${base}/games/recent`);
}
// top scorers from most recent day of games
// data called on content load and passed through here to have max number of entries for row state
async function makeRgTopScorersTbl(numRows, data_in) {
    foldedLog(`%cattempting to build RgTopScorers table...`, SBL);
    let data;
    try {
        if (!data_in) {
            data = await getRGData();
        }
        else {
            data = data_in;
        }
    }
    catch (e) {
        throw new Error(`error fetching RG data: ${e}`);
    }
    if (!data?.top_scorers?.length || !data?.recent_games?.length) {
        throw new Error("RGData missing or malformed, skipping table build");
    }
    const rows = Math.min(numRows, data.top_scorers.length);
    try {
        new Tbl('tstbl', `Top ${rows} Scorers | ${data.recent_games[0].game_date}`, rows, data, [
            {
                header: 'rank',
                value: (_, i) => String(i + 1),
            }, {
                header: 'name | team',
                value: (d, i) => {
                    const { player, game } = getRGRow(d, i);
                    return `${player.player} | ${game.team}`;
                },
                button: {
                    onClick: async (v) => fetchAndBuildPlayerDash('button', v.split(" | ")[0]),
                },
            }, {
                header: 'matchup',
                value: (d, i) => {
                    const { game } = getRGRow(d, i);
                    return game.matchup;
                },
            }, {
                header: 'wl | score',
                value: (d, i) => {
                    const { game } = getRGRow(d, i);
                    return `${game.wl} | ${game.points}-${game.opp_points}`;
                },
            }, {
                header: 'points',
                value: (d, i) => String(d.top_scorers[i].points),
            },
        ]).init();
    }
    catch (e) {
        throw new Error(`error building RgTopScorers table: ${e}`);
    }
    foldedLog(`%cRgTopScorers table built successfully`, MSG);
}
async function getLGData(numRows) {
    return await fetchJSON(`${base}/league/scoring-leaders?num=${numRows}`);
}
// top scorers in the current season
async function makeLgTopScorersTbl(numRows) {
    foldedLog(`%cattempting to build LgTopScorers table...`, SBL);
    let data;
    try {
        data = await getLGData(numRows);
    }
    catch (e) {
        throw new Error(`error fetching LG data: ${e}`);
    }
    const rows = Math.min(numRows, data.nba.length);
    try {
        new Tbl('nba_tstbl', `Scoring Leaders | NBA/WNBA Top ${rows}`, rows, data, [
            {
                header: "rank",
                value: (_, i) => String(i + 1),
            },
            {
                header: `nba | ${data.nba[0].season}`,
                value: (d, i) => `${d.nba[i].player}`,
                button: {
                    onClick: async (v) => fetchAndBuildPlayerDash('button', v.split(" | ")[0]),
                }
            },
            {
                header: "points",
                value: (d, i) => String(d.nba[i].points),
            },
            {
                header: `wnba | ${data.wnba[0].season}`,
                value: (d, i) => `${d.wnba[i].player} | ${d.wnba[i].team}`,
                button: {
                    onClick: async (v) => fetchAndBuildPlayerDash('button', v.split(" | ")[0]),
                },
            },
            {
                header: "points",
                value: (d, i) => String(d.wnba[i].points),
            },
        ]).init();
    }
    catch (e) {
        throw new Error(`error building LgTopScorers table: ${e}`);
    }
    foldedLog(`%cLgTopScorers table built successfully`, MSG);
}
async function getTRData() {
    return await fetchJSON(`${base}/teamrecs`);
}
async function makeTeamRecordsTbl(numRows) {
    foldedLog(`%cattempting to build TeamRecs table...`, SBL);
    let data;
    try {
        data = await getTRData();
    }
    catch (e) {
        throw new Error(`error fetching TR data: ${e}`);
    }
    const rows = Math.min(numRows, data.nba_team_records.length);
    try {
        new Tbl("trtbl", "NBA/WNBA Regular Season Team Records", rows, data, [
            { header: "rank", value: (_, i) => String(i + 1) },
            {
                header: `nba | ${data.nba_team_records[0].season}`,
                value: (d, i) => d.nba_team_records[i].team_long,
            },
            {
                header: "record",
                value: (d, i) => `${d.nba_team_records[i].wins}-${d.nba_team_records[i].losses}`,
            },
            {
                header: `wnba | ${data.wnba_team_records[0].season ?? '-'}`,
                value: (d, i) => {
                    if (i >= d.wnba_team_records.length)
                        return '-';
                    return d.wnba_team_records[i].team_long ?? '-';
                },
            },
            {
                header: "record",
                value: (d, i) => {
                    if (i >= d.wnba_team_records.length)
                        return '-';
                    return `${d.wnba_team_records[i].wins ?? ''}-${d.wnba_team_records[i].losses ?? ''}`;
                },
            },
        ]).init();
    }
    catch (e) {
        throw new Error(`error building RgTopScorers table: ${e}`);
    }
    foldedLog(`%cTeamRecs table built successfully`, MSG);
}
//# sourceMappingURL=onload.js.map