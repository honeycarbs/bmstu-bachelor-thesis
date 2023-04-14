package entity

type Frame struct {
	ID           string `db:"uuid"`
	SampleHash   string `db:"sample_hash"`
	Index        int    `db:"index"`
	ClusterIndex string `db:"cluster_uuid"`
	MFCCs        []float64
}
