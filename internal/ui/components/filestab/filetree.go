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
	bytesTotal     int64
	bytesCompleted int64
}

func buildFileTree(files []*transmissionrpc.TorrentFile,
	fileStats []*transmissionrpc.TorrentFileStat,
) *Directory {
	if len(files) == 0 {
		return &Directory{}
	}

	dirTree := &Directory{
		name: strings.Split(files[0].Name, "/")[0],
	}

	for _, file := range files {
		currentDir := dirTree
		dirs := strings.Split(file.Name, "/")

		for idx, dir := range dirs[1:] {
			if idx == len(dirs)-2 {
				// TODO: add priority
				file := &File{
					name:           dir,
					bytesTotal:     file.Length,
					bytesCompleted: file.BytesCompleted,
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
