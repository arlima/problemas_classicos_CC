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
	ret.init(randDimensions, "", 0)
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
		newCenter.init(means, "", 0)
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

type album struct {
	name   string
	year   int
	length float64
	tracks int
}

func main() {
	rand.Seed(90)
	albums := []album{album{"Got to Be There", 1972, 35.45, 10}, album{"Ben", 1972, 31.31, 10},
		album{"Music & Me", 1973, 32.09, 10}, album{"Forever, Michael", 1975, 33.36, 10},
		album{"Off the Wall", 1979, 42.28, 10}, album{"Thriller", 1982, 42.19, 9},
		album{"Bad", 1987, 48.16, 10}, album{"Dangerous", 1991, 77.03, 14},
		album{"HIStory: Past, Present and Future, Book I", 1995, 148.58, 30}, album{"Invincible", 2001, 77.05, 16}}

	points := []datapoint{}
	for k := 0; k < len(albums); k++ {
		points = append(points, datapoint{})
		points[k].init([]float64{albums[k].length, float64(albums[k].tracks)}, albums[k].name, albums[k].year)
	}
	km := kmeans{}
	km.init(2, points)
	result := km.run(42)
	fmt.Printf("Cluster, Name, Year, Length, Tracks\n")
	for i, cluster := range result {
		for _, p := range cluster.points {
			fmt.Printf("%d, %v, %v, %v, %v \n", i, p.name, p.year, p.originals[0], p.originals[1])
		}
	}

	for _, cluster := range result {
		fmt.Println(cluster.centroid.dimensions)
	}
}
