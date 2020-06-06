package main

import (
	"log"
	"math"
	"math/rand"

	"github.com/gonum/floats"
	"github.com/gonum/stat"
	"golang.org/x/exp/errors/fmt"
)

type cluster struct {
	points   []datapoint
	centroid datapoint
}

func (c *cluster) clearPoints() {
	c.points = nil
}

type kmeans struct {
	points   []datapoint
	clusters []cluster
}

func (k *kmeans) init(ka int, points []datapoint) {
	if ka < 1 {
		log.Fatal("k must be >= 1")
	}
	k.points = points
	k.zscoreNormalize()
	for i := 0; i < ka; i++ {
		randPoint := k.randomPoint()
		c := cluster{[]datapoint{}, randPoint}
		k.clusters = append(k.clusters, c)
	}
}

func (k *kmeans) centroids() []datapoint {
	ret := []datapoint{}
	for _, v := range k.clusters {
		ret = append(ret, v.centroid)
	}
	return ret
}

func (k *kmeans) dimensionSlices(dimension int) []float64 {
	ret := []float64{}
	for _, x := range k.points {
		ret = append(ret, x.dimensions[dimension])
	}
	return ret
}

func zScores(original []float64) []float64 {
	ret := make([]float64, len(original))
	avg, std := stat.MeanStdDev(original, nil)
	if std == 0.0 {
		return ret
	}
	for i, x := range original {
		ret[i] = (x - avg) / std
	}
	return ret
}

func (k *kmeans) zscoreNormalize() {
	zscored := [][]float64{}
	for points := 0; points < len(k.points); points++ {
		zscored = append(zscored, []float64{})
	}
	for dimension := 0; dimension < k.points[0].numDimensions(); dimension++ {
		dimensionSlice := k.dimensionSlices(dimension)
		for i, zscore := range zScores(dimensionSlice) {
			zscored[i] = append(zscored[i], zscore)
		}
	}
	for points := 0; points < len(k.points); points++ {
		k.points[points].dimensions = zscored[points]
	}
}

func (k *kmeans) randomPoint() datapoint {
	randDimensions := []float64{}
	for dimension := 0; dimension < k.points[0].numDimensions(); dimension++ {
		values := k.dimensionSlices(dimension)
		randValue := floats.Min(values) + rand.Float64()*(floats.Max(values)-floats.Min(values))
		randDimensions = append(randDimensions, randValue)
	}
	ret := datapoint{}
	ret.init(randDimensions, "")
	return ret
}

func (k *kmeans) assignCluster() {
	for _, point := range k.points {
		min := math.MaxFloat64
		minIndex := 0
		for k, cluster := range k.clusters {
			dist := point.distance(cluster.centroid)
			if dist < min {
				min = dist
				minIndex = k
			}
		}
		k.clusters[minIndex].points = append(k.clusters[minIndex].points, point)
	}
}

func (k *kmeans) generateCentroids() {
	for c, cluster := range k.clusters {
		if len(cluster.points) == 0 {
			continue
		}
		means := []float64{}
		for dimension := 0; dimension < cluster.points[0].numDimensions(); dimension++ {
			dimensionSlice := []float64{}
			for _, point := range cluster.points {
				dimensionSlice = append(dimensionSlice, point.dimensions[dimension])
			}
			means = append(means, stat.Mean(dimensionSlice, nil))
		}
		newCenter := datapoint{}
		newCenter.init(means, "")
		k.clusters[c].centroid = newCenter
	}
}

func isEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (k *kmeans) run(maxIterations int) []cluster {
	for i := 0; i < maxIterations; i++ {
		fmt.Println("Interaction: ", i+1)
		for j := range k.clusters {
			k.clusters[j].clearPoints()
		}
		k.assignCluster()
		oldCentroids := k.centroids()
		k.generateCentroids()
		count := 0
		for i, c := range k.centroids() {
			if !isEqual(c.dimensions, oldCentroids[i].dimensions) {
				break
			}
			count++
		}
		if count == len(oldCentroids) {
			fmt.Printf("Converged after %d iterations.\n", i+1)
			return k.clusters
		}
	}
	return k.clusters
}

type governor struct {
	long  float64
	age   int
	state string
}

func main() {
	rand.Seed(211)
	governors := []governor{governor{-86.79113, 72, "Alabama"}, governor{-152.404419, 66, "Alaska"},
		governor{-111.431221, 53, "Arizona"}, governor{-92.373123, 66, "Arkansas"},
		governor{-119.681564, 79, "California"}, governor{-105.311104, 65, "Colorado"},
		governor{-72.755371, 61, "Connecticut"}, governor{-75.507141, 61, "Delaware"},
		governor{-81.686783, 64, "Florida"}, governor{-83.643074, 74, "Georgia"},
		governor{-157.498337, 60, "Hawaii"}, governor{-114.478828, 75, "Idaho"},
		governor{-88.986137, 60, "Illinois"}, governor{-86.258278, 49, "Indiana"},
		governor{-93.210526, 57, "Iowa"}, governor{-96.726486, 60, "Kansas"},
		governor{-84.670067, 50, "Kentucky"}, governor{-91.867805, 50, "Louisiana"},
		governor{-69.381927, 68, "Maine"}, governor{-76.802101, 61, "Maryland"},
		governor{-71.530106, 60, "Massachusetts"}, governor{-84.536095, 58, "Michigan"},
		governor{-93.900192, 70, "Minnesota"}, governor{-89.678696, 62, "Mississippi"},
		governor{-92.288368, 43, "Missouri"}, governor{-110.454353, 51, "Montana"},
		governor{-98.268082, 52, "Nebraska"}, governor{-117.055374, 53, "Nevada"},
		governor{-71.563896, 42, "New Hampshire"}, governor{-74.521011, 54, "New Jersey"},
		governor{-106.248482, 57, "New Mexico"}, governor{-74.948051, 59, "New York"},
		governor{-79.806419, 60, "North Carolina"}, governor{-99.784012, 60, "North Dakota"},
		governor{-82.764915, 65, "Ohio"}, governor{-96.928917, 62, "Oklahoma"},
		governor{-122.070938, 56, "Oregon"}, governor{-77.209755, 68, "Pennsylvania"},
		governor{-71.51178, 46, "Rhode Island"}, governor{-80.945007, 70, "South Carolina"},
		governor{-99.438828, 64, "South Dakota"}, governor{-86.692345, 58, "Tennessee"},
		governor{-97.563461, 59, "Texas"}, governor{-111.862434, 70, "Utah"},
		governor{-72.710686, 58, "Vermont"}, governor{-78.169968, 60, "Virginia"},
		governor{-121.490494, 66, "Washington"}, governor{-80.954453, 66, "West Virginia"},
		governor{-89.616508, 49, "Wisconsin"}, governor{-107.30249, 55, "Wyoming"}}

	points := []datapoint{}
	for k := 0; k < len(governors); k++ {
		points = append(points, datapoint{})
		points[k].init([]float64{governors[k].long, float64(governors[k].age)}, governors[k].state)
	}
	km := kmeans{}
	km.init(4, points)
	result := km.run(42)
	fmt.Printf("Cluster, State, coordinates\n")
	for i, cluster := range result {
		for _, p := range cluster.points {
			fmt.Printf("%d, %v, %v, %v \n", i, p.label, p.dimensions[0], p.dimensions[1])
		}
	}

	for _, cluster := range result {
		fmt.Println(cluster.centroid.dimensions)
	}
}
