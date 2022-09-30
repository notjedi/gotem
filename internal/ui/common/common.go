package common

// TODO: remove unwnated fields
var allTorrentInfoFields = []string{
	"id", "hashString", "name", "status", "rateDownload", "rateUpload", "eta", "uploadRatio",
	"sizeWhenDone", "haveValid", "haveUnchecked", "uploadedEver", "recheckProgress",
	"peersConnected", "uploadLimited", "downloadLimited", "bandwidthPriority", "peersSendingToUs",
	"peersGettingFromUs", "seedRatioLimit", "trackerStats", "magnetLink", "honorsSessionLimits",
	"metadataPercentComplete", "percentDone", "downloadDir", "files", "pieceCount", "pieceSize",
	"leftUntilDone", "corruptEver", "downloadLimit", "uploadLimit", "comment", "creator",
	"isPrivate", "dateCreated", "addedDate", "startDate", "activityDate", "doneDate", "queuePosition",
}

// TODO: remove unwnated fields
var torrentInfoFields = []string{
	"id", "hashString", "name", "status", "rateDownload", "rateUpload", "eta", "uploadRatio",
	"sizeWhenDone", "haveValid", "haveUnchecked", "uploadedEver", "recheckProgress",
	"peersConnected", "uploadLimited", "downloadLimited", "bandwidthPriority", "peersSendingToUs",
	"peersGettingFromUs", "seedRatioLimit", "trackerStats", "magnetLink", "honorsSessionLimits",
	"metadataPercentComplete", "percentDone", "downloadDir", "files", "pieceCount", "pieceSize",
	"leftUntilDone", "corruptEver", "downloadLimit", "uploadLimit", "comment", "creator",
	"isPrivate", "dateCreated", "addedDate", "startDate", "activityDate", "doneDate",
}
