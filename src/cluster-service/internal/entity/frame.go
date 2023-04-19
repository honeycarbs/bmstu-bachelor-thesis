package entity

type Frame struct {
	ID           string `db:"uuid"`
	SampleUUID   string `db:"sample_uuid"`
	Index        int    `db:"index"`
	ClusterIndex string `db:"cluster_uuid"`
	MFCCs        []float64
}
