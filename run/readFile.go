package run

import (
	v "apiGO/structFile"

	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReadFileGet(file string) ([]v.Flower, []v.Car, []v.Furniture, error) {

	s, err := os.Open(file)
	if err != nil {
		return nil, nil, nil, err
	}
	defer s.Close()
	byteValue, _ := io.ReadAll(s)
	var data0 []v.Inventory
	var data1 []v.Flower
	var data2 []v.Car
	var data3 []v.Furniture
	err = json.Unmarshal(byteValue, &data0)
	if err != nil {
		fmt.Println(err)
		return nil, nil, nil, err
	}

	data1 = data0[0].Flowers
	fmt.Println(data1)
	data2 = data0[0].Cars
	fmt.Println(data2)
	data3 = data0[0].Furniture
	fmt.Println(data3)
	return data1, data2, data3, nil
}
