package resp

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jdetok/go-api-jdeko.me/pkg/memd"
	"github.com/jdetok/go-api-jdeko.me/pkg/pgdb"
)

// after verifying player exists, query db for their stats
func (r *RespPlayerDash) GetPlayerDashV2(db *sql.DB, sm *memd.StMaps, iq *PQueryIds) ([]byte, error) {

	// query player, scan to structs, call struct functions
	// appends RespObj to r.Results
	if err := r.BuildPlayerRespV2(db, sm, iq); err != nil {
		msg := fmt.Sprintf("failed to query playerId %d seasonId %d", iq.PId, iq.SId)
		return nil, fmt.Errorf("%s\n%v", msg, err)
	}

	// marshall Resp struct to JSON, return as []byte
	js, err := json.Marshal(r)
	if err != nil {
		msg := "failed to marshal structs to json"
		return nil, fmt.Errorf("%s\n%v", msg, err)
	}
	return js, nil
}

func (r *RespPlayerDash) BuildPlayerRespV2(db *sql.DB, sm *memd.StMaps, iq *PQueryIds) error {
	pOrT := "plr"
	q := pgdb.TmSznPlr
	args := []any{iq.PId, iq.TId, iq.SId}

	if iq.SId == 0 || iq.SId == 88888 {
		maxSzn, err := sm.GetSznFromPlrId(iq.PId)
		if err != nil {
			return err
		}
		iq.SId = maxSzn
	}

	if iq.TId == 0 {
		q = pgdb.PlayerDash
		args = []any{iq.PId, iq.SId}
		fmt.Println("team 0 args:", args)
	} else {
		if strconv.Itoa(iq.SId)[1:] == "9999" {
			args = []any{iq.PId, iq.TId, strconv.Itoa(iq.SId)[0]}
			q = pgdb.AggPlrTm
		}
	}

	if err := r.ProcessRows(db, pOrT, q, args...); err != nil {
		return err
	}
	return nil
}
