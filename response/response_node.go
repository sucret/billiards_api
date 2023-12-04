package response

type NodeTree struct {
	NodeID       int32      `json:"node_id"`
	Title        string     `json:"title"`
	Name         string     `json:"name"`
	Path         string     `json:"path"`
	Icon         string     `json:"icon"`
	Type         int        `json:"type"`
	ParentNodeID int32      `json:"parent_node_id"`
	Component    string     `json:"component"`
	Redirect     string     `json:"redirect"`
	Children     []NodeTree `json:"children"`
}
