import { foldedLog, MSG_BOLD, RED_BOLD } from "../global.js";
import { makeLgTopScorersTbl, makeRgTopScorersTbl, makeTeamRecordsTbl } from "./tbl.js";
import { buildLoadDash } from "./player.js";
import { submitPlayerSearch, randPlayerBtn } from "./listeners.js";
import { rowsState, setupExclusiveCheckboxes, clearSearchBtn, clearSearch, lgRadioBtns, holdPlayerBtn, setup_jump_btns, loadSznOptions, loadTeamOptions, expandedListBtns } from "./inputs.js";
// CALL ENTRYPOINT
await LoadContent();
// ENTRYPOINT DEFINITION
export async function LoadContent() {
    // create state class to track number of rows displayed per table
    let ROWSTATE = new rowsState();
    document.addEventListener('DOMContentLoaded', async () => {
        foldedLog(`%cloading content for page {${window.innerWidth}px x ${window.innerHeight}px}...`, MSG_BOLD);
        clearSearch();
        await lgRadioBtns();
        await setup_jump_btns();
        await expandedListBtns(ROWSTATE);
        await setupExclusiveCheckboxes('post', 'reg');
        await setupExclusiveCheckboxes('nbaTm', 'wnbaTm');
        await loadSznOptions();
        await loadTeamOptions();
        await buildOnLoadElements(ROWSTATE);
        await clearSearchBtn();
        await submitPlayerSearch();
        await randPlayerBtn();
        await holdPlayerBtn();
    });
}
// create top numrows tables, load player dash for top scrorer from most recent day of games
export async function buildOnLoadElements(rs) {
    try {
        await makeLgTopScorersTbl(rs.lgRowNum.value);
        await makeRgTopScorersTbl(rs.rgRowNum.value);
        await makeTeamRecordsTbl(rs.startRows);
        await buildLoadDash();
    }
    catch (err) {
        foldedLog(`%cerror building on load elements: ${err}`, RED_BOLD);
    }
}
//# sourceMappingURL=main.js.map