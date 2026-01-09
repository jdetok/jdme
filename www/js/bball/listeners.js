import { clearSearch } from "./btns.js";
import { foldedLog, MSG_BOLD } from "../global.js";
let NUMPL = window.innerWidth <= 700 ? 5 : 10;
export async function LoadContent() {
    document.addEventListener('DOMContentLoaded', async () => {
        foldedLog(`%cloading content for page {${window.innerWidth}px x ${window.innerHeight}px}...`, MSG_BOLD);
        // await buildOnLoadElements();
        // await searchPlayer();
        // await ui.randPlayerBtn();
        clearSearch();
        // await ui.holdPlayerBtn();
    });
}
//# sourceMappingURL=listeners.js.map