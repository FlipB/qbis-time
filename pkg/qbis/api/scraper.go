package api

import (
	"fmt"
	"io"
	"regexp"

	"golang.org/x/net/html"
)

//getCurrentUserFromCurrentLoginBlock parses javascript looking for currentLogin-object
//and returns the currentUser value (employee ID)
// TODO improve this to get more than just the employee ID
func getCurrentUserFromCurrentLoginBlock(script string) (string, error) {
	regex, err := regexp.Compile(".*var currentLogin = {((?:\n|.)*)};.*")
	if err != nil {
		return "", err
	}
	bs := regex.FindSubmatch([]byte(script))
	if len(bs) != 2 {
		return "", fmt.Errorf("unexpected number of matching groups in javascript block")
	}

	json := bs[1]

	// currentUser is the employee Id
	regex, _ = regexp.Compile(`.*currentUser:\s+'(\d+)'.*`)
	groups := regex.FindSubmatch(json)

	if len(groups) != 2 {
		return "", fmt.Errorf("unexpected number of matching groups in javascript block")
	}
	return string(groups[1]), nil

	/*
		// The following code was an attempted to convert the javascript to JSON to be unmarshalled - still want to get this working to easier get access to the other fields
		json = []byte("{" + string(json) + "}")

		regex, _ = regexp.Compile(`(,|{)([a-zA-Z0-9]+?):(.+?)\n`)
		json = regex.ReplaceAllFunc(json, func(match []byte) []byte {
			r, _ := regexp.Compile(`(,|{)([a-zA-Z0-9]+?):(.+?)\n`)
			groups := r.FindSubmatch(match)
			// TODO guard for index out of bound errors
			if len(groups) == 0 {
				return match
			}
			leadingSeparator := string(groups[1])
			key := string(groups[2])
			quotedValue := string(groups[3])

			quotedValue = strings.TrimSpace(quotedValue)
			Value := strings.Trim(quotedValue, `'`)

			return []byte(leadingSeparator + `"` + key + `": "` + Value + `"`)
		})
	*/

}

//getEmbeddedScriptsInHTML returns the content of all script tags in the reader
func getEmbeddedScriptsInHTML(r io.Reader) ([]string, error) {

	scripts := make([]string, 0)

	htmlTokenizer := html.NewTokenizer(r)
	getNext := false
lbreak:
	for {
		tt := htmlTokenizer.Next()
		switch tt {
		case html.ErrorToken:
			if htmlTokenizer.Err() == io.EOF {
				break lbreak
			}
			return nil, fmt.Errorf("error parsing html: %v", htmlTokenizer.Err())
		case html.TextToken:
			if getNext {
				scripts = append(scripts, string(htmlTokenizer.Text()))
				getNext = false
			}
		case html.StartTagToken, html.EndTagToken:
			getNext = false
			tn, _ := htmlTokenizer.TagName()
			if string(tn) == "script" {
				getNext = true
			}
		}
	}

	return scripts, nil
}
