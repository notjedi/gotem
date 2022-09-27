# Overview

## General Info

Name:               {{.Name}}
Hash:               {{.HashString}}
ID:                 {{.ID}}
Location:           {{.DownloadDir}}
Files:              {{len .Files}}
Chunks:             {{.PieceCount}};  {{.PieceSize}} each
***

## Size Info

Size:               {{humanizeBytes .SizeWhenDone.Byte}}
Downloaded:         {{humanizeBytes .HaveValid}}
Uploaded:           {{humanizeBytes .UploadedEver}}
Left until done:    {{humanizeBytes .LeftUntilDone}}
Verified:           {{humanizeBytes .HaveValid}}
Corrupt:            {{humanizeCorrupt .CorruptEver}}
Ratio:              {{.UploadRatio}}
***

## Bandwidth Info

Download limit:     {{.DownloadLimit}}
Upload limit:       {{.UploadLimit}}
Comment:            {{.Comment}}
Creator:            {{.Creator}}
{{if .IsPrivate}} Privacy:            Private torrent {{else}} Privacy:            Public torrent {{end}}
***

## Time Info

Created at:         {{humanizeTime .DateCreated}}
Added at:           {{humanizeTime .AddedDate}}
Started at:         {{humanizeTime .StartDate}}
Last activity at:   {{humanizeTime .ActivityDate}}
Completed at:       {{humanizeTime .DoneDate}}
