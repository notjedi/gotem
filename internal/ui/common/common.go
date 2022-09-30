package common

var allTorrentInfoFields = []string{
	"id", "hashString", "name", "status", "rateDownload", "rateUpload", "eta", "uploadRatio",
	"sizeWhenDone", "haveValid", "haveUnchecked", "uploadedEver", "recheckProgress",
	"peersConnected", "trackerStats", "metadataPercentComplete", "percentDone", "queuePosition",
}

var torrentInfoFields = []string{
	"id", "hashString", "name", "uploadRatio", "sizeWhenDone", "haveValid", "uploadedEver",
	"uploadLimited", "downloadLimited", "magnetLink", "downloadDir", "files", "pieceCount",
	"pieceSize", "leftUntilDone", "corruptEver", "downloadLimit", "uploadLimit", "comment",
	"creator", "isPrivate", "dateCreated", "addedDate", "startDate", "activityDate", "doneDate",
}
