package store

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/jdetok/go-api-jdeko.me/pgdb"
)

func (sm *StMaps) MapPlayersCC(db *sql.DB) error {
	fmt.Println("mapping all players (concurrent workers)")

	rows, err := db.Query(pgdb.QPlayerStore)
	if err != nil {
		return err
	}
	defer rows.Close()

	// concurrency controls
	const maxWorkers = 100
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	results := make(chan *StPlayer)
	errCh := make(chan error, 1)

	go func() {
		for p := range results {
			sm.MapPlrIDDtlCC(p)
			sm.MapPlrNmDtlCC(p)
			sm.MapPlrIdNmCC(p)
			sm.MapPlrNmIdCC(p)
		}
	}()

	count := 0
	// scan rows serially, then spawn worker goroutine per player
	for rows.Next() {
		count++
		var id uint64
		var name, lowrStr, lg, tms string
		var maxR, minR, maxP, minP int

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

		sem <- struct{}{}
		wg.Add(1)
		// to make concurrent safe, need to make functions accept *StPlayer rather than StMaps
		go func(p *StPlayer, lowrStr, tms string) {
			defer wg.Done()
			defer func() { <-sem }()

			fmt.Println("worker", count, "running")

			p.Lowr = RemoveDiacritics(lowrStr)

			sm.MapPlrIdCC(p)
			sm.MapPlrNmCC(p)

			sm.MapPlrIdToSznCC(p)
			sm.MapPlrNmToSznCC(p)

			if err := sm.MapSznTmPlrCC(db, p); err != nil {
				errCh <- fmt.Errorf(
					"error occured mapping player season/teams %s | %d\n%v",
					p.Lowr, p.Id, err)
			}

			select {
			case results <- p:
			case <-errCh:
			}
		}(p, lowrStr, tms)
	}

	// check for error in rows
	if err := rows.Err(); err != nil {
		return err
	}

	// close results when all workers finish
	go func() {
		wg.Wait()
		close(results)
		fmt.Println("finished with", count, "rows")
		// logging goes here
	}()
	return nil
}

// map player name to season played by player
func (sm *StMaps) MapPlrNmToSznCC(p *StPlayer) {
	sm.mu.Lock()
	// set 0 season
	sm.SeasonPlrNms[0][p.Lowr] = p.Id
	for s := p.MinRSzn; s <= p.MaxRSzn; s++ {
		sm.SeasonPlrNms[int(s)][p.Lowr] = p.Id
	}
	sm.mu.Unlock()
}

// map player id to each season played by player
func (sm *StMaps) MapPlrIdToSznCC(p *StPlayer) {
	sm.mu.Lock()
	// set 0 season
	sm.SeasonPlrIds[0][p.Id] = p.Lowr
	for s := p.MinRSzn; s <= p.MaxRSzn; s++ {
		sm.SeasonPlrIds[int(s)][p.Id] = p.Lowr
	}
	sm.mu.Unlock()
}

// map a player id to a map of team ids that is mapped to a map of seasons
// this datastructure enables verifying whether x player played for y team in z season
func (sm *StMaps) MapSznTmPlrCC(db *sql.DB, p *StPlayer) error {
	fmt.Printf("mapping %s|%d to season team maps from %d - %d\n", p.Lowr, p.Id, p.MinRSzn, p.MaxRSzn)
	q := `
select szn_id, string_agg(distinct team_id::text, ',')
from stats.pbox
where player_id = $1 
and szn_id between $2 and $3
group by player_id, szn_id`

	tmsRows, err := db.Query(q, p.Id, p.MinRSzn, p.MaxRSzn)
	if err != nil {
		return fmt.Errorf("season team player query failed %d: %w", p.Id, err)
	}
	fmt.Println("MapSznTmPlrCC query finished, processesing rows player", p.Lowr)

	for tmsRows.Next() {
		var szn int
		var tmStr string
		err = tmsRows.Scan(&szn, &tmStr)
		if err != nil {
			return err
		}

		// scan writes a comma seperated string of team ids to tmStr
		// split to slice of strings & iterate over each
		tmsItr := strings.SplitSeq(tmStr, ",")
		for t := range tmsItr {
			// access uint64 version of team id created early in sm.TeamIDs
			teamId := sm.getTeamIDUintCC(t)

			sm.mu.Lock()
			// create empty map for the seasonid to safely create team id maps
			if sm.SznTmPlrIds[szn] == nil {
				sm.SznTmPlrIds[szn] = map[uint64]map[uint64]string{}
			}
			// create empty map for each team id in each season
			if sm.SznTmPlrIds[szn][teamId] == nil {
				sm.SznTmPlrIds[szn][teamId] = map[uint64]string{}
			}
			// add this player's id to the corresponding season/team inner map
			sm.SznTmPlrIds[szn][teamId][p.Id] = p.Lowr
			sm.mu.Unlock()
		}

	}
	return nil
}

// access a teamId from sm.TeamIds concurrently
func (sm *StMaps) getTeamIDUintCC(t string) uint64 {
	sm.mu.RLock()
	teamId := sm.TeamIds[t]
	sm.mu.RUnlock()
	return teamId
}

// insert player id and cleaned player name as keys in PlrIds and PlrNms maps
func (sm *StMaps) MapPlrIdCC(p *StPlayer) {
	sm.mu.Lock()
	sm.PlrIds[p.Id] = struct{}{}
	sm.mu.Unlock()
}

func (sm *StMaps) MapPlrNmCC(p *StPlayer) {
	sm.mu.Lock()
	sm.PlrNms[p.Lowr] = struct{}{}
	sm.mu.Unlock()
}

func (sm *StMaps) MapPlrIDDtlCC(p *StPlayer) {
	sm.mu.Lock()
	sm.PlayerIdDtl[p.Id] = p
	sm.mu.Unlock()
}

func (sm *StMaps) MapPlrNmDtlCC(p *StPlayer) {
	sm.mu.Lock()
	sm.PlayerNameDtl[p.Lowr] = p
	sm.mu.Unlock()
}

func (sm *StMaps) MapPlrNmIdCC(p *StPlayer) {
	sm.mu.Lock()
	sm.PlayerNameId[p.Lowr] = p.Id
	sm.mu.Unlock()
}

func (sm *StMaps) MapPlrIdNmCC(p *StPlayer) {
	sm.mu.Lock()
	sm.PlayerIdName[p.Id] = p.Lowr
	sm.mu.Unlock()
}
