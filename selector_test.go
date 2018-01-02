package warpcache

import "testing"

func TestSelector(t *testing.T) {
	tests := []struct {
		selector Selector
		Result   string
	}{
		{
			Selector{
				Classname: "money",
			}, "money{}",
		},
		{
			Selector{
				Classname: "money",
				Labels: map[string]string{
					"a": "b",
				},
			}, "money{a=b}",
		},
		{
			Selector{
				Classname: "money",
				Labels: map[string]string{
					"a": "b",
					"c": "d",
				},
			}, "money{a=b,c=d}",
		},
	}

	for _, test := range tests {
		result := test.selector.String()
		if result != test.Result {
			t.Errorf("Expected %s, got %s", test.Result, result)
		}
	}
}
