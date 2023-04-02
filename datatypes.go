// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package main

type SPDXLicenseCollection struct {
	Licenses	[]License	`xml:"license"`
//	Exceptions	[]Exception	`xml:"exception"`
}

type License struct {
	IsOSIApproved		bool		`xml:"isOsiApproved,attr"`
	Name			string		`xml:"name,attr"`
	LicenseID		string		`xml:"licenseId,attr"`
	ListVersionAdded	string		`xml:"listVersionAdded,attr"`
//	ObsoltedBys -- need to parse sub-elements?
//	CrossRefs -- need to parse sub-elements?
	Notes			Text		`xml:"notes"`
	Text			Text		`xml:"text"`
//	StandardLicenseHeader -- outside Text?
}

// defining a separate Text type in order to define UnmarshalXML on it
type Text []TextContent

type TextContentType int

const (
	TCUnknown TextContentType = iota
	TCCharData
	TCParagraph
	TCBullet
	TCList
	TCOptional
	TCAlt
	TCBr
	TCTitleText
	TCCopyrightText
	TCStandardLicenseHeader
)

type SpacingType int

const (
	SpacingUnspecified SpacingType = iota
	SpacingNone
	SpacingBefore
	SpacingAfter
	SpacingBoth
	SpacingUnknown
)

// getSpacingType determines which type of spacing is chosen for this attribute.
func getSpacingType(s string) SpacingType {
	switch s {
	case "":
		return SpacingUnspecified
	case "none":
		return SpacingNone
	case "before":
		return SpacingBefore
	case "after":
		return SpacingAfter
	case "both":
		return SpacingBoth
	default:
		return SpacingUnknown
	}
}

// getSpacingString returns the standard string corresponding to the
// specified type of spacing.
func getSpacingString(st SpacingType) string {
	switch st {
	case SpacingUnspecified:
		return "<unspecified>"
	case SpacingNone:
		return "none"
	case SpacingBefore:
		return "before"
	case SpacingAfter:
		return "after"
	case SpacingBoth:
		return "both"
	case SpacingUnknown:
		return "UNKNOWN"
	default:
		return "DEFAULT-UNKNOWN"
	}
}

// TextContent is one element of text content, which may include various
// nested TextContent elements as well.
// It mostly corresponds to the formattedAltStandardLicenseHeaderTextType
// type in the license-list-XML ListedLicense.xsd schema.
type TextContent struct {
	// Type is the type of TextContent. It should be verified
	// before using any of the following attributes.
	Type		TextContentType

	// If Type == TCCharData:
	// CharData is plain character data which is not embedded
	// within any other elements within the surrounding element.
	CharData	string

	// If Type == TCParagraph:
	// ParaData is the slice of TextContent within this <p> element.
	ParaData	[]TextContent

	// If Type == TCBullet:
	// Bullet is the string within this <bullet> element.
	Bullet		string

	// If Type == TCList:
	// ListItems are the <item> elements within this <list>.
	ListItems	[]ListItem

	// If Type == TCOptional:
	// OptionalSpacing is the spacing attribute for this <optional> element.
	OptionalSpacing	SpacingType
	// OptionalData is the slice of TextContent within this <optional> element.
	OptionalData	[]TextContent

	// If Type == TCAlt:
	// AltFixedText is the character data content within the <alt> tags.
	AltFixedText	string
	// AltName is the name attribute for the <alt> element.
	AltName		string
	// AltMatch is the match attribute regular expression for the <alt> element.
	AltMatch	string
	// AltSpacing is the spacing attribute for this <alt> element.
	AltSpacing	SpacingType

	// If Type == TCBr:
	// no content needed for <br />, empty type

	// If Type == TCTitleText:
	// TitleData is the slice of TextContent within this <titleText> element.
	TitleData	[]TextContent

	// If Type == TCCopyrightText:
	// CopyrightData is the slice of TextContent within this <copyrightText> element.
	CopyrightData	[]TextContent

	// If Type == TCStandardLicenseHeader:
	// SLHData is the slice of TextContent within this <standardLicenseHeader> element.
	SLHData		[]TextContent

}

// ListItem is a slice of all TextContent elements contained within one
// <item> element in a <list>.
type ListItem []TextContent

// getTextContentType determines which type of TextContent we are entering
// (or exiting?) based on the element's name.
func getTextContentType(tag string) TextContentType {
	switch tag {
	case "p":
		return TCParagraph
	case "bullet":
		return TCBullet
	case "list":
		return TCList
	case "optional":
		return TCOptional
	case "alt":
		return TCAlt
	case "br":
		return TCBr
	case "titleText":
		return TCTitleText
	case "copyrightText":
		return TCCopyrightText
	case "standardLicenseHeader":
		return TCStandardLicenseHeader
	default:
		return TCUnknown
	}
}

// getTextContentString returns the standard element name string
// corresponding to the element's type.
func getTextContentString(t TextContentType) string {
	switch t {
	case TCCharData:
		return "CHARDATA"
	case TCParagraph:
		return "p"
	case TCBullet:
		return "bullet"
	case TCList:
		return "list"
	case TCOptional:
		return "optional"
	case TCAlt:
		return "alt"
	case TCBr:
		return "br"
	case TCTitleText:
		return "titleText"
	case TCCopyrightText:
		return "copyrightText"
	case TCStandardLicenseHeader:
		return "standardLicenseHeader"
	case TCUnknown:
		return "UNKNOWN"
	default:
		return "DEFAULT-UNKNOWN"
	}
}
