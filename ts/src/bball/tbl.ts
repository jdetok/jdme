// replaces /www/js/bball/dynamic_table.js
// dynamic replacement for lg_ldg_scorers.js, rg_ldg_scorer.js, reamrecs.js

import { playerBtnListener } from "./listeners.js";
import { base, foldedLog, MSG } from "../global.js";

type CellValue<T> = (data: T | any, idx: number) => string;

type Column<T> = {
    header: string;
    value: CellValue<T>;
    button?: {
        onClick: (value: string, data: T, idx: number) => Promise<void>;
    };
};

export class Tbl<T> {
    private data!: T;

    constructor(
        private elId: string,
        private title: string,
        private rowCount: number,
        private url: string,
        private columns: Column<T>[],
    ) {}

    async init(): Promise<void> {
        const r = await fetch(this.url);
        this.data = await r.json();
        this.build();
    }

    private build(): void {
        let tbl = document.getElementById(this.elId) as HTMLTableElement;
        tbl.innerHTML = '';

        tbl.appendChild(this.makeTitle());

        tbl.appendChild(this.makeHdrRow());

        for (let i = 0; i < this.rowCount; i++) {
            tbl.appendChild(this.makeRow(i));
        }
    }

    makeTitle(): HTMLTableCaptionElement {
    let capt = document.createElement('caption') as HTMLTableCaptionElement;
        capt.textContent = this.title;
        return capt;
    }

    makeHdrRow(): HTMLTableSectionElement {
        const thead = document.createElement('thead');
        const tr = document.createElement('tr');
        for (const col of this.columns) {
            const td = document.createElement('td');
            td.textContent = col.header;
            tr.appendChild(td);
        }
        thead.appendChild(tr);
        return thead;
    }

    makeRow(idx: number): HTMLTableRowElement {
        const tr = document.createElement("tr");

        this.columns.forEach(col => {
            const td = document.createElement('td');
            const val = col.value(this.data, idx);
            if (col.button) {
                const btn = document.createElement("button");
                btn.type = "button";
                btn.textContent = val;
                btn.onclick = () => col.button!.onClick(val, this.data, idx);
                td.appendChild(btn);
            } else {
                td.textContent = val;
            }

            tr.appendChild(td);
        })
        return tr;
    }
}

const getRGRow = (d: any, i: number) => {
    const player = d.top_scorers[i];
    const game = d.recent_games.find(
        (g: any) => g.player_id === player.player_id
    );
    return { player, game };
};

// top scorers from most recent day of games
export async function makeRgTopScorersTbl(numRows: number): Promise<void> {
    let datasrc = `${base}/games/recent`;
    let r = await fetch(datasrc);
    const data = await r.json();
    foldedLog(`attempting to build RgTopScorers table...`);
    console.log(data);
    new Tbl(
        'tstbl', 
        `Top ${Math.min(numRows, data.top_scorers.length)} Scorers | ${data.recent_games[0].game_date}`,
        numRows, datasrc, [
            {
                header: 'rank',
                value: (_, i) => String(i + 1),
            },
            {
                header: 'name | team',
                value: (d, i) => {
                    const { player, game } = getRGRow(d, i);
                    return `${player.player} | ${game.team}`;

                },
                button: {
                    onClick: async (v) => playerBtnListener(v.split(" | ")[0]),
                },
            },
            {
                header: 'matchup',
                value: (d, i) => {
                    const { game } = getRGRow(d, i);
                    return game.matchup;
                },
            },
            {
                header: 'wl | score',
                value: (d, i) => {
                    const { game } = getRGRow(d, i);
                    return `${game.wl} | ${game.points}-${game.opp_points}`;
                },
                button: {
                    onClick: async (v) => playerBtnListener(v.split(" | ")[0]),
                },
            },
            {
                header: 'points',
                value: (d, i) => String(d.top_scorers[i].points),
            },
        ]
    ).init();
}

// top scorers in the current season
export async function makeLgTopScorersTbl(numRows: number): Promise<void> {
    let datasrc = `${base}/league/scoring-leaders?num=${numRows}`;
    let r = await fetch(datasrc);
    const data = await r.json();
    new Tbl(
        'nba_tstbl', 
        `Scoring Leaders | NBA/WNBA Top ${numRows}`,
        numRows, datasrc, [
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
            },
            {
                header: "points",
                value: (d, i) => String(d.nba[i].points),
            },
            {
                header: `wnba | ${data.wnba[0].season}`,
                value: (d: any, i) => `${d.wnba[i].player} | ${d.wnba[i].team}`,
                button: {
                    onClick: async (v) => playerBtnListener(v.split(" | ")[0]),
                },
            },
            {
                header: "points",
                value: (d, i) => String(d.wnba[i].points),
            },
        ]
    ).init();
}

export async function makeTeamRecordsTbl(numRows: number): Promise<void> {
    let datasrc = `${base}/teamrecs`;
    let r = await fetch(datasrc);
    const data = await r.json();
    new Tbl(
    "trtbl",
    "NBA/WNBA Regular Season Team Records", numRows, datasrc, 
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
            header: `wnba | ${data.wnba_team_records[0].season}`,
            value: (d, i) => d.wnba_team_records[i].team_long,
        },
        {
            header: "record",
            value: (d, i) => `${d.wnba_team_records[i].wins}-${d.wnba_team_records[i].losses}`,
        },
    ],
    ).init();
}