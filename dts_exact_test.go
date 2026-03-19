package raymond

import (
	"testing"
)

// TestExactUserDTSPattern tests the exact patterns from the user's
// PoracleJS DTS to verify raymond handles them identically.
func TestExactUserDTSPattern(t *testing.T) {
	// Register minimal helpers needed
	RegisterHelper("eq", func(a, b interface{}, options *Options) interface{} {
		if a == b {
			return options.Fn()
		}
		return options.Inverse()
	})

	tests := []struct {
		name     string
		template string
		ctx      map[string]interface{}
		want     string
	}{
		{
			// PoracleJS original: {{#if}} at start of line with content, closing
			// tag inline, followed by another {{#if}} inline
			name: "PoracleJS streetName+size pattern",
			template: "📊 CP 500 | Lvl 15\n" +
				"{{#if streetName}}🧭 Main St{{/if}}{{#if size}}\n" +
				"📐Big{{/if}}\n" +
				"Moves",
			ctx:  map[string]interface{}{"streetName": "Main St", "size": 5},
			want: "📊 CP 500 | Lvl 15\n🧭 Main St\n📐Big\nMoves",
		},
		{
			name: "PoracleJS streetName+size - no size",
			template: "📊 CP 500 | Lvl 15\n" +
				"{{#if streetName}}🧭 Main St{{/if}}{{#if size}}\n" +
				"📐Big{{/if}}\n" +
				"Moves",
			ctx:  map[string]interface{}{"streetName": "Main St"},
			want: "📊 CP 500 | Lvl 15\n🧭 Main St\nMoves",
		},
		{
			name: "PoracleJS streetName+size - no street",
			template: "📊 CP 500 | Lvl 15\n" +
				"{{#if streetName}}🧭 Main St{{/if}}{{#if size}}\n" +
				"📐Big{{/if}}\n" +
				"Moves",
			ctx:  map[string]interface{}{"size": 5},
			// The empty {{#if streetName}} block between the content and {{#if size}}
			// prevents {{#if size}} from being standalone, so the \n in its body is preserved.
			want: "📊 CP 500 | Lvl 15\n\n📐Big\nMoves",
		},
		{
			name: "PoracleJS streetName+size - neither",
			template: "📊 CP 500 | Lvl 15\n" +
				"{{#if streetName}}🧭 Main St{{/if}}{{#if size}}\n" +
				"📐Big{{/if}}\n" +
				"Moves",
			ctx:  map[string]interface{}{},
			// Both inline blocks produce empty strings; the \n before {{#if size}}
			// body survives because {{#if size}} is not standalone.
			want: "📊 CP 500 | Lvl 15\n\nMoves",
		},
		{
			// PoracleJS: eq block at end of line
			name: "PoracleJS eq block at end of line",
			template: "📊 CP 500 | Lvl 15{{#eq pokemonId 570}}\n" +
				"ℹ️ Change your buddy{{/eq}}\n" +
				"Moves",
			ctx:  map[string]interface{}{"pokemonId": 570},
			want: "📊 CP 500 | Lvl 15\nℹ️ Change your buddy\nMoves",
		},
		{
			name: "PoracleJS eq block at end of line - not matching",
			template: "📊 CP 500 | Lvl 15{{#eq pokemonId 570}}\n" +
				"ℹ️ Change your buddy{{/eq}}\n" +
				"Moves",
			ctx:  map[string]interface{}{"pokemonId": 25},
			want: "📊 CP 500 | Lvl 15\nMoves",
		},
		{
			// weatherChange pattern from user's PoracleJS DTS
			name: "PoracleJS weatherChange inline",
			template: "Moves line{{#if weatherChange}}\n" +
				"⚠️ Weather may change{{/if}}\n" +
				"Next line",
			ctx:  map[string]interface{}{"weatherChange": true},
			want: "Moves line\n⚠️ Weather may change\nNext line",
		},
		{
			name: "PoracleJS weatherChange false",
			template: "Moves line{{#if weatherChange}}\n" +
				"⚠️ Weather may change{{/if}}\n" +
				"Next line",
			ctx:  map[string]interface{}{},
			want: "Moves line\nNext line",
		},
		{
			// futureEvent pattern from user's PoracleJS DTS (standalone on own line)
			name: "PoracleJS futureEvent standalone",
			template: "Line before\n" +
				"{{#if futureEvent}}\n" +
				"⚠️ Event happening\n" +
				"{{/if}}\n" +
				"Line after",
			ctx:  map[string]interface{}{"futureEvent": true},
			want: "Line before\n⚠️ Event happening\nLine after",
		},
		{
			name: "PoracleJS futureEvent false standalone",
			template: "Line before\n" +
				"{{#if futureEvent}}\n" +
				"⚠️ Event happening\n" +
				"{{/if}}\n" +
				"Line after",
			ctx:  map[string]interface{}{},
			want: "Line before\nLine after",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.template, tt.ctx)
			if err != nil {
				t.Fatalf("Render error: %v", err)
			}
			if got != tt.want {
				t.Errorf("\ngot:  %q\nwant: %q", got, tt.want)
			}
		})
	}
}
