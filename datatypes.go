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
	SpacingUnknown SpacingType = iota
	SpacingNone
	SpacingBefore
	SpacingAfter
	SpacingBoth
)

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
