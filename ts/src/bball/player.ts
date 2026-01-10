import { base, bytes_in_resp, checkBoxEls, scrollIntoBySize, MSG, foldedLog, MSG_BOLD, RED_BOLD } from "../global.js";

export async function fetchPlayer(base: string, player: string, 
    season: string, team: string, lg: string, errEl: string
) { // add season & team
    const errmsg = document.getElementById(errEl); // sErr is elId
    if (!errmsg) throw new Error(`%ccould not find error string element at ${errEl}`)
    if (errmsg.style.display === "block") {
        errmsg.style.display = 'none';
    }

    // encode passed args to be ready for query string
    const s = encodeURIComponent(season)
    const p = encodeURIComponent(player).toLowerCase();

    const req = `${base}/v2/players?player=${p}&season=${s}&team=${team}&league=${lg}`;
    // attempt to fetch from /player endpoint with encoded params
    try {
        const r = await fetch(req);
        if (!r.ok) {
            throw new Error(`HTTP Error (${r.status}) attempting to fetch ${player}`);
        }
        foldedLog(`%c ${await bytes_in_resp(r)} bytes received from ${req}}`, MSG)

        const js = await r.json()
        if (js) {
            if (js.error_string) {
                errmsg.textContent = js.error_string;
                errmsg.style.display = "block";
                return;
            } else {
                return js;
            }
        }
    } catch(err) {
        errmsg.textContent = `can't find ${player}`;
        errmsg.style.display = "block";
        console.error(`an error occured attempting to fetch ${player}\n${err}`);
        
    }
}

