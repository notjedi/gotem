package common

/*
i can prolly use package specific fields to squeeze a tiny tiny amount of performance.
since the response to the request made comes from c, i can kinda assume that adding
more fields is basically free. so the performance gain comes down to json serialization and
deserialization? and as go is also kinda fast, ig the performance gain here is immeasurable?
*/
var torrentFields = []string{
	"id", "hashString", "name", "status", "rateDownload", "rateUpload", "eta", "uploadRatio",
	"sizeWhenDone", "haveValid", "haveUnchecked", "uploadedEver", "recheckProgress",
	"peersConnected", "uploadLimited", "downloadLimited", "bandwidthPriority", "peersSendingToUs",
	"peersGettingFromUs", "seedRatioLimit", "trackerStats", "magnetLink", "honorsSessionLimits",
	"metadataPercentComplete", "percentDone", "downloadDir", "files", "pieceCount", "pieceSize",
	"leftUntilDone", "corruptEver", "downloadLimit", "uploadLimit", "comment", "creator",
	"isPrivate", "dateCreated", "addedDate", "startDate", "activityDate", "doneDate",
}
