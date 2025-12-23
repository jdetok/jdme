package memd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jdetok/go-api-jdeko.me/pkg/clnd"
	"github.com/jdetok/go-api-jdeko.me/pkg/logd"
	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
	"golang.org/x/sync/errgroup"
)

// query all players from DB, loop through rows serially, then launch goroutines
// to process and map the data for each player concurrently
func (sm *StMaps) MapPlayersCC(ctx context.Context, db pgdb.DB, lgd *logd.Logd) error {
	start := time.Now()
	lgd.Debugf("mapping all players (concurrent workers)")

	rows, err := db.QueryContext(ctx, pgdb.QPlayerStore)
	if err != nil {
		return err
	}
	defer rows.Close()

	// concurrency controls

	const maxWorkers = 20
	sem := make(chan struct{}, maxWorkers)
	errCh := make(chan error, maxWorkers)
	g, ctx := errgroup.WithContext(ctx)
	results := make(chan *StPlayer, maxWorkers)

	sm.PlayerNameId["random"] = 77777

	// read the results channel, add player to maps
	go func() {
		// WAITGROUP SHOULD NOT WAIT HERE. func hangs if so
		defer func() {
			if r := recover(); r != nil {
				select {
				case errCh <- fmt.Errorf("results reader panicing: %v", r):
				default:
				}
			}
		}()

		for p := range results {
			lgd.Debugf("%s complete", p.Name)
		}
	}()

	// var wg = &sync.WaitGroup{}
	count := 0
	// scan rows serially, then spawn worker goroutine per player
	for rows.Next() {
		if ctx.Err() != nil {
			break
		}

		count++
		var id uint64
		var name, lowrStr, tms string
		var maxR, minR, maxP, minP, lg int

		if err := rows.Scan(&id, &name, &lowrStr, &lg, &maxR, &minR, &maxP, &minP, &tms); err != nil {
			return err
		}

		p := &StPlayer{
			Id:      id,
			Name:    name,
			Lg:      lg,
			MaxRSzn: maxR,
			MinRSzn: minR,
			MaxPSzn: maxP,
			MinPSzn: minP,
		}

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			return ctx.Err()
		}
		// capture loop vars
		l := lowrStr
		t := tms
		wrkNum := count

		// wg.Add(1)
		g.Go(func() error {
			// defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					select {
					case errCh <- fmt.Errorf("worker %d panicing: %v", wrkNum, r):
					default:
					}
				}
			}()

			defer func() { <-sem }()
			if ctx.Err() != nil {
				return ctx.Err()
			}

			lgd.Debugf("worker %d running\n", wrkNum)

			// clean the player name (lower case, remove accents)
			p.Lowr = clnd.ConvToASCII(l)

			// split comma separated string with teams to p.Teams
			tmIds, err := sm.SplitTeams(p, t)
			if err != nil {
				return err
			}
			p.Teams = tmIds

			sm.MapPlrIDDtlCC(p)
			sm.MapPlrNmDtlCC(p)
			sm.MapPlrIdNmCC(p)
			sm.MapPlrNmIdCC(p)
			// player exists maps
			sm.MapPlrIdCC(p)
			sm.MapPlrNmCC(p)

			// query team(s) played for each season from min-max player season,
			// add player to map for each team played for in each season played
			if err := sm.MapSeasonTeamPlayers(lgd, db, p); err != nil {
				return fmt.Errorf(
					"error occured mapping player season/teams %s | %d\n%v",
					p.Lowr, p.Id, err)
			}
			select {
			case results <- p:
			case <-ctx.Done():
			}
			return nil
		})
	}

	// check for error in rows
	if err := rows.Err(); err != nil {
		return err
	}

	// close results when all workers finish
	// wg.Wait()
	g.Wait()
	close(results)
	lgd.Debugf("finished with %d rows after %v", count, time.Since(start))

	return nil
}

// func (sm *StMaps) MapPlrIDDtlCC() {

// }

// map a player id to a map of team ids that is mapped to a map of seasons
// this datastructure enables verifying whether x player played for y team in z season
func (sm *StMaps) MapSeasonTeamPlayers(lg *logd.Logd, db pgdb.DB, p *StPlayer) error {
	lg.Debugf("mapping %s|%d to season team maps from %d - %d\n", p.Lowr, p.Id, p.MinRSzn, p.MaxRSzn)
	q := ` 
select szn_id, string_agg(distinct team_id::text, ',')
from stats.pbox
where player_id = $1 
and substr(szn_id::text, 2, 4)::int between 
substr($2::text, 2, 4)::int and substr($3::text, 2, 4)::int
group by player_id, szn_id
`

	rows, err := db.Query(q, p.Id, p.MinRSzn, p.MaxRSzn)
	if err != nil {
		return fmt.Errorf("season team player query failed %d: %w", p.Id, err)
	}

	for rows.Next() {
		var szn int
		var tmStr string
		err = rows.Scan(&szn, &tmStr)
		if err != nil {
			return err
		}

		// scan writes a comma seperated string of team ids to tmStr
		// split to slice of strings & iterate over each
		tmsItr := strings.SplitSeq(tmStr, ",")
		for t := range tmsItr {
			// access uint64 version of team id created early in sm.TeamIDs
			teamId, err := sm.GetTeamIDUintCC(t)
			if err != nil {
				return err
			}
			sm.mu.Lock()
			sm.MapSznPlr(szn, p)
			sm.MapSznTmPlr(szn, teamId, p)
			sm.mu.Unlock()
		}

	}
	return nil
}

func (sm *StMaps) MapSznPlr(szn int, p *StPlayer) {
	var szns = []int{0, 29999, 49999, szn}
	for _, s := range szns {
		sm.SeasonPlrIds[s][p.Id] = p.Lowr
		sm.SeasonPlrNms[s][p.Lowr] = p.Id
	}
}

func (sm *StMaps) MapSznTmPlr(szn int, tId uint64, p *StPlayer) {
	var plOff bool = (szn >= 40000 && szn < 50000)
	var plOffSzn int = 49999
	var rgSzn int = 29999

	// season team player
	sm.SznTmPlrIds[szn][tId][p.Id] = p.Lowr
	sm.SznTmPlrIds[szn][0][p.Id] = p.Lowr
	switch p.Lg {
	case 0:
		sm.NSznTmPlrIds[szn][tId][p.Id] = p.Lowr
		sm.NSznTmPlrIds[szn][0][p.Id] = p.Lowr
		if plOff {
			sm.NSznTmPlrIds[plOffSzn][tId][p.Id] = p.Lowr
			sm.NSznTmPlrIds[plOffSzn][0][p.Id] = p.Lowr
		} else {
			sm.NSznTmPlrIds[rgSzn][tId][p.Id] = p.Lowr
			sm.NSznTmPlrIds[rgSzn][0][p.Id] = p.Lowr
		}
	case 1:
		sm.WSznTmPlrIds[szn][tId][p.Id] = p.Lowr
		sm.WSznTmPlrIds[szn][0][p.Id] = p.Lowr
		if plOff {
			sm.WSznTmPlrIds[plOffSzn][tId][p.Id] = p.Lowr
			sm.WSznTmPlrIds[plOffSzn][0][p.Id] = p.Lowr
		} else {
			sm.WSznTmPlrIds[rgSzn][tId][p.Id] = p.Lowr
			sm.WSznTmPlrIds[rgSzn][0][p.Id] = p.Lowr
		}
	}
}

// playoff safe copy
func (sm *StMaps) MapSznTmPlPO(lg *logd.Logd, db pgdb.DB, p *StPlayer) error {
	lg.Debugf("mapping %s|%d to season team maps from %d - %d\n", p.Lowr, p.Id, p.MinRSzn, p.MaxRSzn)
	q := `
select szn_id, string_agg(distinct team_id::text, ',')
from stats.pbox
where player_id = $1 
and substr(szn_id::text, 2, 4)::int between 
substr($2::text, 2, 4)::int and substr($3::text, 2, 4)::int
group by player_id, szn_id
	`

	tmsRows, err := db.Query(q, p.Id, p.MinRSzn, p.MaxRSzn)
	if err != nil {
		return fmt.Errorf("season team player query failed %d: %w", p.Id, err)
	}

	for tmsRows.Next() {
		var szn int
		var tmStr string
		err = tmsRows.Scan(&szn, &tmStr)
		if err != nil {
			return err
		}
		//
		// scan writes a comma seperated string of team ids to tmStr
		// split to slice of strings & iterate over each
		tmsItr := strings.SplitSeq(tmStr, ",")
		for t := range tmsItr {
			// access uint64 version of team id created early in sm.TeamIDs
			teamId, err := sm.GetTeamIDUintCC(t)
			if err != nil {
				return err
			}

			sm.mu.Lock()
			// create empty map for the seasonid to safely create team id maps
			if sm.SznTmPlrIds[szn] == nil {
				sm.SznTmPlrIds[szn] = map[uint64]map[uint64]string{}
			}
			// create empty map for each team id in each season
			if sm.SznTmPlrIds[szn][teamId] == nil {
				sm.SznTmPlrIds[szn][teamId] = map[uint64]string{}
			}

			if sm.SznTmPlrIds[szn][0] == nil {
				sm.SznTmPlrIds[szn][0] = map[uint64]string{}
			}
			if sm.SznTmPlrIds[szn][1] == nil {
				sm.SznTmPlrIds[szn][1] = map[uint64]string{}
			}

			// add this player's id to the corresponding season/team inner map
			sm.SznTmPlrIds[szn][teamId][p.Id] = p.Lowr
			sm.SznTmPlrIds[szn][0][p.Id] = p.Lowr

			switch p.Lg {
			case 0:
				if sm.NSznTmPlrIds[szn] == nil {
					sm.NSznTmPlrIds[szn] = map[uint64]map[uint64]string{}
				}
				if sm.NSznTmPlrIds[szn][teamId] == nil {
					sm.NSznTmPlrIds[szn][teamId] = map[uint64]string{}
				}
				if sm.NSznTmPlrIds[szn][0] == nil {
					sm.NSznTmPlrIds[szn][0] = map[uint64]string{}
				}
				sm.NSznTmPlrIds[szn][teamId][p.Id] = p.Lowr
				sm.NSznTmPlrIds[szn][0][p.Id] = p.Lowr

			case 1:
				if sm.WSznTmPlrIds[szn] == nil {
					sm.WSznTmPlrIds[szn] = map[uint64]map[uint64]string{}
				}
				if sm.WSznTmPlrIds[szn][teamId] == nil {
					sm.WSznTmPlrIds[szn][teamId] = map[uint64]string{}
				}
				if sm.WSznTmPlrIds[szn][0] == nil {
					sm.WSznTmPlrIds[szn][0] = map[uint64]string{}
				}
				sm.WSznTmPlrIds[szn][teamId][p.Id] = p.Lowr
				sm.WSznTmPlrIds[szn][0][p.Id] = p.Lowr
			}
			sm.mu.Unlock()
		}

	}
	return nil
}

// access a teamId from sm.TeamIds concurrently
func (sm *StMaps) GetTeamIDUintCC(t string) (uint64, error) {
	sm.mu.RLock()
	var teamId uint64
	var ok bool
	var err error
	if teamId, ok = sm.TeamIds[t]; !ok {
		// convert and add to map if doesn't already exist
		teamId, err = strconv.ParseUint(t, 10, 64)
		if err != nil {
			return 0, err
		}
		sm.TeamIds[t] = teamId
	}
	sm.mu.RUnlock()
	return teamId, nil
}

// add player id as key to sm.PlrIds
func (sm *StMaps) MapPlrIdCC(p *StPlayer) {
	sm.mu.Lock()
	sm.PlrIds[p.Id] = struct{}{}
	sm.mu.Unlock()
}

// add cleaned player name as key to sm.PlrIds
func (sm *StMaps) MapPlrNmCC(p *StPlayer) {
	sm.mu.Lock()
	sm.PlrNms[p.Lowr] = struct{}{}
	sm.mu.Unlock()
}

// map StPlayer struct to player id
func (sm *StMaps) MapPlrIDDtlCC(p *StPlayer) {
	sm.mu.Lock()
	sm.PlayerIdDtl[p.Id] = p
	sm.mu.Unlock()
}

// map StPlayer struct to cleaned player name
func (sm *StMaps) MapPlrNmDtlCC(p *StPlayer) {
	sm.mu.Lock()
	sm.PlayerNameDtl[p.Lowr] = p
	sm.mu.Unlock()
}

// map player id to cleaned player name
func (sm *StMaps) MapPlrNmIdCC(p *StPlayer) {
	sm.mu.Lock()
	sm.PlayerNameId[p.Lowr] = p.Id
	sm.mu.Unlock()
}

// map cleaned player name to player id
func (sm *StMaps) MapPlrIdNmCC(p *StPlayer) {
	sm.mu.Lock()
	sm.PlayerIdName[p.Id] = p.Lowr
	sm.mu.Unlock()
}
