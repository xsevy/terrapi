package bubbles

import "testing"

func TestButtonFocus(t *testing.T) {
	tcs := []struct {
		scenario      func(b *ButtonModel)
		name          string
		initialState  bool
		expectedState bool
	}{
		{
			name:          "initial state is false",
			scenario:      func(b *ButtonModel) {},
			initialState:  false,
			expectedState: false,
		},
		{
			name:          "initial state is true",
			scenario:      func(b *ButtonModel) {},
			initialState:  true,
			expectedState: true,
		},
		{
			name: "initial state is false -> blur",
			scenario: func(b *ButtonModel) {
				b.Blur()
			},
			initialState:  false,
			expectedState: false,
		},
		{
			name: "initial state is true -> blur",
			scenario: func(b *ButtonModel) {
				b.Blur()
			},
			initialState:  true,
			expectedState: false,
		},
		{
			name: "initial state is true -> focus",
			scenario: func(b *ButtonModel) {
				b.Focus()
			},
			initialState:  true,
			expectedState: true,
		},
		{
			name: "initial state is false -> focus",
			scenario: func(b *ButtonModel) {
				b.Focus()
			},
			initialState:  false,
			expectedState: true,
		},
		{
			name: "initial state is true -> focus -> blur",
			scenario: func(b *ButtonModel) {
				b.Focus()
				b.Blur()
			},
			initialState:  true,
			expectedState: false,
		},
		{
			name: "initial state is false -> blur -> focus",
			scenario: func(b *ButtonModel) {
				b.Blur()
				b.Focus()
			},
			initialState:  false,
			expectedState: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			btn := NewButtonModel("text", tc.initialState)

			tc.scenario(btn)

			if btn.focused != tc.expectedState {
				t.Errorf("%s expected %t, got %t", tc.name, tc.expectedState, btn.focused)
			}
		})
	}
}

func TestButtonValue(t *testing.T) {
	btn := NewButtonModel("text", false)

	if btn.Value() != "" {
		t.Error("Button Value should return empty string")
	}
}
