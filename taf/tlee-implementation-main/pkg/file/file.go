package file

import (
	"bufio"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/dtos"
	"connect.informatik.uni-ulm.de/coordination/tlee-implementation/pkg/utils"
	"encoding/csv"
	"errors"
	"github.com/gocarina/gocsv"
	"io"
	"os"
	"strings"
)

/*
This method takes as a paramenter a string that identifies the name of the file.
Through its functions it allows us to open the file and read.
IMPORTANT ----> [In the first line of the readFromCsv you have to specify the path of where the file is located
if it is not placed in the conf/csv/ project folder.]
*/
func GetOpinion(getOpinionMode string) ([]dtos.OpinionDTOFlatten, error) {
	if strings.Contains(getOpinionMode, ".csv") {
		return readFromCsv(getOpinionMode), nil
	}
	return nil, errors.New("this metod actually work just with csv")
}

/*
IMPORTANT ----> [In the first line of the readFromCsv you have to specify the path of where the file is located
if it is not placed in the conf/csv/ project folder.]
*/
func readFromCsv(nameFile string) []dtos.OpinionDTOFlatten {
	opinionFile, err := os.OpenFile(nameFile, os.O_RDWR, os.ModePerm)
	utils.Must(err)
	setFileReader(opinionFile)
	defer opinionFile.Close()
	var opinions []dtos.OpinionDTOFlatten
	utils.Must(gocsv.UnmarshalFile(opinionFile, &opinions))
	return opinions
}

func setFileReader(file io.Reader) {
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(bufio.NewReader(file))
		r.Comma = ','
		return r
	})
}
