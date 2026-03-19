package raymond

import "testing"

func TestNarrowBug(t *testing.T) {
	tests := []struct {
		name     string
		template string
		ctx      map[string]interface{}
		want     string
	}{
		{
			// From the JSmon template: the pattern around futureEvent
			// "...{{/if}}\n{{#if futureEvent}}\n⚠️ Event happening{{/if}}\n{{{quickMoveEmoji}}}..."
			// Note: the content inside the block does NOT end with \n before {{/if}}
			// The {{/if}} is on the SAME line as the content
			name:     "standalone open, inline close - false",
			template: "Before\n{{#if show}}\nContent{{/if}}\nAfter",
			ctx:      map[string]interface{}{},
			want:     "Before\nAfter",
		},
		{
			name:     "standalone open, inline close - true",
			template: "Before\n{{#if show}}\nContent{{/if}}\nAfter",
			ctx:      map[string]interface{}{"show": true},
			want:     "Before\nContent\nAfter",
		},
		{
			// The exact pattern: {{/if}} on same line as content, then \n, then standalone block
			name:     "inline close with content, then newline, then false standalone",
			template: "{{#if a}}A content{{/if}}\n{{#if b}}\nB content{{/if}}\nAfter",
			ctx:      map[string]interface{}{"a": true, "b": false},
			want:     "A content\nAfter",
		},
		{
			// Even simpler - the futureEvent block pattern
			name:     "futureEvent pattern exact",
			template: "Weather line\n{{#if futureEvent}}\n⚠️ Event{{/if}}\nMoves line",
			ctx:      map[string]interface{}{"futureEvent": false},
			want:     "Weather line\nMoves line",
		},
		{
			name:     "futureEvent pattern exact - true",
			template: "Weather line\n{{#if futureEvent}}\n⚠️ Event{{/if}}\nMoves line",
			ctx:      map[string]interface{}{"futureEvent": true},
			want:     "Weather line\n⚠️ Event\nMoves line",
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
