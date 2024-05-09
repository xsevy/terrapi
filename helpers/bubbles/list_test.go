package bubbles

import "testing"

func TestListFocus(t *testing.T) {
	tcs := []struct {
		scenario      func(b *ListModel)
		name          string
		initialState  bool
		expectedState bool
	}{
		{
			name:          "initial state is false",
			scenario:      func(b *ListModel) {},
			initialState:  false,
			expectedState: false,
		},
		{
			name:          "initial state is true",
			scenario:      func(b *ListModel) {},
			initialState:  true,
			expectedState: true,
		},
		{
			name: "initial state is false -> blur",
			scenario: func(b *ListModel) {
				b.Blur()
			},
			initialState:  false,
			expectedState: false,
		},
		{
			name: "initial state is true -> blur",
			scenario: func(b *ListModel) {
				b.Blur()
			},
			initialState:  true,
			expectedState: false,
		},
		{
			name: "initial state is true -> focus",
			scenario: func(b *ListModel) {
				b.Focus()
			},
			initialState:  true,
			expectedState: true,
		},
		{
			name: "initial state is false -> focus",
			scenario: func(b *ListModel) {
				b.Focus()
			},
			initialState:  false,
			expectedState: true,
		},
		{
			name: "initial state is true -> focus -> blur",
			scenario: func(b *ListModel) {
				b.Focus()
				b.Blur()
			},
			initialState:  true,
			expectedState: false,
		},
		{
			name: "initial state is false -> blur -> focus",
			scenario: func(b *ListModel) {
				b.Blur()
				b.Focus()
			},
			initialState:  false,
			expectedState: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			l := NewListModel("text", []string{}, false, tc.initialState)

			tc.scenario(l)

			if l.focused != tc.expectedState {
				t.Errorf("%s expected %t, got %t", tc.name, tc.expectedState, l.focused)
			}
		})
	}
}

func TestListValue(t *testing.T) {
	i := []string{"a", "b", "c"}
	l := NewListModel("text", i, false, false)

	if l.Value() != i[0] {
		t.Errorf("expected %s got %s", i[0], l.Value())
	}

	maxNext := len(i)

	l.selected.Next(maxNext)
	if l.Value() != i[1] {
		t.Errorf("expected %s got %s", i[1], l.Value())
	}

	l.selected.Prev()
	if l.Value() != i[0] {
		t.Errorf("expected %s got %s", i[0], l.Value())
	}
}
