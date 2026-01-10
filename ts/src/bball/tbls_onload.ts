import { Tbl } from "./tbl.js"
import { base } from "../global.js";
import { playerBtnListener } from "./listeners.js";

export async function makeLgTopScorersTbl(numPl: number): Promise<void> {
    new Tbl(
        'nba_tstbl', 
        `Scoring Leaders | NBA/WNBA Top ${numPl}`,
        numPl, `${base}/league/scoring-leaders?num=${numPl}`, [
            {
                header: "rank", 
                value: (_, i) => String(i + 1),
            }, 
            {
                header: `nba | ${data.nba[0].season}`, 
                value: (d, i) => `${d.nba[i].player}`,
                button: {
                    onClick: async (v) => playerBtnListener(v.split(" | ")[0]),
                }
            }
        ]
    ).init();
}   