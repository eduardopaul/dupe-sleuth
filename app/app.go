package app

type File struct {
	Path string
	Size int64
	Hash string
	FirstBytes string
}

type AppType struct {
	Duplicates map[string][]File
}


