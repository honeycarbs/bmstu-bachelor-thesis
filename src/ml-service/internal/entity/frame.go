package entity

type Frame struct {
	ID           string `db:"uuid"`
	SampleID     string `db:"sample_uuid"`
	Index        int    `db:"frame_index"`
	ClusterIndex int    `db:"cluster_index"`
}
