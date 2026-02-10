package desktop

import (
	"image/color"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func TestNewCustomTheme(t *testing.T) {
	customTheme := &CustomTheme{}
	assert.NotNil(t, customTheme)
}

func TestCustomTheme_Color(t *testing.T) {
	customTheme := &CustomTheme{}

	tests := []struct {
		name      string
		colorName fyne.ThemeColorName
		variant   fyne.ThemeVariant
		expected  color.Color
	}{
		{
			name:      "Primary color",
			colorName: theme.ColorNamePrimary,
			variant:   theme.VariantLight,
			expected:  color.RGBA{R: 99, G: 102, B: 241, A: 255},
		},
		{
			name:      "Background dark",
			colorName: theme.ColorNameBackground,
			variant:   theme.VariantDark,
			expected:  color.RGBA{R: 17, G: 24, B: 39, A: 255},
		},
		{
			name:      "Background light",
			colorName: theme.ColorNameBackground,
			variant:   theme.VariantLight,
			expected:  color.RGBA{R: 248, G: 250, B: 252, A: 255},
		},
		{
			name:      "Foreground dark",
			colorName: theme.ColorNameForeground,
			variant:   theme.VariantDark,
			expected:  color.RGBA{R: 248, G: 250, B: 252, A: 255},
		},
		{
			name:      "Foreground light",
			colorName: theme.ColorNameForeground,
			variant:   theme.VariantLight,
			expected:  color.RGBA{R: 15, G: 23, B: 42, A: 255},
		},
		{
			name:      "Hover color",
			colorName: theme.ColorNameHover,
			variant:   theme.VariantLight,
			expected:  color.RGBA{R: 99, G: 102, B: 241, A: 50},
		},
		{
			name:      "Focus color",
			colorName: theme.ColorNameFocus,
			variant:   theme.VariantLight,
			expected:  color.RGBA{R: 99, G: 102, B: 241, A: 100},
		},
		{
			name:      "Button color",
			colorName: theme.ColorNameButton,
			variant:   theme.VariantLight,
			expected:  color.RGBA{R: 99, G: 102, B: 241, A: 255},
		},
		{
			name:      "Disabled color",
			colorName: theme.ColorNameDisabled,
			variant:   theme.VariantLight,
			expected:  color.RGBA{R: 148, G: 163, B: 184, A: 255},
		},
		{
			name:      "Input background",
			colorName: theme.ColorNameInputBackground,
			variant:   theme.VariantLight,
			expected:  color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := customTheme.Color(tt.colorName, tt.variant)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCustomTheme_Font(t *testing.T) {
	customTheme := &CustomTheme{}

	// Test that font returns default theme font
	font := customTheme.Font(fyne.TextStyle{Bold: true})
	assert.NotNil(t, font)

	font = customTheme.Font(fyne.TextStyle{Italic: true})
	assert.NotNil(t, font)

	font = customTheme.Font(fyne.TextStyle{})
	assert.NotNil(t, font)
}

func TestCustomTheme_Icon(t *testing.T) {
	customTheme := &CustomTheme{}

	// Test that icon returns default theme icons
	icon := customTheme.Icon(theme.IconNameConfirm)
	assert.NotNil(t, icon)

	icon = customTheme.Icon(theme.IconNameDelete)
	assert.NotNil(t, icon)

	icon = customTheme.Icon(theme.IconNameSearch)
	assert.NotNil(t, icon)
}

func TestCustomTheme_Size(t *testing.T) {
	customTheme := &CustomTheme{}

	tests := []struct {
		name     string
		sizeName fyne.ThemeSizeName
		expected float32
	}{
		{
			name:     "Text size",
			sizeName: theme.SizeNameText,
			expected: 16,
		},
		{
			name:     "Heading text size",
			sizeName: theme.SizeNameHeadingText,
			expected: 24,
		},
		{
			name:     "Subheading text size",
			sizeName: theme.SizeNameSubHeadingText,
			expected: 20,
		},
		{
			name:     "Caption text size",
			sizeName: theme.SizeNameCaptionText,
			expected: 12,
		},
		{
			name:     "Padding size",
			sizeName: theme.SizeNamePadding,
			expected: 12,
		},
		{
			name:     "Inner padding size",
			sizeName: theme.SizeNameInnerPadding,
			expected: 8,
		},
		{
			name:     "Scrollbar size",
			sizeName: theme.SizeNameScrollBar,
			expected: 12,
		},
		{
			name:     "Scrollbar small size",
			sizeName: theme.SizeNameScrollBarSmall,
			expected: 8,
		},
		{
			name:     "Separator thickness",
			sizeName: theme.SizeNameSeparatorThickness,
			expected: 1,
		},
		{
			name:     "Input border size",
			sizeName: theme.SizeNameInputBorder,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := customTheme.Size(tt.sizeName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCustomTheme_Size_Default(t *testing.T) {
	customTheme := &CustomTheme{}

	// Test that unknown size names fall back to default theme
	result := customTheme.Size("unknown-size")
	assert.NotNil(t, result)
}
