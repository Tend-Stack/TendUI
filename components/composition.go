package components

// Space selects a density-aware gap. Empty and unknown values use SpaceMedium.
type Space string

const (
	SpaceNone       Space = "none"
	SpaceExtraSmall Space = "xs"
	SpaceSmall      Space = "sm"
	SpaceMedium     Space = "md"
	SpaceLarge      Space = "lg"
	SpaceExtraLarge Space = "xl"
)

// LayoutAlign positions items across the flow direction without reordering them.
// Baseline applies to Cluster; Stack normalizes it to stretch.
type LayoutAlign string

const (
	LayoutAlignStart    LayoutAlign = "start"
	LayoutAlignCenter   LayoutAlign = "center"
	LayoutAlignEnd      LayoutAlign = "end"
	LayoutAlignStretch  LayoutAlign = "stretch"
	LayoutAlignBaseline LayoutAlign = "baseline"
)

// LayoutJustify distributes Cluster items along each wrapped line.
type LayoutJustify string

const (
	LayoutJustifyStart   LayoutJustify = "start"
	LayoutJustifyCenter  LayoutJustify = "center"
	LayoutJustifyEnd     LayoutJustify = "end"
	LayoutJustifyBetween LayoutJustify = "between"
)

type StackProps struct {
	ID  string
	Gap Space
	// Align defaults to stretch. Start and end follow the document writing direction.
	Align LayoutAlign
}

type ClusterProps struct {
	ID  string
	Gap Space
	// Align defaults to center; Justify defaults to start. Wrapping is always enabled.
	Align   LayoutAlign
	Justify LayoutJustify
}

func layoutSpace(value Space) string {
	return string(allowed(value, SpaceMedium, SpaceNone, SpaceExtraSmall, SpaceSmall, SpaceMedium, SpaceLarge, SpaceExtraLarge))
}

func stackAlign(value LayoutAlign) string {
	return string(allowed(value, LayoutAlignStretch, LayoutAlignStart, LayoutAlignCenter, LayoutAlignEnd, LayoutAlignStretch))
}

func clusterAlign(value LayoutAlign) string {
	return string(allowed(value, LayoutAlignCenter, LayoutAlignStart, LayoutAlignCenter, LayoutAlignEnd, LayoutAlignStretch, LayoutAlignBaseline))
}

func clusterJustify(value LayoutJustify) string {
	return string(allowed(value, LayoutJustifyStart, LayoutJustifyStart, LayoutJustifyCenter, LayoutJustifyEnd, LayoutJustifyBetween))
}
