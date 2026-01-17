import { buildOnLoadDash } from "./player.js";
import { foldedLog, MSG_BOLD, RED_BOLD } from "../global.js";
import { clearSearch, lgRadioBtns, loadSznOptions, loadTeamOptions,  } from "./inputs.js";
import { makeLgTopScorersTbl, makeRgTopScorersTbl, makeTeamRecordsTbl, getRGData } from "./tbls_onload.js";
import { rowsState, makeExpandTblBtns } from "./rowstate.js"
import { submitPlayerSearch, randPlayerBtn, holdPlayerBtn, setup_jump_btns, setupExclusiveCheckboxes, clearSearchBtn } from "./listeners.js";
import { makeLogoImgs } from "./img.js";

// CALL ENTRYPOINT
await LoadContent();

// ENTRYPOINT DEFINITION
async function LoadContent(): Promise<void> {
    // create state class to track number of rows displayed per table

    document.addEventListener('DOMContentLoaded', async () => {

        const recent_game_data = await getRGData();

        let ROWSTATE = new rowsState(recent_game_data); 
        foldedLog(`%cloading content for page {${window.innerWidth}px x ${window.innerHeight}px}...`, MSG_BOLD)
        
        // setup default buttons / inputs
        clearSearch();
        await lgRadioBtns();
        await setup_jump_btns();
        await makeExpandTblBtns(ROWSTATE);
        await setupExclusiveCheckboxes('post', 'reg');
        await setupExclusiveCheckboxes('nbaTm', 'wnbaTm');
        await loadSznOptions();
        await loadTeamOptions();
        await makeLogoImgs();

        // build tables and recent top scorer dash on initial load
        try {
            await makeLgTopScorersTbl(ROWSTATE.lgRowNum.value);
            await makeRgTopScorersTbl(ROWSTATE.rgRowNum.value, recent_game_data);
            await makeTeamRecordsTbl(ROWSTATE.startRows);
            await buildOnLoadDash();
        } catch (err) {foldedLog(`%cerror building on load elements: ${err}`, RED_BOLD)}

        // listen for submissions
        await clearSearchBtn();
        await submitPlayerSearch();
        await randPlayerBtn();
        await holdPlayerBtn(); 
    });
}
