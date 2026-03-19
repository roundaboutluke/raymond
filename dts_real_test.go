package raymond

import (
	"fmt"
	"testing"
)

func init() {
	// Register helpers needed for the DTS patterns
	RegisterHelper("compare", func(a interface{}, op string, b interface{}, options *Options) interface{} {
		sa, sb := fmt.Sprint(a), fmt.Sprint(b)
		switch op {
		case "!=":
			if sa != sb {
				return options.Fn()
			}
		case "<=":
			return options.Fn()
		case ">":
			return options.Fn()
		}
		return ""
	})
	RegisterHelper("or", func(a, b interface{}, options *Options) interface{} {
		if a != nil && a != false && a != "" && a != 0 {
			return options.Fn()
		}
		if b != nil && b != false && b != "" && b != 0 {
			return options.Fn()
		}
		return ""
	})
	RegisterHelper("lte", func(a, b interface{}) bool { return true })
	RegisterHelper("gte", func(a, b interface{}) bool { return true })
	RegisterHelper("capitalize", func(s string) string { return s })
	RegisterHelper("lowercase", func(s string) string { return s })
	RegisterHelper("map", func(a, b string) string { return "" })
	RegisterHelper("addCommas", func(v interface{}) string { return fmt.Sprint(v) })
	RegisterHelper("round", func(v interface{}) string { return fmt.Sprint(v) })
	RegisterHelper("getEmoji", func(s string) string { return "[" + s + "]" })
	RegisterHelper("getPowerUpCost", func(a, b interface{}, options *Options) interface{} {
		return options.FnWith(map[string]interface{}{"stardust": 5000, "candy": 10, "xlCandy": 2})
	})
	RegisterHelper("pvpSlug", func(s string) string { return s })
	RegisterHelper("isnt", func(a, b interface{}, options *Options) interface{} {
		if fmt.Sprint(a) != fmt.Sprint(b) {
			return options.Fn()
		}
		return ""
	})
}

// TestJSMonRenderedOutput renders the original PoracleJS DTS template through
// raymond and shows what raymond produces vs what Node handlebars would produce.
func TestJSMonRenderedOutput(t *testing.T) {
	// The original JSmon.txt template (simplified - focus on the problem areas)
	jsTemplate := "⏰ {{time}} ({{tthm}}m {{tths}}s){{#if confirmedTime}}✅{{else}}❌{{/if}}\n" +
		"📊 CP {{cp}} | Lvl {{level}}\n" +
		"{{#if streetName}}🧭 {{{addr}}}{{/if}}{{#if size}}{{#or (lte size 1) (gte size 5)}}\n" +
		"📐{{sizeName}}{{/or}}{{/if}}{{#eq pokemonId 570}}\n" +
		"\n" +
		"ℹ️ Zorua info{{/eq}}{{#if weatherChange}}\n" +
		"⚠️ Weather may change{{/if}}\n" +
		"{{#if futureEvent}}\n" +
		"⚠️ Event happening{{/if}}\n" +
		"{{{quickMoveEmoji}}} {{quickMoveName}} {{{chargeMoveEmoji}}} {{chargeMoveName}}\n" +
		"Disappears  <t:{{disappear_time}}:R>"

	ctx := map[string]interface{}{
		"time":          "12:00",
		"tthm":          5,
		"tths":          30,
		"confirmedTime": true,
		"cp":            500,
		"level":         15,
		"streetName":    "Main St",
		"addr":          "123 Main St",
		"size":          5,
		"sizeName":      "Big",
		"pokemonId":     25,  // not Zorua
		"weatherChange": true,
		"weatherCurrentEmoji": "☀️",
		"weatherNextEmoji":    "🌧️",
		"weatherNextName":     "rain",
		"futureEvent":         false,
		"futureEventName":     "",
		"quickMoveEmoji":      "⚡",
		"quickMoveName":       "Thunder Shock",
		"chargeMoveEmoji":     "💥",
		"chargeMoveName":      "Thunder",
		"disappear_time":      1234567890,
	}

	// What Node handlebars v4 would produce:
	// The key is that {{/if}}{{#if size}} are NOT standalone (they have content
	// around them on the line), so Node does NOT strip the line.
	// And {{/if}}{{#if weatherChange}} are also NOT standalone.
	nodeExpected := "⏰ 12:00 (5m 30s)✅\n" +
		"📊 CP 500 | Lvl 15\n" +
		"🧭 123 Main St\n" +
		"📐Big\n" +
		"⚠️ Weather may change\n" +
		"⚡ Thunder Shock 💥 Thunder\n" +
		"Disappears  <t:1234567890:R>"

	got, err := Render(jsTemplate, ctx)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	if got != nodeExpected {
		t.Errorf("Raymond produced different output than Node handlebars would.\n\nraymond got:\n%s\n\nNode expected:\n%s\n\nraymond quoted: %q\nNode quoted:    %q", got, nodeExpected, got, nodeExpected)
	}
}

// TestInlineCloseOpenNotStandalone verifies the core pattern:
// content{{/if}}{{#if other}}\n  should NOT be treated as standalone
func TestInlineCloseOpenNotStandalone(t *testing.T) {
	tests := []struct {
		name     string
		template string
		ctx      map[string]interface{}
		want     string
	}{
		{
			name:     "close+open inline, both true",
			template: "Before\n{{#if a}}A{{/if}}{{#if b}}\nB{{/if}}\nAfter",
			ctx:      map[string]interface{}{"a": true, "b": true},
			want:     "Before\nA\nB\nAfter",
		},
		{
			name:     "close+open inline, first true only",
			template: "Before\n{{#if a}}A{{/if}}{{#if b}}\nB{{/if}}\nAfter",
			ctx:      map[string]interface{}{"a": true, "b": false},
			want:     "Before\nA\nAfter",
		},
		{
			name:     "close+open inline, second true only",
			template: "Before\n{{#if a}}A{{/if}}{{#if b}}\nB{{/if}}\nAfter",
			ctx:      map[string]interface{}{"a": false, "b": true},
			want:     "Before\n\nB\nAfter",
		},
		{
			name:     "close+open inline, neither true",
			template: "Before\n{{#if a}}A{{/if}}{{#if b}}\nB{{/if}}\nAfter",
			ctx:      map[string]interface{}{"a": false, "b": false},
			want:     "Before\n\nAfter",
		},
		{
			// This is the critical one - content before closing tag means
			// it's NOT standalone, so the \n after {{/if}} should be preserved
			name:     "content before close, then open - newline preserved",
			template: "Line 1\n{{#if a}}content{{/if}}{{#if b}}\nB line{{/if}}\nLine 3",
			ctx:      map[string]interface{}{"a": false, "b": false},
			want:     "Line 1\n\nLine 3",
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
