package raymond

import "testing"

// TestStandaloneFalseBlockAfterInlineClose isolates the bug:
// when a false standalone block appears right after an inline {{/if}}\n,
// the \n before the standalone block should be consumed.
func TestStandaloneFalseBlockAfterInlineClose(t *testing.T) {
	tests := []struct {
		name     string
		template string
		ctx      map[string]interface{}
		want     string
	}{
		{
			// The core bug pattern:
			// "...content{{/if}}\n{{#if false}}\ncontent\n{{/if}}\nnext"
			// The {{#if false}} block is standalone. When false, the entire
			// block (including its surrounding whitespace lines) should vanish.
			// But the \n at end of "...content{{/if}}\n" is NOT part of the
			// standalone block's whitespace — it belongs to the previous content.
			// Node handlebars strips the \n before the standalone {{#if}} line
			// when the block is false, because the standalone stripping consumes
			// the trailing whitespace of the previous content node.
			name:     "false standalone block after inline close",
			template: "{{#if weather}}Weather{{/if}}\n{{#if futureEvent}}\nEvent\n{{/if}}\nMoves",
			ctx:      map[string]interface{}{"weather": true},
			want:     "Weather\nMoves",
		},
		{
			// Same but with the if wrapper
			name:     "inline if then false standalone if",
			template: "{{#if weather}}Weather{{/if}}\n{{#if event}}\nEvent\n{{/if}}\nMoves",
			ctx:      map[string]interface{}{"weather": true},
			want:     "Weather\nMoves",
		},
		{
			name:     "inline if then true standalone if",
			template: "{{#if weather}}Weather{{/if}}\n{{#if event}}\nEvent\n{{/if}}\nMoves",
			ctx:      map[string]interface{}{"weather": true, "event": true},
			want:     "Weather\nEvent\nMoves",
		},
		{
			// Simpler: standalone false block should not leave blank line
			name:     "content then false standalone block",
			template: "Line 1\n{{#if show}}\nContent\n{{/if}}\nLine 2",
			ctx:      map[string]interface{}{},
			want:     "Line 1\nLine 2",
		},
		{
			// Two sequential standalone false blocks
			name:     "two sequential false standalone blocks",
			template: "Start\n{{#if a}}\nA\n{{/if}}\n{{#if b}}\nB\n{{/if}}\nEnd",
			ctx:      map[string]interface{}{},
			want:     "Start\nEnd",
		},
		{
			// Mixed: first true, second false
			name:     "first true then false standalone",
			template: "Start\n{{#if a}}\nA\n{{/if}}\n{{#if b}}\nB\n{{/if}}\nEnd",
			ctx:      map[string]interface{}{"a": true},
			want:     "Start\nA\nEnd",
		},
		{
			// Mixed: first false, second true
			name:     "first false then true standalone",
			template: "Start\n{{#if a}}\nA\n{{/if}}\n{{#if b}}\nB\n{{/if}}\nEnd",
			ctx:      map[string]interface{}{"b": true},
			want:     "Start\nB\nEnd",
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
