package gapi

import "google.golang.org/api/sheets/v4"

func createRichTextCell(links []Link) *sheets.CellData {
	var text string
	var textRuns []*sheets.TextFormatRun
	offset := int64(0)

	for _, link := range links {
		if text != "" {
			text += "\n"
			offset += int64(len("\n"))
		}
		text += link.Text
		textRuns = append(textRuns, &sheets.TextFormatRun{
			StartIndex: offset,
			Format: &sheets.TextFormat{
				Link: &sheets.Link{Uri: link.URL},
			},
		})
		offset += int64(len(link.Text))

		if link.Extras != "" {
			textRuns = append(textRuns, &sheets.TextFormatRun{
				StartIndex: offset,
				Format:     &sheets.TextFormat{},
			})

			t := " (" + link.Extras + ")"
			text += t
			offset += int64(len(t))
		}
	}

	return &sheets.CellData{
		UserEnteredValue: &sheets.ExtendedValue{StringValue: &text},
		TextFormatRuns:   textRuns,
	}
}
