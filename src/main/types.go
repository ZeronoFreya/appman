package main

type SL struct {
	Dirpath string `json:"dirpath"`
	Symlink string `json:"symlink"`
}
type SymlinkData struct {
	Dir  string `json:"dir"`
	Link []SL   `json:"link"`
}
type Config struct {
	Nouac   []string                 `json:"nouac"`
	Symlink map[string][]SymlinkData `json:"symlink"`
}
