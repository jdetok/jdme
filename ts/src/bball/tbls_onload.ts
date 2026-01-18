
import { Tbl } from "./tbl.js";
import { searchPlayer } from "./player.js";
import { base, fetchJSON, foldedLog, RED_BOLD, SBL, MSG, foldedErr } from "../global.js";

const getRGRow = (d: any, i: number) => {
    const player = d.top_scorers[i];
    const game = d.recent_games.find(
        (g: any) => g.player_id === player.player_id
    );
    return { player, game };
};

export type recentGameTopScorer = {
    player_id: number,
    team_id: number,
    player: string,
    league: "NBA" | "WNBA",
    points: number, 
    assists: number, 
    rebounds: number,
};

export type recentGame = {
    game_id: number,
    team_id: number,
    player_id: number,
    player: string,
    league: "NBA" | "WNBA",
    team: string,
    team_name: string,
    game_date: string,
    matchup: string,
    wl: string,
    points: number,
    opp_points: number,
};

export type RGData = {
    top_scorers: recentGameTopScorer[],
    recent_games: recentGame[],
}

export async function getRGData(): Promise<RGData> {
    return await fetchJSON(`${base}/games/recent`);
}

// top scorers from most recent day of games
// data called on content load and passed through here to have max number of entries for row state
export async function makeRgTopScorersTbl(numRows: number, data_in?: RGData): Promise<void> {
    foldedLog(`%cattempting to build RgTopScorers table...`, SBL);
    let data: RGData;
    try {
        if (!data_in) {
            data = await getRGData();
        } else {
            data = data_in;
        }
    } catch (e) {
        throw new Error(`error fetching RG data: ${e}`);
    }

    if (!data?.top_scorers?.length || !data?.recent_games?.length) {
        throw new Error("RGData missing or malformed, skipping table build");
    }

    const rows = Math.min(
        numRows,
        data.top_scorers.length,
    );
    try {
        new Tbl(
            'tstbl',
            `Top ${rows} Scorers | ${data.recent_games[0].game_date}`,
            rows, data, [
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
        ]
        ).init();
    } catch (e) {
        throw new Error(`error building RgTopScorers table: ${e}`);
    }
    foldedLog(`%cRgTopScorers table built successfully`, MSG);
}

export type scoringLeader = {
    player_id: number,
    player: string,
    season: string,
    team: string,
    point: number,
};

export type LGData = {
    nba: scoringLeader[];
    wnba: scoringLeader[];
}

export async function getLGData(numRows: number): Promise<LGData> {
    return await fetchJSON(`${base}/league/scoring-leaders?num=${numRows}`);
}

// top scorers in the current season
export async function makeLgTopScorersTbl(numRows: number): Promise<void> {
    foldedLog(`%cattempting to build LgTopScorers table...`, SBL);

    let data: LGData;
    try {
        data = await getLGData(numRows);
    } catch (e) {
        throw new Error(`error fetching LG data: ${e}`);
    }

    const rows = Math.min(numRows, data.nba.length);

    try {
        new Tbl(
            'nba_tstbl', 
            `Scoring Leaders | NBA/WNBA Top ${rows}`,
            rows, data, [
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
                    value: (d: any, i) => `${d.wnba[i].player} | ${d.wnba[i].team}`,
                    button: {
                        onClick: async (v) => searchPlayer('button',v.split(" | ")[0]),
                    },
                },
                {
                    header: "points",
                    value: (d, i) => String(d.wnba[i].points),
                },
            ]
        ).init();
    } catch (e) {
        throw new Error(`error building LgTopScorers table: ${e}`);
    }
    foldedLog(`%cLgTopScorers table built successfully`, MSG);
}

export type TeamRec = {
    league: "NBA" | "WNBA",
    season_id: number,
    season: string,
    season_desc: string,
    team_id: number,
    team: string,
    team_long: string,
    wins: number,
    losses: number,
};

export type TRData = {
    nba_team_records: TeamRec[];
    wnba_team_records: TeamRec[];
};

export async function getTRData(): Promise<TRData> {
    return await fetchJSON(`${base}/teamrecs`);
}


export async function makeTeamRecordsTbl(numRows: number): Promise<void> {
    foldedLog(`%cattempting to build TeamRecs table...`, SBL);
    
    let data: TRData;
    try {
        data = await getTRData();
    } catch (e) {
        throw new Error(`error fetching TR data: ${e}`);
    }

    const rows = Math.min(numRows, data.nba_team_records.length);
    
    try {
        new Tbl(
            "trtbl",
            "NBA/WNBA Regular Season Team Records", rows, data,
            [
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
                        if (i >= d.wnba_team_records.length) return '-';
                        return d.wnba_team_records[i].team_long ?? '-';
                    },
                },
                {
                    header: "record",
                    value: (d, i) => {
                        if (i >= d.wnba_team_records.length) return '-';
                        return `${d.wnba_team_records[i].wins ?? ''}-${d.wnba_team_records[i].losses ?? ''}`;
                    },
                },
            ],
        ).init();
    } catch (e) {
        throw new Error(`error building RgTopScorers table: ${e}`);
    }
    foldedLog(`%cTeamRecs table built successfully`, MSG);
}