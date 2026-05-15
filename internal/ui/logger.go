package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/rs/zerolog"
	"image/color"
	"os"
	"time"
)

type ColoredTextSegment struct {
	*widget.TextSegment
	Color color.Color
}

func (t *ColoredTextSegment) Visual() fyne.CanvasObject {
	obj := canvas.NewText(t.Text, t.Color)

	t.Update(obj)
	return obj
}

func (t *ColoredTextSegment) Update(o fyne.CanvasObject) {
	obj := o.(*canvas.Text)
	obj.Text = t.Text
	obj.Color = t.Color
	obj.Alignment = t.Style.Alignment
	obj.TextStyle = t.Style.TextStyle
	obj.TextSize = 15.0
	obj.Refresh()
}

var levelColors = map[zerolog.Level]color.NRGBA{
	zerolog.DebugLevel: {0, 200, 0, 255},              // Green
	zerolog.InfoLevel:  {R: 0, G: 0, B: 0xff, A: 255}, // Blue
	zerolog.WarnLevel:  {255, 165, 0, 255},            // Orange
	zerolog.ErrorLevel: {255, 0, 0, 255},              // Red
	zerolog.FatalLevel: {128, 0, 128, 255},            // Purple
	zerolog.PanicLevel: {255, 0, 255, 255},            // Magenta
}

func getColor(level zerolog.Level) color.Color {
	logColor, exists := levelColors[level]
	if !exists {
		logColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	}

	return logColor
}

func GetUIReporter(receiver *widget.RichText) zerolog.ConsoleWriter {
	wr := zerolog.ConsoleWriter{
		PartsOrder: []string{
			zerolog.LevelFieldName,
			zerolog.TimestampFieldName,
			zerolog.CallerFieldName,
			zerolog.MessageFieldName,
		},
		TimeFormat: time.TimeOnly,
		Out:        os.Stdout,
		FormatMessage: func(i interface{}) string {
			receiver.Segments = append(receiver.Segments, &widget.TextSegment{
				Text:  i.(string) + "\n",
				Style: widget.RichTextStyleInline,
			})
			receiver.Refresh()
			return i.(string)
		},
		FormatLevel: func(i interface{}) string {
			lvl, _ := zerolog.ParseLevel(i.(string))
			receiver.Segments = append(receiver.Segments, &ColoredTextSegment{
				TextSegment: &widget.TextSegment{
					Text:  "[" + i.(string) + "]",
					Style: widget.RichTextStyleInline,
				},
				Color: getColor(lvl),
			})
			receiver.Refresh()
			return i.(string)
		},
	}

	wr.FormatTimestamp = func(i interface{}) string {
		ts, _ := time.ParseInLocation(zerolog.TimeFieldFormat, i.(string), time.Local)
		v := ts.In(time.Local).Format(wr.TimeFormat)

		receiver.Segments = append(receiver.Segments, &widget.TextSegment{
			Text:  "[" + v + "] ",
			Style: widget.RichTextStyleInline,
		})
		receiver.Refresh()
		return v
	}

	return wr
}
