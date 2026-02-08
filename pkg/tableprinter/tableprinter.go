package tableprinter

import (
	"fmt"
	"reflect"
	"strings"
)

// NestedTableConfig defines the headers and fields for nested tables
type NestedTableConfig struct {
	Headers []string
	Fields  []string
}

// PrintTable is the main function to print a table for any struct slice with dynamic nested structures
func PrintTable(data interface{}, headers []string, fields []string, nestedConfig map[string]NestedTableConfig, indent int) {
	val := reflect.ValueOf(data)

	if val.Kind() != reflect.Slice {
		fmt.Println("Data is not a slice")
		return
	}

	// Print headers
	printHeaders(headers, indent)
	fmt.Println()

	// Print each struct in the slice
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		printRow(elem, fields, nestedConfig, indent)
	}
}

func printHeaders(headers []string, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print("  ")
	}
	for _, header := range headers {
		fmt.Printf("%-20s ", header)
	}
	fmt.Println()
}

func printRow(elem reflect.Value, fields []string, nestedConfig map[string]NestedTableConfig, indent int) {
	// Dereference pointer if necessary
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}

	for i := 0; i < indent; i++ {
		fmt.Print("  ")
	}
	for _, field := range fields {
		fieldVal := getNestedField(elem, field)
		fmt.Printf("%-20v ", fieldVal)
	}
	fmt.Println()

	// Check for nested structures based on the provided configuration
	for nestedField, config := range nestedConfig {
		nestedFieldValue := elem.FieldByName(nestedField)
		if nestedFieldValue.IsValid() && nestedFieldValue.Kind() == reflect.Slice {
			printNestedTable(nestedFieldValue.Interface(), config.Headers, config.Fields, nestedConfig, indent+1)
		}
	}
}

func printNestedTable(data interface{}, headers []string, fields []string, nestedConfig map[string]NestedTableConfig, indent int) {
	val := reflect.ValueOf(data)

	printHeaders(headers, indent)
	fmt.Println()

	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i)
		for j := 0; j < indent; j++ {
			fmt.Print("  ")
		}
		printNestedRow(elem, fields, nestedConfig, indent)
	}
}

func printNestedRow(elem reflect.Value, fields []string, nestedConfig map[string]NestedTableConfig, indent int) {
	// Dereference pointer if necessary
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}

	for _, field := range fields {
		fieldVal := getNestedField(elem, field)
		fmt.Printf("%-20v ", fieldVal)
	}
	fmt.Println()

	// Recursively handle further nested structures if any
	for nestedField, config := range nestedConfig {
		nestedFieldValue := elem.FieldByName(nestedField)
		if nestedFieldValue.IsValid() && nestedFieldValue.Kind() == reflect.Slice {
			printNestedTable(nestedFieldValue.Interface(), config.Headers, config.Fields, nestedConfig, indent+1)
		}
	}
}

func getNestedField(v reflect.Value, field string) interface{} {
	fields := splitFieldPath(field)
	for _, f := range fields {
		// Dereference pointer if necessary
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		v = v.FieldByName(f)
	}
	return v.Interface()
}

func splitFieldPath(path string) []string {
	return strings.Split(path, ".")
}
