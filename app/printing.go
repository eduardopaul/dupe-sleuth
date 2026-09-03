package app

import "fmt"

func getPlural(num int) string {
	if num == 1 {
		return ""
	} else {
		return "s"
	}
}

func PrintGroups(filesByHash map[string][]File) {
	numGroups := len(filesByHash)
	numFiles := 0
	for _, sliceOfFile := range filesByHash {
		numFiles += len(sliceOfFile)
	}

	fmt.Printf(
		"\nFound %d duplicate group%s, totaling %d file%s.\n",
		numGroups,
		getPlural(numGroups),
		numFiles,
		getPlural(numFiles),
	)

	fmt.Println()
	fileIdx := 1
	for _, sliceOfFile := range filesByHash {
		PrintGroup(sliceOfFile, fileIdx)
	}
}

func PrintGroup(group []File, fileIdx int) {
	nFiles := len(group)
	size := group[0].Size

	fmt.Printf(
		"Group %d, %d files, %d byte%s each:\n",
		fileIdx,
		nFiles,
		size,
		getPlural(int(size)),
	)

	fileIdx++

	for idx, file := range group {
		fmt.Printf("  %d. %s\n", idx+1, file.Path)
	}
	fmt.Println()
}

