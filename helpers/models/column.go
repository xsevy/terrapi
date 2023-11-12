package models

type ColumnModel struct {
	focused bool
}

func (c ColumnModel) GetFocused() bool {
	return c.focused
}

func (c *ColumnModel) SetFocused(f bool) {
	c.focused = f
}
