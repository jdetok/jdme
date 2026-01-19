import { getRGData, initEventListeners, buildOnLoadContent, initUIElements } from "./onload.js";
import { foldedLog, foldedErr, MSG_BOLD, SBL } from "../global.js";
import { rowsState } from "./rowstate.js";
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
            await initUIElements(tblRowState);
        }
        catch (e) {
            foldedErr(`error setting up page elements: ${e}`);
        }
        foldedLog(`%cbuilding tables with games data through ${gameDate}...`, SBL);
        try {
            await buildOnLoadContent(tblRowState, recentGameData);
        }
        catch (e) {
            foldedErr(`error building on load elements: ${e}`);
        }
        foldedLog(`%csetting up button listeners...`, SBL);
        try {
            await initEventListeners();
        }
        catch (e) {
            foldedErr(`error starting submission listeners: ${e}`);
        }
        foldedLog(`%ccontent loaded succesfully`, MSG_BOLD);
    });
}
//# sourceMappingURL=main.js.map