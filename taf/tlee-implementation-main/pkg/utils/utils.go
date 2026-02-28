package utils

import (
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/config"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/dtos"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
)

func Must(err error) {
	if err != nil {
		Should(err)
		panic(err)
	}
}

func CartesianProduct(sources, targets []string) [][]string {
	var product [][]string

	for _, source := range sources {
		for _, target := range targets {
			arr := []string{}
			arr = append(arr, source)
			arr = append(arr, target)
			product = append(product, arr)
		}
	}
	return product
}

func ListToTuples(aList []string) [][]string {
	n := make([][]string, 0)
	for i := 0; i < len(aList)-1; i++ {
		n = append(n, []string{aList[i], aList[i+1]})
	}
	return n
}

func Should(err error) {
	if err != nil {
		log.Error().Err(err).Send()
	}
}

func SortByLength(arrays [][]string) {
	sort.Slice(arrays, func(i, j int) bool {
		return len(arrays[i]) < len(arrays[j])
	})
}

func SlicePop[T any](s []T, i int) ([]T, T) {
	elem := s[i]
	s = append(s[:i], s[i+1:]...)
	return s, elem
}

func SortValue(myMap map[string]bool) map[string]bool {
	var data []struct {
		Key   string
		Value bool
	}

	for k, v := range myMap {
		data = append(data, struct {
			Key   string
			Value bool
		}{k, v})
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Key < data[j].Key
	})

	sortedMap := make(map[string]bool)
	for _, item := range data {
		sortedMap[item.Key] = item.Value
	}

	return sortedMap
}

func CreateBenchmarkFile(time float64, numEdgeInput int64, numEdgeOutput int64) {

	var fileName string = "banchmark.json"
	var benchMark []dtos.BenchMark

	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {

		benchMark = make([]dtos.BenchMark, 0)

	} else {

		err = json.Unmarshal(fileContent, &benchMark)
		Must(err)

	}

	benchMark = append(benchMark, dtos.BenchMark{Time: time, NumEdge: numEdgeInput, CompressionRate: float64(numEdgeOutput) / float64(numEdgeInput)})

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
	Must(err)

	defer file.Close()

	benchmarkJSON, err := json.Marshal(benchMark)
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString(fmt.Sprintf("%s\n\n", benchmarkJSON))
	Must(err)
}

func CreateExpDot(a1 *dtos.DataChildDTO, a2 *dtos.DataChildExpressionDTO, txt string) {
	var fileName string
	if a1 == nil {
		dotOutput, _ := a2.GenerateDot(0)
		dotGraph := "digraph G {\n" + dotOutput + "}\n"
		fileName = config.OutputPath + "/debug/tree_" + txt + ".dot"
		err := os.WriteFile(fileName, []byte(dotGraph), 0644)
		if err != nil {
			config.Logger.Error(fmt.Sprintf("Error writing to file: %s\n", err))

		}
	} else {
		dotOutput, _ := a1.GenerateDotD(0)
		dotGraph := "digraph G {\n" + dotOutput + "}\n"
		fileName = config.OutputPath + "/debug/tree_" + txt + ".dot"
		err := os.WriteFile(fileName, []byte(dotGraph), 0644)
		if err != nil {
			config.Logger.Error(fmt.Sprintf("Error writing to file: %s", err))
		}

	}
	outputImage := config.OutputPath + "/debug/tree_" + txt + ".png"
	cmd := exec.Command("dot", "-Tpng", fileName, "-o", outputImage)
	if err := cmd.Run(); err != nil {

		config.Logger.Error(fmt.Sprintf("Error running dot command %s", err))
	}

}
