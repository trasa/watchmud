package loader

type raceFileEntry struct {
	Id            string `json:"id"`
	RaceGroupName string `json:"race_group_name"`
	RaceName      string `json:"race_name"`
	StrBonus      int32  `json:"str_bonus"`
	DexBonus      int32  `json:"dex_bonus"`
	ConBonus      int32  `json:"con_bonus"`
	IntBonus      int32  `json:"int_bonus"`
	WisBonus      int32  `json:"wis_bonus"`
	ChaBonus      int32  `json:"cha_bonus"`
}
