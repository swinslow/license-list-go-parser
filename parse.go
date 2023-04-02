// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package main

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type parser struct {
	// tc is the Text ([]TextContent) slice we are building
	tc []TextContent

	// stack will be our stack of nested TextContent elements
	// so that we can track which one we're currently inside
	stack []TextContent
}

func (t *Text) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return parseText(t, d)
}

// addElement adds the current TextContent element either to the current
// top of stack element (in the appropriate place), or else to the main
// Text slice, as applicable.
func (p *parser) addElement(elt TextContent) error {
	if len(p.stack) > 0 {
		// there's at least one element on the stack,
		// so we'll fold this one into the appropriate holder
		// FIXME handle <list></list> elements appropriately
		topEltPtr := &(p.stack)[len(p.stack)-1]
		ttype := (*topEltPtr).Type
		fmt.Printf("adding %s to topElt %s\n", getTextContentString(elt.Type),
			getTextContentString(ttype))
		switch ttype {
		case TCParagraph:
			(*topEltPtr).ParaData = append((*topEltPtr).ParaData, elt)
			return nil
		case TCList:
			// FIXME figure out how to handle
			fmt.Printf("===> NOT HANDLED YET: TCList\n")
			return nil
		case TCOptional:
			(*topEltPtr).OptionalData = append((*topEltPtr).OptionalData, elt)
			return nil
		case TCAlt:
			// FIXME confirm -- should only work if this is CharData?
			if elt.Type == TCCharData {
				(*topEltPtr).AltFixedText = elt.CharData
				return nil
			} else {
				return fmt.Errorf("got non-CharData within <alt> tags\n")
			}
		case TCTitleText:
			(*topEltPtr).TitleData = append((*topEltPtr).TitleData, elt)
			return nil
		case TCCopyrightText:
			(*topEltPtr).CopyrightData = append((*topEltPtr).CopyrightData, elt)
			return nil
		case TCStandardLicenseHeader:
			(*topEltPtr).SLHData = append((*topEltPtr).SLHData, elt)
			return nil
		default:
			return fmt.Errorf("unexpected element type on stack: %s\n", getTextContentString(ttype))
		}
	} else {
		// add to main tc Text slice
		fmt.Printf("adding %s to main slice\n", getTextContentString(elt.Type))
		p.tc = append(p.tc, elt)
		return nil
	}
}

func addAttribs(elt *TextContent, se xml.StartElement) error {
	// assumes that the Type has already been filled in by parseText
	switch elt.Type {
	case TCBullet:
		for _, attr := range se.Attr {
			if attr.Name.Local == "spacing" {
				elt.Bullet = attr.Value
			} else {
				return fmt.Errorf("unknown attribute in <bullet> tag: %s", attr.Name.Local)
			}
		}
		return nil
	case TCOptional:
		for _, attr := range se.Attr {
			if attr.Name.Local == "spacing" {
				elt.OptionalSpacing = getSpacingType(attr.Value)
			} else {
				return fmt.Errorf("unknown attribute in <optional> tag: %s", attr.Name.Local)
			}
		}
		return nil
	case TCAlt:
		for _, attr := range se.Attr {
			if attr.Name.Local == "name" {
				elt.AltName = attr.Value
			} else if attr.Name.Local == "match" {
				elt.AltMatch = attr.Value
			} else {
				return fmt.Errorf("unknown attribute in <optional> tag: %s", attr.Name.Local)
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown attributes found in <%s> tag", getTextContentString(elt.Type))
	}
}

func parseText(t *Text, d *xml.Decoder) error {
	p := parser {tc: []TextContent{}, stack: []TextContent{}}

	// walk through each token, and based on what it is,
	// decide what we should do
	for {
		tok, err := d.Token()
		if err != nil {
			if tok == nil {
				// reached the end
				break
			}
			return err
		}

		// do something different depending on the token type
		switch tok.(type) {
			case xml.StartElement:
				se := tok.(xml.StartElement)
				fmt.Printf("StartElement: %s [%+v]\n", se.Name.Local, se.Attr)

				// create new element and add to stack
				newT := TextContent{Type:getTextContentType(se.Name.Local)}
				// handle attribs
				if len(se.Attr) > 0 {
					err := addAttribs(&newT, se)
					if err != nil {
						return err
					}
				}
				p.stack = append(p.stack, newT)

			case xml.EndElement:
				se := tok.(xml.EndElement)
				fmt.Printf("EndElement: %s\n", se.Name.Local)

				// check if we reached the end with the </text> element
				if se.Name.Local == "text" {
					// check that the stack is cleaned up?
					if len(p.stack) != 0 {
						elt := p.stack[len(p.stack)-1]
						return fmt.Errorf("reached </text> closing tag before end of prior tag <%s>", getTextContentString(elt.Type))
					}
					break
				}

				// check stack size before popping
				if len(p.stack) == 0 {
					return fmt.Errorf("reached </%s> closing tag without any active opening tag", se.Name.Local)
				}

				// check that the end elt matches the expected top of stack
				elt := p.stack[len(p.stack)-1]
				if getTextContentString(elt.Type) != se.Name.Local {
					return fmt.Errorf("reached </%s> closing tag but current active opening tag was <%s>", se.Name.Local, getTextContentString(elt.Type))
				}

				// handle popping and inserting into tc or parent, as appropriate
				p.stack = p.stack[:len(p.stack)-1]
				err := p.addElement(elt)
				if err != nil {
					return err
				}

			case xml.CharData:
				cd := tok.(xml.CharData)
				fmt.Printf("CharData: %s\n", string(cd))

				// remove external whitespace
				trimcd := strings.TrimSpace(string(cd))
				if trimcd != "" {
					// remove internal excess whitespace
					trimcd = strings.Join(strings.Fields(trimcd), " ")

					charElt := TextContent{
						Type: TCCharData,
						CharData: trimcd,
					}
					err := p.addElement(charElt)
					if err != nil {
						return err
					}
				}

			case xml.Comment:
				co := tok.(xml.Comment)
				fmt.Printf("Comment: %s\n", string(co))
			case xml.ProcInst:
				pi := tok.(xml.ProcInst)
				fmt.Printf("ProcInst: Target %s, Inst %s\n", pi.Target, string(pi.Inst))
				fmt.Printf("ERROR, processing instructions not supported")
				return fmt.Errorf("Processing instruction not supported: Target %s, Inst %s", pi.Target, string(pi.Inst))
			case xml.Directive:
				di := tok.(xml.Directive)
				fmt.Printf("Directive: %s\n", string(di))
			default:
				fmt.Printf("in UnmarshalXML, tok = %+v\n", tok)
		}
	}

	*t = p.tc
	return nil
}
