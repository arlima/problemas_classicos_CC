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
	ret.init(randDimensions)
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
	for _, cluster := range k.clusters {
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
		newCenter.init(means)
		cluster.centroid = newCenter
	}
}

func (k *kmeans) run(maxIterations int) []cluster {
	for i := 0; i < maxIterations; i++ {
		fmt.Println("Interaction: ", i)
		for j := range k.clusters {
			k.clusters[j].clearPoints()
		}
		k.assignCluster()
		//oldCentroids := k.centroids()
		k.generateCentroids()
		/*
			if oldCentroids == k.centroids() {
				fmt.Printf("Converged after %d iterations.", i)
				return k.clusters
			}*/
		for i, cluster := range k.clusters {
			fmt.Printf("Cluster %d: %v\n", i, cluster.points)
		}
	}
	return k.clusters
}

func main() {
	point1 := datapoint{}
	point1.init([]float64{2.0, 1.0, 1.0})
	point2 := datapoint{}
	point2.init([]float64{2.0, 2.0, 5.0})
	point3 := datapoint{}
	point3.init([]float64{3.0, 1.5, 2, 5})
	km := kmeans{}
	km.init(2, []datapoint{point1, point2, point3})
	result := km.run(100)
	for i, cluster := range result {
		fmt.Printf("Cluster %d: %v\n", i, cluster.points)
	}
}
