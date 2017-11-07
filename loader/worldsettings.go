package loader

func LoadWorldSettings() (settings map[string]string) {
	settings = make(map[string]string)
	// settings for the entire world
	settings["void.zone"] = "void"
	settings["void.room"] = "void"
	settings["start.zone"] = "sample"
	settings["start.room"] = "start"
	return
}
