package office

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// ReadExcel function reads the excel file and return 2 dimension slice
func ReadExcel(path string, sheetName string) ([][]string, error) {
	if sheetName == "" {
		sheetName = "Sheet1"
	}
	var r [][]string
	fd, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return r, err
	}
	rows := fd.GetRows(sheetName)
	r = rows
	return r, err
}

// ReadExcel function reads the excel file and return 2 dimension slice
func ReadExcelFromReader(reader io.Reader, sheetName string) ([][]string, error) {
	if sheetName == "" {
		sheetName = "Sheet1"
	}
	var r [][]string
	fd, err := excelize.OpenReader(reader)
	if err != nil {
		fmt.Println(err.Error())
		return r, err
	}
	rows := fd.GetRows(sheetName)
	r = rows
	return r, err
}

// writeData writes the data to an *excel.File type variable
func writeData(dataIn interface{}, col []string) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := "Sheet1"
	f.NewSheet(sheetName)

	// Get data
	var data [][]interface{}
	v := reflect.ValueOf(dataIn)
	if v.Kind() != reflect.Slice {
		err := errors.New("Data to write should be 2 dim slice.")
		return f, err
	}
	for i := 0; i < v.Len(); i++ {
		itemV := v.Index(i)
		if itemV.Kind() != reflect.Slice {
			err := errors.New("Data to write should be 2 dim slice.")
			return f, err
		}
		var innerSli []interface{}
		for j := 0; j < itemV.Len(); j++ {
			innerSli = append(innerSli, itemV.Index(j))
		}
		data = append(data, innerSli)
	}

	// Write data to excelize.File
	f.SetSheetRow(sheetName, "A1", &col)
	for i, item := range data {
		f.SetSheetRow(sheetName, "A"+strconv.Itoa(i+2), &item)
	}

	return f, nil
}

// WriteExcel writes the data to an xlsx file
func WriteExcel(p string, dataIn interface{}, col []string) error {
	f, err := writeData(dataIn, col)
	if err != nil {
		return err
	}

	// Write to file
	err = f.SaveAs(p)
	if err != nil {
		return err
	}

	return nil
}

// WriteExcelToWriter writes the data to an io.Writer
func WriteExcelToWriter(w io.Writer, dataIn interface{}, col []string) error {
	f, err := writeData(dataIn, col)
	if err != nil {
		return err
	}

	// Write to an io.Writer variable
	err = f.Write(w)
	if err != nil {
		return err
	}

	return nil
}
