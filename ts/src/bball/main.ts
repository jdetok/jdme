import { buildOnLoadDash } from "./player.js";
import { foldedLog, MSG_BOLD, RED_BOLD } from "../global.js";
import { clearSearch, lgRadioBtns, loadSznOptions, loadTeamOptions,  } from "./inputs.js";
import { makeLgTopScorersTbl, makeRgTopScorersTbl, makeTeamRecordsTbl } from "./tbls_onload.js";
import { submitPlayerSearch, randPlayerBtn, holdPlayerBtn, setup_jump_btns, setupExclusiveCheckboxes, clearSearchBtn, expandedListBtns, rowsState} from "./listeners.js";

// CALL ENTRYPOINT
await LoadContent();

// ENTRYPOINT DEFINITION
async function LoadContent(): Promise<void> {
    // create state class to track number of rows displayed per table
    let ROWSTATE = new rowsState();

    document.addEventListener('DOMContentLoaded', async () => {
        foldedLog(`%cloading content for page {${window.innerWidth}px x ${window.innerHeight}px}...`, MSG_BOLD)
        
        // setup default buttons / inputs
        clearSearch();
        await lgRadioBtns();
        await setup_jump_btns();
        await expandedListBtns(ROWSTATE);
        await setupExclusiveCheckboxes('post', 'reg');
        await setupExclusiveCheckboxes('nbaTm', 'wnbaTm');
        await loadSznOptions();
        await loadTeamOptions();

        // build tables and recent top scorer dash on initial load
        try {
            await makeLgTopScorersTbl(ROWSTATE.lgRowNum.value);
            await makeRgTopScorersTbl(ROWSTATE.rgRowNum.value);
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
