package filestab

import (
	"strings"

	"github.com/hekmon/transmissionrpc/v2"
)

type Directory struct {
	name     string
	files    []*File
	children []*Directory
}

type File struct {
	name           string
	priority       int64
	bytesTotal     uint64
	bytesCompleted uint64
	percentDone    float64
}

func buildFileTree(files []*transmissionrpc.TorrentFile,
	fileStats []*transmissionrpc.TorrentFileStat,
) *Directory {
	if len(files) == 0 {
		return &Directory{}
	}

	dirTree := &Directory{
		name: "/",
	}

	for fileIdx, file := range files {
		currentDir := dirTree
		dirs := strings.Split(file.Name, "/")

		for idx, dir := range dirs {
			if idx == len(dirs)-1 {
				file := &File{
					name:           dir,
					bytesTotal:     uint64(file.Length),
					bytesCompleted: uint64(file.BytesCompleted),
					percentDone:    (float64(file.BytesCompleted) / float64(file.Length)) * 100,
					priority:       fileStats[fileIdx].Priority,
				}
				currentDir.files = append(currentDir.files, file)
			} else {
				found := false
				for _, child := range currentDir.children {
					if child.name == dir {
						currentDir = child
						found = true
						break
					}
				}

				if !found {
					child := &Directory{
						name: dir,
					}
					currentDir.children = append(currentDir.children, child)
					currentDir = child
				}
			}
		}
	}

	return dirTree
}
