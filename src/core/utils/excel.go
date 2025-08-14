package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ExcelDropdownConfig represents the configuration for adding a dropdown to Excel
type ExcelDropdownConfig struct {
	SheetName    string   // Name of the sheet to add dropdown to
	CellRange    string   // Cell range like "D2:D10"
	Options      []string // List of options for the dropdown
	AllowBlank   bool     // Whether to allow blank values
	ShowError    bool     // Whether to show error message for invalid input
	ErrorTitle   string   // Title for error message
	ErrorMessage string   // Error message text
	InputTitle   string   // Title for input message
	InputMessage string   // Input message text
}

// AddDropdownToExcel adds a dropdown selector to an Excel file
// filePath: path to the Excel file
// config: configuration for the dropdown
func AddDropdownToExcel(filePath string, config ExcelDropdownConfig) error {
	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	// Validate cell range format
	if err := validateCellRange(config.CellRange); err != nil {
		return fmt.Errorf("invalid cell range: %w", err)
	}

	// Create data validation for dropdown
	dv := excelize.NewDataValidation(true)
	dv.Sqref = config.CellRange
	dv.Type = "list"
	dv.Formula1 = strings.Join(config.Options, ",")
	dv.AllowBlank = config.AllowBlank
	dv.ShowErrorMessage = config.ShowError
	dv.ShowInputMessage = true

	// Add data validation to the sheet
	if err := f.AddDataValidation(config.SheetName, dv); err != nil {
		return fmt.Errorf("failed to add data validation: %w", err)
	}

	// Save the file
	if err := f.Save(); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}

// AddDropdownToExcelWithDefaultConfig adds a dropdown with default configuration
func AddDropdownToExcelWithDefaultConfig(filePath, sheetName, cellRange string, options []string) error {
	config := ExcelDropdownConfig{
		SheetName:    sheetName,
		CellRange:    cellRange,
		Options:      options,
		AllowBlank:   true,
		ShowError:    true,
		ErrorTitle:   "Invalid Input",
		ErrorMessage: "Please select a value from the dropdown list.",
		InputTitle:   "Select Option",
		InputMessage: "Please select an option from the dropdown.",
	}

	return AddDropdownToExcel(filePath, config)
}

// validateCellRange validates the cell range format (e.g., "D2:D10")
func validateCellRange(cellRange string) error {
	// Regex pattern for cell range format like "A1:B10" or "D2:D10"
	pattern := `^[A-Z]+\d+:[A-Z]+\d+$`
	matched, err := regexp.MatchString(pattern, cellRange)
	if err != nil {
		return fmt.Errorf("regex validation error: %w", err)
	}
	if !matched {
		return fmt.Errorf("invalid cell range format. Expected format like 'D2:D10', got: %s", cellRange)
	}
	return nil
}

// parseCellRange parses a cell range string and returns start and end cells
func parseCellRange(cellRange string) (startCell, endCell string, err error) {
	parts := strings.Split(cellRange, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid cell range format")
	}
	return parts[0], parts[1], nil
}

// GetCellAddress converts row and column numbers to Excel cell address
func GetCellAddress(row, col int) string {
	cellName, _ := excelize.CoordinatesToCellName(col, row)
	return cellName
}

// GetCellCoordinates converts Excel cell address to row and column numbers
func GetCellCoordinates(cell string) (row, col int, err error) {
	return excelize.CellNameToCoordinates(cell)
}

// CreateExcelWithDropdown creates a new Excel file with dropdown selectors
func CreateExcelWithDropdown(filePath string, dropdowns []ExcelDropdownConfig) error {
	// Create a new Excel file
	f := excelize.NewFile()
	defer f.Close()

	// Add dropdowns to the file
	for _, dropdown := range dropdowns {
		if err := AddDropdownToExcelFile(f, dropdown); err != nil {
			return fmt.Errorf("failed to add dropdown: %w", err)
		}
	}

	// Save the file
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}

// AddDropdownToExcelFile adds a dropdown to an existing Excel file object
func AddDropdownToExcelFile(f *excelize.File, config ExcelDropdownConfig) error {
	// Validate cell range format
	if err := validateCellRange(config.CellRange); err != nil {
		return fmt.Errorf("invalid cell range: %w", err)
	}

	// Create data validation for dropdown
	dv := excelize.NewDataValidation(true)
	dv.Sqref = config.CellRange
	dv.Type = "list"
	dv.Formula1 = strings.Join(config.Options, ",")
	dv.AllowBlank = config.AllowBlank
	dv.ShowErrorMessage = config.ShowError
	dv.ShowInputMessage = true

	// Add data validation to the sheet
	if err := f.AddDataValidation(config.SheetName, dv); err != nil {
		return fmt.Errorf("failed to add data validation: %w", err)
	}

	return nil
}

// AddMultipleDropdowns adds multiple dropdowns to an Excel file
func AddMultipleDropdowns(filePath string, dropdowns []ExcelDropdownConfig) error {
	// Open the Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	// Add each dropdown
	for i, dropdown := range dropdowns {
		if err := AddDropdownToExcelFile(f, dropdown); err != nil {
			return fmt.Errorf("failed to add dropdown %d: %w", i+1, err)
		}
	}

	// Save the file
	if err := f.Save(); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}

// Example usage function
func ExampleAddDropdown() error {
	// Example: Add dropdown to D2:D10 with options
	config := ExcelDropdownConfig{
		SheetName:    "Sheet1",
		CellRange:    "D2:D10",
		Options:      []string{"Option 1", "Option 2", "Option 3", "Option 4"},
		AllowBlank:   true,
		ShowError:    true,
		ErrorTitle:   "Invalid Selection",
		ErrorMessage: "Please select a valid option from the dropdown.",
		InputTitle:   "Select Option",
		InputMessage: "Choose an option from the dropdown list.",
	}

	return AddDropdownToExcel("example.xlsx", config)
}
