package e621

type E621Post struct {
	ID   int    `json:"id"`
	URL  string `json:"url"`
	Ext  string `json:"file_ext"`
	Size int    `json:"file_size"`
}

type E621PostResponse struct {
	ID   int `json:"id"`
	File struct {
		URL  string `json:"url"`
		Ext  string `json:"ext"`
		Size int    `json:"size"`
	} `json:"file"`
	Sample struct {
		URL  string `json:"url"`
		Has  bool   `json:"has"`
		Alts map[string]struct {
			URLs []*string `json:"urls"`
		} `json:"alternates"`
	} `json:"sample"`
}
