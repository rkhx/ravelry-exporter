package gapi

import (
	"context"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	"time"
)

func createSpreadsheet(ctx context.Context, driveSrv *drive.Service, sheetsSrv *sheets.Service, title string) (string, error) {
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: title,
		},
	}
	resp, err := sheetsSrv.Spreadsheets.Create(spreadsheet).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("unable to create spreadsheet: %w", err)
	}
	spreadsheetID := resp.SpreadsheetId

	permission := &drive.Permission{
		Type:         "user",
		Role:         "writer",
		EmailAddress: "r.khaniukov@gmail.com",
	}
	_, err = driveSrv.Permissions.Create(resp.SpreadsheetId, permission).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("unable to share spreadsheet: %w", err)
	}

	return spreadsheetID, nil
}

func Spreadsheet(ctx context.Context, serviceAccountFileContent, spreadsheetID string, data []RowData, columns []interface{}) error {
	sheetsSrv, driveSrv, err := newGoogleServices(serviceAccountFileContent, ctx)
	if err != nil {
		return err
	}

	if spreadsheetID == "" {
		spreadsheetID, err = createSpreadsheet(ctx, driveSrv, sheetsSrv, "Shared Spreadsheet")
		if err != nil {
			return err
		}
	}

	sheetTitle := fmt.Sprintf("Sheet_%d", time.Now().Unix())
	addSheetRequest := getAddSheetRequest(sheetTitle)
	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{addSheetRequest},
	}

	batchUpdateResponse, err := sheetsSrv.Spreadsheets.BatchUpdate(spreadsheetID, batchUpdateRequest).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to add new sheet '%s' to spreadsheet %s: %w", sheetTitle, spreadsheetID, err)
	}

	sheetID := batchUpdateResponse.Replies[0].AddSheet.Properties.SheetId

	requests := []*sheets.Request{
		getAddHeaderRowRequest(sheetID, columns),
		getFreezeRowRequest(sheetID),
	}

	startRow := 1
	for i, row := range data {
		rowIndex := startRow + i

		rowValues := func() []*sheets.CellData {
			var cells []*sheets.CellData
			for _, col := range columns {
				value := row[col.(string)]
				switch v := value.(type) {
				case []Link:
					cells = append(cells, createRichTextCell(v))
				case Link:
					cells = append(cells, createRichTextCell([]Link{v}))
				default:
					str := fmt.Sprintf("%v", v)
					cells = append(cells, &sheets.CellData{
						UserEnteredValue: &sheets.ExtendedValue{StringValue: &str},
					})
				}
			}
			return cells
		}()

		requests = append(requests, &sheets.Request{
			UpdateCells: &sheets.UpdateCellsRequest{
				Rows:   []*sheets.RowData{{Values: rowValues}},
				Fields: "userEnteredValue,textFormatRuns",
				Range: &sheets.GridRange{
					SheetId:          sheetID,
					StartRowIndex:    int64(rowIndex),
					EndRowIndex:      int64(rowIndex + 1),
					StartColumnIndex: 0,
					EndColumnIndex:   int64(len(columns)),
				},
			},
		})
	}

	requests = append(requests, getSetColumnWidthRequest(sheetID, int64(250)))
	if len(requests) > 0 {
		_, err := sheetsSrv.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
			Requests: requests,
		}).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("batch update failed for spreadsheet %s: %w", spreadsheetID, err)
		}
	}

	return nil
}
