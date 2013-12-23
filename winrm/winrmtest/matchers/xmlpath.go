package matchers

import (
	"bytes"
	"fmt"
	"github.com/onsi/gomega"
	"io"
	"io/ioutil"
	"launchpad.net/xmlpath"
)

func MatchXmlPath(expected string) gomega.OmegaMatcher {
	return &xpathMatcher{
		text: expected,
		path: xmlpath.MustCompile(expected),
	}
}

type xpathMatcher struct {
	text string
	path *xmlpath.Path
}

func (matcher *xpathMatcher) Match(actual interface{}) (success bool, message string, err error) {
	reader, ok := actual.(io.Reader)
	if !ok {
		return false, "", fmt.Errorf("MatchXmlPath expects an io.Reader")
	}

	buffer, _ := ioutil.ReadAll(reader)
	if seeker, ok := reader.(io.Seeker); ok {
		seeker.Seek(0, 0)
	}

	node, err := xmlpath.Parse(bytes.NewReader(buffer))

	if err != nil {
		return false, "", err
	}

	xml := string(buffer)
	_, ok = matcher.path.String(node)

	if ok {
		return true, fmt.Sprintf("Expected\n\t%#v\nnot to match xml-path\n\t%#v", xml, matcher.text), nil
	} else {
		return false, fmt.Sprintf("Expected\n\t%#v\nto to match xml-path\n\t%#v", xml, matcher.text), nil
	}
}
