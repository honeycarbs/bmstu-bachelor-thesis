package kmeans

import "math"

//func SilhouetteCoefficient(observations []Node, centroids []Node) float64 {
//	numClusters := len(centroids)
//	numObservations := len(observations)
//	clusterAssignments := make([]int, numObservations)
//	distances := make([][]float64, numObservations)
//	for i := 0; i < numObservations; i++ {
//		distances[i] = make([]float64, numClusters)
//		for j := 0; j < numClusters; j++ {
//			distances[i][j] = distance(observations[i], centroids[j])
//		}
//		minDist := distances[i][0]
//		minIndex := 0
//		for j := 1; j < numClusters; j++ {
//			if distances[i][j] < minDist {
//				minDist = distances[i][j]
//				minIndex = j
//			}
//		}
//		clusterAssignments[i] = minIndex
//	}
//
//	a := make([]float64, numObservations)
//	b := make([]float64, numObservations)
//
//	for i := 0; i < numObservations; i++ {
//		clusterIndex := clusterAssignments[i]
//
//		var sumDist float64
//		numSameCluster := 0
//		for j := 0; j < numObservations; j++ {
//			if clusterAssignments[j] == clusterIndex && i != j {
//				sumDist += distances[i][clusterIndex]
//				numSameCluster++
//			}
//		}
//		if numSameCluster > 0 {
//			a[i] = sumDist / float64(numSameCluster)
//		}
//
//		minAvgDist := math.Inf(1)
//		for j := 0; j < numClusters; j++ {
//			if j != clusterIndex {
//				var sumDist float64
//				numOtherCluster := 0
//				for k := 0; k < numObservations; k++ {
//					if clusterAssignments[k] == j {
//						sumDist += distances[i][j]
//						numOtherCluster++
//					}
//				}
//				if numOtherCluster > 0 {
//					avgDist := sumDist / float64(numOtherCluster)
//					if avgDist < minAvgDist {
//						minAvgDist = avgDist
//					}
//				}
//			}
//		}
//		b[i] = minAvgDist
//	}
//
//	var sumSC float64
//	for i := 0; i < numObservations; i++ {
//		if b[i] != 0 {
//			sc := (b[i] - a[i]) / math.Max(a[i], b[i])
//			sumSC += sc
//		}
//	}
//	if numObservations > 0 {
//		return sumSC / float64(numObservations)
//	} else {
//		return 0.0
//	}
//}

func SilhouetteCoefficient(observations []Node, centroids []Node) float64 {
	numClusters := len(centroids)
	numObservations := len(observations)
	clusterAssignments, distances := assignToCluster(observations, centroids)

	a := computeAScores(observations, clusterAssignments, distances)
	b := computeBScores(observations, centroids, clusterAssignments, numClusters)

	var sumSC float64
	for i := 0; i < numObservations; i++ {
		if b[i] != 0 {
			sc := (b[i] - a[i]) / math.Max(a[i], b[i])
			sumSC += sc
		}
	}
	if numObservations > 0 {
		return sumSC / float64(numObservations)
	} else {
		return 0.0
	}
}

func assignToCluster(observations []Node, centroids []Node) ([]int, [][]float64) {
	numClusters := len(centroids)
	numObservations := len(observations)
	clusterAssignments := make([]int, numObservations)
	distances := make([][]float64, numObservations)
	for i := 0; i < numObservations; i++ {
		distances[i] = make([]float64, numClusters)
		for j := 0; j < numClusters; j++ {
			distances[i][j] = distance(observations[i], centroids[j])
		}
		minDist := distances[i][0]
		minIndex := 0
		for j := 1; j < numClusters; j++ {
			if distances[i][j] < minDist {
				minDist = distances[i][j]
				minIndex = j
			}
		}
		clusterAssignments[i] = minIndex
	}

	return clusterAssignments, distances
}

func computeAScores(observations []Node, clusterAssignments []int, distances [][]float64) []float64 {
	numObservations := len(observations)
	a := make([]float64, numObservations)

	for i := 0; i < numObservations; i++ {
		clusterIndex := clusterAssignments[i]

		var sumDist float64
		numSameCluster := 0
		for j := 0; j < numObservations; j++ {
			if clusterAssignments[j] == clusterIndex && i != j {
				sumDist += distances[i][clusterIndex]
				numSameCluster++
			}
		}
		if numSameCluster > 0 {
			a[i] = sumDist / float64(numSameCluster)
		}
	}
	return a
}

func computeBScores(observations []Node, centroids []Node, clusterAssignments []int, numClusters int) []float64 {
	numObservations := len(observations)
	distances := make([][]float64, numObservations)
	for i := 0; i < numObservations; i++ {
		distances[i] = make([]float64, numClusters)
		for j := 0; j < numClusters; j++ {
			distances[i][j] = distance(observations[i], centroids[j])
		}
	}

	b := make([]float64, numObservations)
	for i := 0; i < numObservations; i++ {
		clusterIndex := clusterAssignments[i]
		minAvgDist := math.Inf(1)
		for j := 0; j < numClusters; j++ {
			if j != clusterIndex {
				var sumDist float64
				numOtherCluster := 0
				for k := 0; k < numObservations; k++ {
					if clusterAssignments[k] == j {
						sumDist += distances[i][j]
						numOtherCluster++
					}
				}
				if numOtherCluster > 0 {
					avgDist := sumDist / float64(numOtherCluster)
					if avgDist < minAvgDist {
						minAvgDist = avgDist
					}
				}
			}
		}
		b[i] = minAvgDist
	}
	return b
}
