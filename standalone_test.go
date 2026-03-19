package raymond

import "testing"

// TestStandaloneBlockWhitespace verifies that lines containing only block
// helpers (and whitespace) are removed from output, per the Mustache spec.
func TestStandaloneBlockWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		template string
		ctx      map[string]interface{}
		want     string
	}{
		{
			name:     "standalone if",
			template: "Line 1\n{{#if show}}\nContent\n{{/if}}\nLine 2",
			ctx:      map[string]interface{}{"show": true},
			want:     "Line 1\nContent\nLine 2",
		},
		{
			name:     "standalone if false",
			template: "Line 1\n{{#if show}}\nContent\n{{/if}}\nLine 2",
			ctx:      map[string]interface{}{"show": false},
			want:     "Line 1\nLine 2",
		},
		{
			name:     "standalone if/else true",
			template: "Line 1\n{{#if show}}\nYes\n{{else}}\nNo\n{{/if}}\nLine 2",
			ctx:      map[string]interface{}{"show": true},
			want:     "Line 1\nYes\nLine 2",
		},
		{
			name:     "standalone if/else false",
			template: "Line 1\n{{#if show}}\nYes\n{{else}}\nNo\n{{/if}}\nLine 2",
			ctx:      map[string]interface{}{"show": false},
			want:     "Line 1\nNo\nLine 2",
		},
		{
			name:     "standalone each",
			template: "Before\n{{#each items}}\n- {{this}}\n{{/each}}\nAfter",
			ctx:      map[string]interface{}{"items": []string{"a", "b"}},
			want:     "Before\n- a\n- b\nAfter",
		},
		{
			name:     "inline if should not strip",
			template: "Hello {{#if show}}World{{/if}} Bye",
			ctx:      map[string]interface{}{"show": true},
			want:     "Hello World Bye",
		},
		{
			name:     "nested standalone blocks",
			template: "Start\n{{#if outer}}\n{{#if inner}}\nDeep\n{{/if}}\n{{/if}}\nEnd",
			ctx:      map[string]interface{}{"outer": true, "inner": true},
			want:     "Start\nDeep\nEnd",
		},
		{
			name:     "standalone with indentation",
			template: "Line 1\n  {{#if show}}\n  Content\n  {{/if}}\nLine 2",
			ctx:      map[string]interface{}{"show": true},
			want:     "Line 1\n  Content\nLine 2",
		},
		{
			name:     "DTS-style consecutive ifs",
			template: "**{{name}}**\n{{#if iv}}IV: {{iv}}%\n{{/if}}\n{{#if cp}}CP: {{cp}}\n{{/if}}\n{{#if level}}Level: {{level}}\n{{/if}}\nExpires: {{time}}",
			ctx:      map[string]interface{}{"name": "Bulbasaur", "iv": "100", "cp": "500", "level": "15", "time": "12:00"},
			want:     "**Bulbasaur**\nIV: 100%\nCP: 500\nLevel: 15\nExpires: 12:00",
		},
		{
			name:     "DTS-style consecutive ifs partial",
			template: "**{{name}}**\n{{#if iv}}IV: {{iv}}%\n{{/if}}\n{{#if cp}}CP: {{cp}}\n{{/if}}\n{{#if level}}Level: {{level}}\n{{/if}}\nExpires: {{time}}",
			ctx:      map[string]interface{}{"name": "Bulbasaur", "iv": "100", "time": "12:00"},
			want:     "**Bulbasaur**\nIV: 100%\nExpires: 12:00",
		},
		{
			name:     "content on same line as closing tag",
			template: "Start\n{{#if show}}Content{{/if}}\nEnd",
			ctx:      map[string]interface{}{"show": true},
			want:     "Start\nContent\nEnd",
		},
		{
			name:     "content on same line as opening tag",
			template: "Start\n{{#if show}}Content\n{{/if}}\nEnd",
			ctx:      map[string]interface{}{"show": true},
			want:     "Start\nContent\nEnd",
		},
		{
			name:     "empty line between blocks",
			template: "Header\n\n{{#if a}}\nA: {{a}}\n{{/if}}\n\n{{#if b}}\nB: {{b}}\n{{/if}}\n\nFooter",
			ctx:      map[string]interface{}{"a": "yes", "b": "yes"},
			want:     "Header\n\nA: yes\n\nB: yes\n\nFooter",
		},
		{
			name:     "windows line endings",
			template: "Line 1\r\n{{#if show}}\r\nContent\r\n{{/if}}\r\nLine 2",
			ctx:      map[string]interface{}{"show": true},
			want:     "Line 1\r\nContent\r\nLine 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.template, tt.ctx)
			if err != nil {
				t.Fatalf("Render error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got:\n%q\nwant:\n%q", got, tt.want)
			}
		})
	}
}
