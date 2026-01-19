// /seasons
export type Season = {
    season_id: string,
    season: string,
    wseason: string,
};

export type Seasons = Season[];

// /teams
type Team = {
    league: string,
    team_id: string,
    team: string,
    team_long: string,
};

export type Teams = Team[];

// /games/recent
type recentGameTopScorer = {
    player_id: number,
    team_id: number,
    player: string,
    league: "NBA" | "WNBA",
    points: number, 
    assists: number, 
    rebounds: number,
};

type recentGame = {
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

// /league/scoring-leaders
type scoringLeader = {
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

// /teamrecs
type TeamRec = {
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


// /v2/players 
// THIS IS THE ENDPOINT THAT BUILDS PLAYER DASH
export type TopScorer = RGData | null;

type shotTypeStats = {
    made: number,
    attempted: number,
    percentage: string,
}

type shootingStats = {
    "field goals": shotTypeStats,
    "three pointers": shotTypeStats,
    "free throws": shotTypeStats,
}

type boxStats = {
    points: number,
    assists: number,
    rebounds: number,
    steals: number,
    blocks: number,
}

export type playerMeta = {
    player_id: number,
    team_id: number,
    league: string,
    season_id: number,
    player: string,
    team: string,
    team_name: string,
    season: string,
    caption: string,
    caption_short: string,
    cap_box_tot: string,
    cap_box_avg: string,
    cap_shtg_tot: string,
    cap_shtg_avg: string,
    headshot_url: string,
    team_logo_url: string,
};

type statsGroup = {
    box_stats: boxStats,
    shooting: shootingStats,
}

type playerPlaytime = {
    games_played: number,
    minutes: number,
    minutes_pg: number,
}

export type PlayerResp = {
    player_meta: playerMeta,
    playtime: playerPlaytime,
    totals: statsGroup,
    per_game: statsGroup,
}

type requestMeta = {
    request: string,
    requestedUrl: string,
    errorsOccured: number,
};

export type PlayersResp = {
    request_meta: requestMeta,
    player: PlayerResp[],
    error_string?: string,
};