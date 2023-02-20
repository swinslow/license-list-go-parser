// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type SPDXLicenseCollection struct {
	Licenses	[]License	`xml:"license"`
//	Exceptions	[]Exception	`xml:"exception"`
}

type License struct {
	IsOSIApproved		bool	`xml:"isOsiApproved,attr"`
	Name			string	`xml:"name,attr"`
	LicenseID		string	`xml:"licenseId,attr"`
	ListVersionAdded	string	`xml:"listVersionAdded,attr"`
	Text			Text	`xml:"text"`
}

/*
type Text struct {
//	TitleText	[]string	`xml:"titleText"`

}
*/
type Text string

func main() {
	xf, err := os.Open("samples/basic.xml")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer xf.Close()

	b, _ := ioutil.ReadAll(xf)

	var lics SPDXLicenseCollection

	xml.Unmarshal(b, &lics)

	fmt.Printf("%#v", lics)
}
