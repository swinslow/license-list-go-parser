// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package main

import (
	"fmt"
	"encoding/xml"
	"io/ioutil"
	"os"
)

func main() {
	xf, err := os.Open("samples/MIT.xml")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer xf.Close()

	b, _ := ioutil.ReadAll(xf)

	var lics SPDXLicenseCollection

	xml.Unmarshal(b, &lics)

	fmt.Printf("%#v\n\n", lics)

	for _, l := range lics.Licenses {
		fmt.Printf("License: %s (%s)\n", l.Name, l.LicenseID)
		PrettyPrintText(l.Text)
	}
}
