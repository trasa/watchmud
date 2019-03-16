package db

type RaceData struct {
	RaceId        int32  `db:"race_id"`
	RaceGroupName string `db:"race_group_name"`
	RaceName      string `db:"race_name"`
}

func GetRaceData() (result []RaceData, err error) {
	err = watchdb.Select(&result, "select r.race_id, rg.race_group_name, r.race_name from races r inner join race_group rg on r.race_group_id = rg.race_group_id order by r.race_id")
	return
}
