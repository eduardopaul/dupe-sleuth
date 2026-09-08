package app

type File struct {
	Path string
	Size int64
	Hash string
	FirstBytes string
}

type State struct {
	Duplicates map[string][]File
	Order []string
	Marked map[string]File
}

