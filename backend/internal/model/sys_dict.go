package model

// SysDictOption represents a dict option for selection.
type SysDictOption struct {
	Label     string `json:"label"`
	Value     string `json:"value"`
	Color     string `json:"color,omitempty"`
	Icon      string `json:"icon,omitempty"`
	IsDefault bool   `json:"isDefault"`
}

// SysDictLabelOut is the output for fetching a label by value.
type SysDictLabelOut struct {
	Label string `json:"label"`
}
