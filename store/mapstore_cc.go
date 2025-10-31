package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

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
	const maxWorkers = 10
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	results := make(chan *StPlayer)
	errCh := make(chan error, 1)

	// consumerDone := make(chan struct{})

	// single-threaded insertion into maps
	go func() {

		for p := range results {
			fmt.Println("consumer received", p.Id)

			sm.MapPlrIDDtlCC(p)
			sm.MapPlrNmDtlCC(p)
			sm.MapPlrIdNmCC(p)
			sm.MapPlrNmIdCC(p)
		}
		// close(consumerDone)
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

			sm.MapPlrIdCC(p)
			sm.MapPlrNmCC(p)

			// per-player processing (safe to run concurrently)
			p.Lowr = RemoveDiacritics(lowrStr)

			// sm.StoreTeamIDUintCC(p, tms)

			sm.MapPlrIdToSznCC(p)
			sm.MapPlrNmToSznCC(p)

			fmt.Println("worker before MapSznTmPlrCC", p.Id)
			if err := sm.MapSznTmPlrCC(db, p); err != nil {
				fmt.Printf("error occured mapping player season/teams %s | %d\n%v", p.Lowr, p.Id, err)
				// non-fatal: continue
			}
			fmt.Println("worker after MapSznTmPlrCC", p.Id)
			fmt.Println("worker sending result", p.Id)

			select {
			case results <- p:
			case <-errCh:
			}
			// send completed player for single-threaded map insertion
			// results <- p
		}(p, lowrStr, tms)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	// close results when all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()
	// <-consumerDone
	// check for worker-reported error
	select {
	case err := <-errCh:
		return err
	case p := <-results:
		fmt.Println("done with", p.Lowr)
		return nil
	default:
		return nil
	}
}

func (sm *StMaps) StoreTeamIDUintCC(p *StPlayer, tms string) {
	// split tms string from db to slice of strings
	teamsStrArr := strings.SplitSeq(tms, ",")

	// iterate through each team player has played for
	// TODO: map players to team map similar to season maps
	for t := range teamsStrArr {
		// use TeamIds map created with MakeMaps() get the uint64 version of t
		// append to teams slice
		teamId := sm.getTeamIDUintCC(t)
		p.Teams = append(p.Teams, teamId)
	}
}

func (sm *StMaps) MapPlrToSznCC(p *StPlayer) error {
	fmt.Printf("mapping %s|%d to season maps from %d - %d\n", p.Lowr, p.Id, p.MinRSzn, p.MaxRSzn)
	sm.mu.Lock()
	sm.SeasonPlrNms[0][p.Lowr] = p.Id
	sm.SeasonPlrIds[0][p.Id] = p.Lowr
	for s := p.MinRSzn; s <= p.MaxRSzn; s++ {
		sm.SeasonPlrNms[int(s)][p.Lowr] = p.Id
		sm.SeasonPlrIds[int(s)][p.Id] = p.Lowr
	}
	sm.mu.Unlock()
	return nil
}

func (sm *StMaps) MapPlrNmToSznCC(p *StPlayer) {
	sm.mu.Lock()
	sm.SeasonPlrNms[0][p.Lowr] = p.Id
	for s := p.MinRSzn; s <= p.MaxRSzn; s++ {
		sm.SeasonPlrNms[int(s)][p.Lowr] = p.Id
	}
	sm.mu.Unlock()
}

func (sm *StMaps) MapPlrIdToSznCC(p *StPlayer) {
	sm.mu.Lock()
	sm.SeasonPlrIds[0][p.Id] = p.Lowr
	for s := p.MinRSzn; s <= p.MaxRSzn; s++ {
		sm.SeasonPlrIds[int(s)][p.Id] = p.Lowr
	}
	sm.mu.Unlock()
}

func (sm *StMaps) MapSznTmPlrCC(db *sql.DB, p *StPlayer) error {
	fmt.Printf("mapping %s|%d to season team maps from %d - %d\n", p.Lowr, p.Id, p.MinRSzn, p.MaxRSzn)
	q := `
select szn_id, string_agg(distinct team_id::text, ',')
from stats.pbox
where player_id = $1 
and szn_id between $2 and $3
group by player_id, szn_id`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("MapSznTmPlrCC: query start player", p.Id)
	tmsRows, err := db.QueryContext(ctx, q, p.Id, p.MinRSzn, p.MaxRSzn)
	if err != nil {
		return fmt.Errorf("query context failed player %d: %w", p.Id, err)
	}
	defer tmsRows.Close()
	fmt.Println("MapSznTmPlrCC query finished, processesing rows player", p.Lowr)

	for tmsRows.Next() {
		var szn int
		var tmStr string
		err = tmsRows.Scan(&szn, &tmStr)
		if err != nil {
			return err
		}

		tmsItr := strings.SplitSeq(tmStr, ",")
		for t := range tmsItr {
			teamId := sm.getTeamIDUintCC(t)

			sm.mu.Lock()
			if sm.SznTmPlrIds[szn] == nil {
				sm.SznTmPlrIds[szn] = map[uint64]map[uint64]string{}
			}
			// ensure inner map for team exists
			if sm.SznTmPlrIds[szn][teamId] == nil {
				sm.SznTmPlrIds[szn][teamId] = map[uint64]string{}
			}
			sm.SznTmPlrIds[szn][teamId][p.Id] = p.Lowr
			sm.mu.Unlock()

		}

	}
	return nil
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

func (sm *StMaps) getTeamIDUintCC(t string) uint64 {
	sm.mu.RLock()
	teamId := sm.TeamIds[t]
	sm.mu.RUnlock()
	return teamId
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
