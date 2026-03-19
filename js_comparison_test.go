package raymond

import (
	"fmt"
	"testing"
)

// TestJSComparison tests raymond against known handlebars.js v4.7.8 results.
// These are the EXACT outputs from Node.js handlebars v4.
func TestJSComparison(t *testing.T) {
	RemoveHelper("compare")
	RemoveHelper("eq")

	RegisterHelper("compare", func(lvalue interface{}, operator string, rvalue interface{}, options *Options) interface{} {
		if operator == "!=" {
			if Str(lvalue) != Str(rvalue) {
				return options.Fn()
			}
			return options.Inverse()
		}
		if operator == "<=" {
			return options.Inverse()
		}
		return options.Inverse()
	})
	defer RemoveHelper("compare")

	RegisterHelper("eq", func(a interface{}, b interface{}, options *Options) interface{} {
		if Str(a) == Str(b) {
			return options.Fn()
		}
		return options.Inverse()
	})
	defer RemoveHelper("eq")

	tests := []struct {
		name     string
		template string
		ctx      map[string]interface{}
		jsResult string // Exact result from handlebars.js v4.7.8
	}{
		{
			name:     "1. standalone open, inline close - FALSE",
			template: "Line1\n{{#if show}}\nContent{{/if}}\nLine2",
			ctx:      map[string]interface{}{"show": false},
			jsResult: "Line1\n\nLine2",
		},
		{
			name:     "2. standalone open, inline close - TRUE",
			template: "Line1\n{{#if show}}\nContent{{/if}}\nLine2",
			ctx:      map[string]interface{}{"show": true},
			jsResult: "Line1\nContent\nLine2",
		},
		{
			name:     "3. fully standalone - FALSE",
			template: "text\n{{#if show}}\ncontent\n{{/if}}\nmore",
			ctx:      map[string]interface{}{"show": false},
			jsResult: "text\nmore",
		},
		{
			name:     "4. fully standalone - TRUE",
			template: "text\n{{#if show}}\ncontent\n{{/if}}\nmore",
			ctx:      map[string]interface{}{"show": true},
			jsResult: "text\ncontent\nmore",
		},
		{
			name:     "5. two sequential standalone-open - both FALSE",
			template: "PrevLine\n{{#if weatherChange}}\nWeather info{{/if}}\n{{#if futureEvent}}\nEvent info{{/if}}\nMoveLine",
			ctx:      map[string]interface{}{"weatherChange": false, "futureEvent": false},
			jsResult: "PrevLine\n\n\nMoveLine",
		},
		{
			name:     "6. two sequential standalone-open - first TRUE, second FALSE",
			template: "PrevLine\n{{#if weatherChange}}\nWeather info{{/if}}\n{{#if futureEvent}}\nEvent info{{/if}}\nMoveLine",
			ctx:      map[string]interface{}{"weatherChange": true, "futureEvent": false},
			jsResult: "PrevLine\nWeather info\n\nMoveLine",
		},
		{
			name:     "7. mixed inline+standalone blocks - all FALSE",
			template: "CP 2500\n{{#compare addr '!=' 'Unknown'}}addr{{/compare}}{{#if size}}\nbig{{/if}}\n{{#if weatherChange}}\nweather{{/if}}\n{{#if futureEvent}}\nevent\n{{/if}}\nMoves",
			ctx:      map[string]interface{}{"addr": "Unknown", "size": 0, "weatherChange": false, "futureEvent": false},
			jsResult: "CP 2500\n\n\nMoves",
		},
		{
			name:     "8. mixed inline+standalone - weather TRUE",
			template: "CP 2500\n{{#compare addr '!=' 'Unknown'}}addr{{/compare}}{{#if size}}\nbig{{/if}}\n{{#if weatherChange}}\nweather{{/if}}\n{{#if futureEvent}}\nevent\n{{/if}}\nMoves",
			ctx:      map[string]interface{}{"addr": "Unknown", "size": 0, "weatherChange": true, "futureEvent": false},
			jsResult: "CP 2500\n\nweather\nMoves",
		},
		{
			name:     "9. prev block output + empty fully standalone",
			template: "{{#if prev}}output{{/if}}\n{{#if show}}\ncontent\n{{/if}}\nmore",
			ctx:      map[string]interface{}{"prev": true, "show": false},
			jsResult: "output\nmore",
		},
		{
			name:     "10. prev block empty + empty fully standalone",
			template: "{{#if prev}}output{{/if}}\n{{#if show}}\ncontent\n{{/if}}\nmore",
			ctx:      map[string]interface{}{"prev": false, "show": false},
			jsResult: "\nmore",
		},
		{
			name:     "11. three sequential fully standalone - all FALSE",
			template: "Start\n{{#if a}}\nAAA\n{{/if}}\n{{#if b}}\nBBB\n{{/if}}\n{{#if c}}\nCCC\n{{/if}}\nEnd",
			ctx:      map[string]interface{}{"a": false, "b": false, "c": false},
			jsResult: "Start\nEnd",
		},
		{
			name:     "12. three sequential fully standalone - middle TRUE",
			template: "Start\n{{#if a}}\nAAA\n{{/if}}\n{{#if b}}\nBBB\n{{/if}}\n{{#if c}}\nCCC\n{{/if}}\nEnd",
			ctx:      map[string]interface{}{"a": false, "b": true, "c": false},
			jsResult: "Start\nBBB\nEnd",
		},
		{
			name:     "13. eq standalone open, inline close - FALSE",
			template: "Before\n{{#eq pokemonId 570}}\nZorua text{{/eq}}\nAfter",
			ctx:      map[string]interface{}{"pokemonId": 25},
			jsResult: "Before\n\nAfter",
		},
		{
			name:     "14. compare standalone open, inline close - FALSE",
			template: "Before\n{{#compare addr '!=' 'Unknown'}}\naddr text{{/compare}}\nAfter",
			ctx:      map[string]interface{}{"addr": "Unknown"},
			jsResult: "Before\n\nAfter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tpl, err := Parse(tt.template)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			df := NewDataFrame()
			result, err := tpl.ExecWith(tt.ctx, df)
			if err != nil {
				t.Fatalf("Exec error: %v", err)
			}
			if result != tt.jsResult {
				t.Errorf("DIVERGES FROM JS\n  template: %q\n  raymond:  %q\n  js v4:    %q", tt.template, result, tt.jsResult)
				fmt.Printf("  raymond visual:\n---\n%s\n---\n", result)
				fmt.Printf("  js visual:\n---\n%s\n---\n", tt.jsResult)
			}
		})
	}
}
