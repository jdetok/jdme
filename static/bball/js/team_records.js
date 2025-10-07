import { tblCaption } from "./table.js"
import { base } from "./listen.js"

export async function getTeamRecords() {
    try {
        const r = await fetch(`${base}/teamrecs`);
        const js = await r.json();
        if (js) {
            return js;
        }

    } catch(err) {
        throw new Error(`HTTP Error (${r.status}) attempting to fetch ${player}
            \n${err}`);
    }
}

export async function buildTeamRecsTbl(js) {
    console.log(js);
}