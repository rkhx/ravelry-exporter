package gapi

import (
	"fmt"
	"google.golang.org/api/sheets/v4"
)

func getFreezeRowRequest(sheetID int64) *sheets.Request {
	return &sheets.Request{
		UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
			Properties: &sheets.SheetProperties{
				SheetId: sheetID,
				GridProperties: &sheets.GridProperties{
					FrozenRowCount: 1,
				},
			},
			Fields: "gridProperties.frozenRowCount",
		},
	}
}

func getAddSheetRequest(sheetTitle string) *sheets.Request {
	return &sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: sheetTitle,
			},
		},
	}
}

func getSetColumnWidthRequest(sheetID, minWidth int64) *sheets.Request {
	return &sheets.Request{
		UpdateDimensionProperties: &sheets.UpdateDimensionPropertiesRequest{
			Range: &sheets.DimensionRange{
				SheetId:    sheetID,
				Dimension:  "COLUMNS",
				StartIndex: 0,
				EndIndex:   26,
			},
			Properties: &sheets.DimensionProperties{
				PixelSize: minWidth,
			},
			Fields: "pixelSize",
		},
	}
}

func getAddHeaderRowRequest(sheetID int64, columns []interface{}) *sheets.Request {
	return &sheets.Request{
		UpdateCells: &sheets.UpdateCellsRequest{
			Rows: []*sheets.RowData{
				{Values: func() []*sheets.CellData {
					var cells []*sheets.CellData
					for _, col := range columns {
						colText := fmt.Sprintf("%v", col)
						cells = append(cells, &sheets.CellData{
							UserEnteredValue: &sheets.ExtendedValue{StringValue: &colText},
							UserEnteredFormat: &sheets.CellFormat{
								TextFormat: &sheets.TextFormat{Bold: true},
							},
						})
					}
					return cells
				}()},
			},
			Fields: "userEnteredValue,userEnteredFormat.textFormat.bold",
			Range: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    0,
				EndRowIndex:      1,
				StartColumnIndex: 0,
				EndColumnIndex:   int64(len(columns)),
			},
		},
	}
}
