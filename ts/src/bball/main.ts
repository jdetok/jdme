import { getRGData, initEventListeners, buildOnLoadContent, initUIElements } from "./onload.js";
import { foldedLog, foldedErr, MSG_BOLD, SBL, wsize } from "../global.js";
import { RGData } from "./resp_types.js";
import { rowsState } from "./rowstate.js"

// CALL ENTRYPOINT
await LoadContent();

// ENTRYPOINT DEFINITION
async function LoadContent(): Promise<void> {
    foldedLog(`%cbuilding UI once DOM content loads...`, SBL);
    document.addEventListener('DOMContentLoaded', async () => {
        foldedLog(`%cDOM loaded for ${wsize()} page... `, SBL);
        
        let recentGameData: RGData;
        let tblRowState: rowsState;
        let gameDate: string;

        try {
            recentGameData = await getRGData();
            gameDate = recentGameData.recent_games[0].game_date;
            tblRowState = new rowsState(recentGameData);
        } catch (e) {
            foldedErr(`error getting recent game data: ${e}`);
            return;
        }
        
        try {
            await initUIElements(tblRowState);
        } catch (e) {
            foldedErr(`error setting up page elements: ${e}`);
        }
        
        try {
            await buildOnLoadContent(tblRowState, recentGameData);
        } catch (e) {
            foldedErr(`error building on load elements: ${e}`);
        }

        try {
            await initEventListeners();
        } catch (e) {
            foldedErr(`error starting submission listeners: ${e}`);
        }

        foldedLog(`%ccontent loaded succesfully`, MSG_BOLD);
    });
}
