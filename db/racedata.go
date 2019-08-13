package db

type RaceData struct {
	RaceId        int32  `db:"race_id"`
	RaceGroupName string `db:"race_group_name"`
	RaceName      string `db:"race_name"`
	StrBonus int32 `db:"str_bonus"`
	DexBonus int32 `db:"dex_bonus"`
	ConBonus int32 `db:"con_bonus"`
	IntBonus int32 `db:"int_bonus"`
	WisBonus int32 `db:"wis_bonus"`
	ChaBonus int32 `db:"cha_bonus"`
}


func GetRaceData() (result []RaceData, err error) {
	err = watchdb.Select(&result, "select r.race_id, rg.race_group_name, r.race_name, r.str_bonus, r.dex_bonus, r.con_bonus, r.int_bonus, r.wis_bonus, r.cha_bonus " +
		"from races r " +
		"inner join race_group rg on r.race_group_id = rg.race_group_id " +
		"order by r.race_id")
	return
}
