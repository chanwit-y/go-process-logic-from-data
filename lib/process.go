package lib

import (
	"encoding/json"
	"fmt"
)

type FlowType int

const (
	Assign    FlowType = 0
	Operator           = 1
	IF                 = 2
	Loop               = 3
	Query              = 4
	Transfrom          = 5
	Printf             = 6
)

type DataType int

const (
	Int       DataType = 0
	String             = 1
	Boolean            = 2
	JSONArray          = 3
)

type AssignType struct {
	Name     string
	DataType DataType
	Default  string // array object and object use json string
}

type IfType struct{}

type LoopType struct {
	Key          string
	IsFromString bool // set true when loop 2 step up
	Flow         []Flow
}

type PrintfType struct {
	Key     string
	Message string
}

type Flow struct {
	Type   FlowType
	Assign AssignType
	If     IfType
	Loop   LoopType
	Printf PrintfType
}

func IfProcess(flow []Flow, variables *map[string]string) {}

func ForEachObjProcess(data []map[string]any, flow []Flow, variables map[string]any) map[string]any {

	for _, d := range data {
		// key := fmt.Sprintf("foreachobjprocess%v", i)
		// variables[key] = d
		for key, element := range d {
			variables[string(key)] = element
		}
		for key, element := range variables {
			fmt.Printf("key: %+v \n", key)
			fmt.Printf("element: %+v \n", element)
		}
		variables = RunProcess(flow, variables)
	}

	return variables
}

func RunProcess(flow []Flow, variables map[string]any) map[string]any {
	for _, f := range flow {
		// fmt.Printf("type: %v \n", f.Type)
		switch f.Type {
		case 0:
			{
				variables[f.Assign.Name] = f.Assign.Default
				break
			}
		case 3:
			{
				arr := []map[string]any{}
				if f.Loop.IsFromString {
					json.Unmarshal([]byte(variables[f.Loop.Key].(string)), &arr)
				} else {

					j, _ := json.Marshal(variables[f.Loop.Key])
					jsonStr := string(j)
					json.Unmarshal([]byte(jsonStr), &arr)
				}
				variables = ForEachObjProcess(arr, f.Loop.Flow, variables)
			}
		case 6:
			{
				fmt.Printf(f.Printf.Message, variables[f.Printf.Key])
				break
			}
		}
	}
	return variables
}

func Process() {
	data := []map[string]any{
		{
			"Code": "1",
			"Name": "test 1",
			"Json2": []map[string]any{
				{
					"A": "data a 1",
					"B": "data b 1",
				},
				{
					"A": "data a 2",
					"B": "data b 2",
				},
			},
		},
		{
			"Code": "2",
			"Name": "test 2",
		},
	}

	json, _ := json.Marshal(data)
	jsonStr := string(json)

	flows := []Flow{
		{
			Type: 0,
			Assign: AssignType{
				Name:     "json1",
				DataType: 1,
				Default:  jsonStr,
			},
		},
		{
			Type: 0,
			Assign: AssignType{
				Name:     "str1",
				DataType: 1,
				Default:  "test 1",
			},
		},
		{
			Type: 3,
			Loop: LoopType{
				Key:          "json1",
				IsFromString: true,
				Flow: []Flow{
					{
						Type: 6,
						Printf: PrintfType{
							Message: "Code: %v ",
							Key:     "Code",
						},
					},
					{
						Type: 6,
						Printf: PrintfType{
							Message: "Name: %v \n",
							Key:     "Name",
						},
					},
					{
						Type: 3,
						Loop: LoopType{
							Key:          "Json2",
							IsFromString: false,
							Flow: []Flow{
								{
									Type: 6,
									Printf: PrintfType{
										Message: "A: %v \n",
										Key:     "A",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Type: 6,
			Printf: PrintfType{
				Key:     "str1",
				Message: "log str1: %v \n",
			},
		},
	}

	variables := map[string]any{}
	variables = RunProcess(flows, variables)

	fmt.Printf("Hi \n")
}
