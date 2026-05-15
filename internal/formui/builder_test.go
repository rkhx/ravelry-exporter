package formui

import (
	"fyne.io/fyne/v2/test"
	"testing"

	"fyne.io/fyne/v2/widget"
)

type TestFormStruct struct {
	Username string `fyneform:"label=Username,placeholder=Enter your name,required"`
	Email    string `fyneform:"label=Email,placeholder=user@example.com"`
}

func TestBuildFormFromStruct(t *testing.T) {
	test.NewApp()
	defer test.NewApp().Quit()
	
	form := BuildFormFromStruct(TestFormStruct{})

	if form == nil {
		t.Fatal("form should not be nil")
	}

	if len(form.Items) != 2 {
		t.Fatalf("expected 2 form items, got %d", len(form.Items))
	}

	item1 := form.Items[0]
	item2 := form.Items[1]

	// Check labels
	if item1.Text != "Username" {
		t.Errorf("expected label 'Username', got '%s'", item1.Text)
	}
	if item2.Text != "Email" {
		t.Errorf("expected label 'Email', got '%s'", item2.Text)
	}

	// Check placeholders
	entry1, ok1 := item1.Widget.(*widget.Entry)
	entry2, ok2 := item2.Widget.(*widget.Entry)

	if !ok1 || !ok2 {
		t.Fatal("form items must contain widget.Entry")
	}

	if entry1.PlaceHolder != "Enter your name" {
		t.Errorf("expected placeholder for Username to be 'Enter your name', got '%s'", entry1.PlaceHolder)
	}
	if entry2.PlaceHolder != "user@example.com" {
		t.Errorf("expected placeholder for Email to be 'user@example.com', got '%s'", entry2.PlaceHolder)
	}

	// Check validator on required field
	entry1.SetText("") // simulate empty input
	err := entry1.Validate()
	if err == nil {
		t.Error("expected validation error for required Username field, got nil")
	}
}
