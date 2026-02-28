This method takes as a paramenter a string that identifies the name of the file.
Through its functions it allows us to open the file and read.

### Parameters

- `getOpinionMode`: A string that identifies the name of the file.

### Important

- In the first line of the `readFromCsv` you have to specify the path of where the file is located if it is not placed in the conf/csv/ project folder.

**Snippet**

```go
func GetOpinion(getOpinionMode string) ([]dtos.OpinionDTO, error) {
 if strings.Contains(getOpinionMode, ".csv") {
  return readFromCsv(getOpinionMode), nil
 }
 return nil, errors.New("this metod actually work just with csv")
}

func readFromCsv(nameFile string) []dtos.OpinionDTO {
 opinionFile, err := os.OpenFile("conf/csv/"+nameFile, os.O_RDWR, os.ModePerm)
 utils.Must(err)
 setFileReader(opinionFile)
 defer opinionFile.Close()
 var opinions []dtos.OpinionDTO
 utils.Must(gocsv.UnmarshalFile(opinionFile, &opinions))

 return opinions
}

// Function to set the delimiter on csv file
func setFileReader(file io.Reader) {
 gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
  r := csv.NewReader(bufio.NewReader(file))
  r.Comma = ','
  return r
 })
}
