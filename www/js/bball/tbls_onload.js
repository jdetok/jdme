import { Tbl } from "./tbl.js";
import { searchPlayer } from "./player.js";
import { base, fetchJSON, foldedLog, RED_BOLD, SBL, MSG } from "../global.js";
const getRGRow = (d, i) => {
    const player = d.top_scorers[i];
    const game = d.recent_games.find((g) => g.player_id === player.player_id);
    return { player, game };
};
export async function getRGData() {
    return await fetchJSON(`${base}/games/recent`);
}
// top scorers from most recent day of games
// data called on content load and passed through here to have max number of entries for row state
export async function makeRgTopScorersTbl(numRows, data_in) {
    foldedLog(`%cattempting to build RgTopScorers table...`, SBL);
    let data;
    if (!data_in) {
        data = await getRGData();
    }
    else {
        data = data_in;
    }
    if (!data?.top_scorers?.length || !data?.recent_games?.length) {
        foldedLog("%cRGData missing or malformed, skipping table build", RED_BOLD);
        return;
    }
    const rows = Math.min(numRows, data.top_scorers.length);
    new Tbl('tstbl', `Top ${rows} Scorers | ${data.recent_games[0].game_date}`, rows, data, [
        {
            header: 'rank',
            value: (_, i) => String(i + 1),
        }, {
            header: 'name | team',
            value: (d, i) => {
                const { player, game } = getRGRow(d, i);
                return `${player.player} | ${game.team}`;
            },
            button: {
                onClick: async (v) => searchPlayer('button', v.split(" | ")[0]),
            },
        }, {
            header: 'matchup',
            value: (d, i) => {
                const { game } = getRGRow(d, i);
                return game.matchup;
            },
        }, {
            header: 'wl | score',
            value: (d, i) => {
                const { game } = getRGRow(d, i);
                return `${game.wl} | ${game.points}-${game.opp_points}`;
            },
        }, {
            header: 'points',
            value: (d, i) => String(d.top_scorers[i].points),
        },
    ]).init();
    foldedLog(`%cRgTopScorers table built successfully`, MSG);
}
export async function getLGData(numRows) {
    return await fetchJSON(`${base}/league/scoring-leaders?num=${numRows}`);
}
// top scorers in the current season
export async function makeLgTopScorersTbl(numRows) {
    foldedLog(`%cattempting to build LgTopScorers table...`, SBL);
    const data = await getLGData(numRows);
    const rows = Math.min(numRows, data.nba.length);
    new Tbl('nba_tstbl', `Scoring Leaders | NBA/WNBA Top ${rows}`, rows, data, [
        {
            header: "rank",
            value: (_, i) => String(i + 1),
        },
        {
            header: `nba | ${data.nba[0].season}`,
            value: (d, i) => `${d.nba[i].player}`,
            button: {
                onClick: async (v) => searchPlayer('button', v.split(" | ")[0]),
            }
        },
        {
            header: "points",
            value: (d, i) => String(d.nba[i].points),
        },
        {
            header: `wnba | ${data.wnba[0].season}`,
            value: (d, i) => `${d.wnba[i].player} | ${d.wnba[i].team}`,
            button: {
                onClick: async (v) => searchPlayer('button', v.split(" | ")[0]),
            },
        },
        {
            header: "points",
            value: (d, i) => String(d.wnba[i].points),
        },
    ]).init();
    foldedLog(`%cLgTopScorers table built successfully`, MSG);
}
export async function getTRData() {
    return await fetchJSON(`${base}/teamrecs`);
}
export async function makeTeamRecordsTbl(numRows) {
    foldedLog(`%cattempting to build TeamRecs table...`, SBL);
    const data = await getTRData();
    const rows = Math.min(numRows, data.nba_team_records.length);
    new Tbl("trtbl", "NBA/WNBA Regular Season Team Records", rows, data, [
        { header: "rank", value: (_, i) => String(i + 1) },
        {
            header: `nba | ${data.nba_team_records[0].season}`,
            value: (d, i) => d.nba_team_records[i].team_long,
        },
        {
            header: "record",
            value: (d, i) => `${d.nba_team_records[i].wins}-${d.nba_team_records[i].losses}`,
        },
        {
            header: `wnba | ${data.wnba_team_records[0].season ?? '-'}`,
            value: (d, i) => {
                if (i >= d.wnba_team_records.length)
                    return '-';
                return d.wnba_team_records[i].team_long ?? '-';
            },
        },
        {
            header: "record",
            value: (d, i) => {
                if (i >= d.wnba_team_records.length)
                    return '-';
                return `${d.wnba_team_records[i].wins ?? ''}-${d.wnba_team_records[i].losses ?? ''}`;
            },
        },
    ]).init();
    foldedLog(`%cTeamRecs table built successfully`, MSG);
}
//# sourceMappingURL=tbls_onload.js.map