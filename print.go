// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package main

import (
	"fmt"
)

// PrettyPrintText displays a pretty-printed version of the Text content tree
// of tags and license text content.
func PrettyPrintText(t []TextContent) {
	pprintHelper(t, 0)
}

func pprintHelper(t []TextContent, indent int) {
	indentBump := 2
	for _, elt := range t {
		for i := 0; i < indent; i++ {
			fmt.Printf(" ")
		}
		switch elt.Type {
		case TCCharData:
			fmt.Printf("- TCCharData: %s\n", elt.CharData)
		case TCParagraph:
			fmt.Printf("- TCParagraph:\n")
			pprintHelper(elt.ParaData, indent + indentBump)
		case TCBullet:
			fmt.Printf("- TCBullet: %s\n", elt.Bullet)
		case TCList:
			fmt.Printf("- TCList:\n")
			for _, lelt := range elt.ListItems {
				for i := 0; i < indent + indentBump; i++ {
					fmt.Printf(" ")
				}
				fmt.Printf("- TCItem:\n")
				pprintHelper(lelt, indent + (indentBump * 2))
			}
		case TCOptional:
			fmt.Printf("- TCOptional (spacing: %s):\n", getSpacingString(elt.OptionalSpacing))
			pprintHelper(elt.OptionalData, indent + indentBump)
		case TCAlt:
			fmt.Printf("- TCAlt (spacing: %s, name: %s, match: %s): %s\n", getSpacingString(elt.AltSpacing),
				elt.AltName, elt.AltMatch, elt.AltFixedText)
		case TCBr:
			fmt.Printf("- TCBr\n")
		case TCTitleText:
			fmt.Printf("- TCTitleText:\n")
			pprintHelper(elt.TitleData, indent + indentBump)
		case TCCopyrightText:
			fmt.Printf("- TCCopyrightText:\n")
			pprintHelper(elt.CopyrightData, indent + indentBump)
		case TCStandardLicenseHeader:
			fmt.Printf("- TCStandardLicenseHeader:\n")
			pprintHelper(elt.SLHData, indent + indentBump)
		case TCUnknown:
			fmt.Printf("- TCUnknown\n")
		default:
			fmt.Printf("- DEFAULT-UNKNOWN\n")
		}
	}
}
