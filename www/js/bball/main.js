import { buildOnLoadDash } from "./player.js";
import { foldedLog, foldedErr, MSG_BOLD, SBL } from "../global.js";
import { clearSearch, lgRadioBtns, loadSznOptions, loadTeamOptions, } from "./inputs.js";
import { makeLgTopScorersTbl, makeRgTopScorersTbl, makeTeamRecordsTbl, getRGData } from "./tbls_onload.js";
import { rowsState, makeExpandTblBtns } from "./rowstate.js";
import { submitPlayerSearch, randPlayerBtn, holdPlayerBtn, setup_jump_btns, setupExclusiveCheckboxes, clearSearchBtn } from "./listeners.js";
import { makeLogoImgs } from "./img.js";
// CALL ENTRYPOINT
await LoadContent();
// ENTRYPOINT DEFINITION
async function LoadContent() {
    foldedLog(`%cbuilding UI once DOM content loads...`, SBL);
    document.addEventListener('DOMContentLoaded', async () => {
        const wsize = `W:${window.innerWidth}px X H:${window.innerHeight}px`;
        foldedLog(`%cDOM loaded for ${wsize} page... `, SBL);
        let recentGameData;
        let tblRowState;
        let gameDate;
        try {
            recentGameData = await getRGData();
            gameDate = recentGameData.recent_games[0].game_date;
            tblRowState = new rowsState(recentGameData);
        }
        catch (e) {
            foldedErr(`error getting recent game data: ${e}`);
            return;
        }
        foldedLog(`%csetting up UI elements...`, SBL);
        try {
            clearSearch();
            await lgRadioBtns();
            await setup_jump_btns();
            await makeExpandTblBtns(tblRowState);
            await setupExclusiveCheckboxes('post', 'reg');
            await setupExclusiveCheckboxes('nbaTm', 'wnbaTm');
            await loadSznOptions();
            await loadTeamOptions();
            await makeLogoImgs();
        }
        catch (e) {
            foldedErr(`error setting up page elements: ${e}`);
        }
        foldedLog(`%cbuilding tables with games data through ${gameDate}...`, SBL);
        try {
            await makeLgTopScorersTbl(tblRowState.lgRowNum.value);
            await makeRgTopScorersTbl(tblRowState.rgRowNum.value, recentGameData);
            await makeTeamRecordsTbl(tblRowState.startRows);
            await buildOnLoadDash(recentGameData);
        }
        catch (e) {
            foldedErr(`error building on load elements: ${e}`);
        }
        foldedLog(`%csetting up button listeners...`, SBL);
        try {
            await clearSearchBtn();
            await submitPlayerSearch();
            await randPlayerBtn();
            await holdPlayerBtn();
        }
        catch (e) {
            foldedErr(`error starting submission listeners: ${e}`);
        }
        foldedLog(`%ccontent loaded succesfully`, MSG_BOLD);
    });
}
//# sourceMappingURL=main.js.map