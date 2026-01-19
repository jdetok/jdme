import { base, MSG, foldedLog, MSG_BOLD, fetchJSON } from "../global.js";
import { getInputVals, inputVals } from "./inputs.js";
import { RGData, PlayersResp } from "./resp_types.js";

export type PlayerSearchType = 'onload' | 'random' | 'submit' | 'button';

export async function searchPlayer(pst: PlayerSearchType = 'submit',
    playerOverride: string | null, rgData?: RGData, searchEl: string = 'pSearch', holdEl: string = 'pHold'
): Promise<PlayersResp> {
    const input = document.getElementById(searchEl) as HTMLInputElement;
    if (!input) throw new Error(`couldn't get element at Id ${searchEl}`);

    let iv: inputVals;
    try {
        iv = await getInputVals();
    } catch (e) {
        throw new Error(`error getting input values`);
    }
    
    let player = '';

    try {
        player = handlePlayerSearchType(pst, input, iv, playerOverride, rgData);
        if (player === '') {
            const pHoldEl = document.getElementById(holdEl) as HTMLInputElement;
            if (!pHoldEl) throw new Error(`could not find input element with id ${holdEl}`);
            player = pHoldEl.value;
        }
    } catch (e) {
        throw new Error(`error handling player search type: ${pst}`);
    }

    foldedLog(`%csearching for player ${player} | season ${iv.season} | team ${iv.team} | league ${iv.lg}`, MSG_BOLD);

    try {
        const data = await fetchPlayer(base, player, iv.season, iv.team, iv.lg);
        if (!data) throw new Error(`failed to get data for player ${player} | season ${iv.season} | team ${iv.team}`);
        return data;
    } catch (e) {
        throw new Error(`error fetching data for player ${player} | season ${iv.season} | team ${iv.team}: ${e}`);
    }
}

function handlePlayerSearchType(pst: PlayerSearchType, input: HTMLInputElement, iv: inputVals,
    playerOverride?: string | null, rgData?: RGData
): string {
        let player = '';
        switch (pst) {
        case 'onload':
            if (!rgData) throw new Error(`%cmust pass recent game data for onload mode`) 
            player = rgData.top_scorers[0].player;
            foldedLog(`%cfetched onload player ${player} | league: ${iv.lg} | season ${iv.season}`, MSG);
            break;
        case 'submit': 
            player = input.value.trim();
            break;
        case 'random': 
            foldedLog(`%csearching random player | league: ${iv.lg} | season ${iv.season}`, MSG);
            player = String(pst);
            break;
        case 'button':
            foldedLog(`%cfetching player ${player} from button | league: ${iv.lg} | season ${iv.season}`, MSG);
            if (!playerOverride) throw new Error(`playerOverride must be passed if called with pst='button'`);
            player = playerOverride;
            break;
        default: 
            throw new Error(`option passed as pst "${pst}" not valid`);
        }
        return player;
    }


async function fetchPlayer(base: string, player: string | number, 
    season: string | number, team: string | number, lg: string, errEl = 'sErr'
): Promise<PlayersResp> {
    const errmsg = document.getElementById(errEl);
    if (!errmsg) throw new Error(`%ccould not find error string element at ${errEl}`);
    if (errmsg.style.display === "block") {
        errmsg.style.display = 'none';
    }

    const s = encodeURIComponent(season);
    const p = encodeURIComponent(player).toLowerCase();
    const url = `${base}/v2/players?player=${p}&season=${s}&team=${team}&league=${lg}`;

    return await fetchJSON(url);
}