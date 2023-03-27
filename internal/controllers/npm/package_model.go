package npm

//https://github.com/npm/registry/blob/master/docs/responses/package-metadata.md

import "time"

type (
	Meta struct {
		Name     string             `json:"name"`
		DistTags DistTags           `json:"dist-tags"`
		Versions map[string]Version `json:"versions"`
		Modified time.Time          `json:"modified"`
	}
	DistTags struct {
		Latest string `json:"latest"`
	}
	Signatures struct {
		Keyid string `json:"keyid"`
		Sig   string `json:"sig"`
	}
	Dist struct {
		Shasum       string       `json:"shasum"`
		Integrity    string       `json:"integrity"`
		Tarball      string       `json:"tarball"`
		FileCount    int          `json:"fileCount"`
		UnpackedSize int          `json:"unpackedSize"`
		Signatures   []Signatures `json:"signatures"`
		NpmSignature string       `json:"npm-signature"`
	}
	Version struct {
		Name             string            `json:"name"`
		Version          string            `json:"version"`
		Dependencies     map[string]string `json:"dependencies"`
		PeerDependencies map[string]string `json:"peerDependencies"`
		Dist             Dist              `json:"dist"`
		Funding          string            `json:"funding"`
	}
)

type (
	DBInfo struct {
		DbName             string `json:"db_name"`
		Engine             string `json:"engine"`
		DocCount           int    `json:"doc_count"`
		DocDelCount        int    `json:"doc_del_count"`
		UpdateSeq          int    `json:"update_seq"`
		PurgeSeq           int    `json:"purge_seq"`
		CompactRunning     bool   `json:"compact_running"`
		Sizes              Sizes  `json:"sizes"`
		DiskSize           int64  `json:"disk_size"`
		DataSize           int64  `json:"data_size"`
		Other              Other  `json:"other"`
		InstanceStartTime  string `json:"instance_start_time"`
		DiskFormatVersion  int    `json:"disk_format_version"`
		CommittedUpdateSeq int    `json:"committed_update_seq"`
		CompactedSeq       int    `json:"compacted_seq"`
		UUID               string `json:"uuid"`
	}
	Sizes struct {
		Active   int64 `json:"active"`
		External int64 `json:"external"`
		File     int64 `json:"file"`
	}
	Other struct {
		DataSize int64 `json:"data_size"`
	}
)

type (
	Publish struct {
		ID          string                `json:"_id"`
		Name        string                `json:"name"`
		DistTags    DistTags              `json:"dist-tags"`
		Versions    map[string]Version    `json:"versions"`
		Readme      string                `json:"readme"`
		Attachments map[string]Attachment `json:"_attachments"`
	}

	Attachment struct {
		ContentType string `json:"content_type"`
		Data        string `json:"data"`
		Length      int    `json:"length"`
	}
)
