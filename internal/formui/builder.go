package formui

import (
	"errors"
	"fyne.io/fyne/v2/widget"
	"reflect"
)

type BuiltForm struct {
	Form    *widget.Form
	Entries map[string]*widget.Entry
}

func BuildFormFromStruct(v any) *BuiltForm {
	t := reflect.TypeOf(v)
	entries := make(map[string]*widget.Entry)
	var formItems []*widget.FormItem

	var walkFields func(reflect.Type, string)
	walkFields = func(t reflect.Type, prefix string) {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldPath := field.Name
			if prefix != "" {
				fieldPath = prefix + "." + field.Name
			}

			if !field.IsExported() {
				continue
			}

			if field.Type.Kind() == reflect.Struct {
				walkFields(field.Type, fieldPath)
				continue
			}

			tag := ParseFyneFormTag(field.Tag.Get("formui"))
			label := tag["label"]
			placeholder := tag["placeholder"]

			var entry *widget.Entry
			if tag["type"] == "textarea" {
				entry = widget.NewMultiLineEntry()
			} else {
				entry = widget.NewEntry()
			}

			entry.SetPlaceHolder(placeholder)

			if tag["required"] == "true" {
				entry.Validator = func(s string) error {
					if s == "" {
						return errors.New("this field is required")
					}
					return nil
				}
			}

			entries[fieldPath] = entry
			item := widget.NewFormItem(label, entry)
			formItems = append(formItems, item)
		}
	}

	walkFields(t, "")
	return &BuiltForm{
		Form:    &widget.Form{Items: formItems},
		Entries: entries,
	}
}
