package main

import (
	mymath "math/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // chars of randomly generated comment
	length = 16 // length of randomly generated comment
)

type snippet struct {
	Body []byte `datastore:",noindex"`
}

// get sha256 with slight adjustments (what were they thinking with static salts it's so pointless)
func (s *snippet) ID() string {
	h := sha256.New()
	io.WriteString(h, "Go playground salt\n")
	h.Write(s.Body)
	sum := h.Sum(nil)
	b := make([]byte, base64.URLEncoding.EncodedLen(len(sum)))
	base64.URLEncoding.Encode(b, sum)

	//  If there is an underscore at the end of the substring, extend it until there is not.
	hashLen := 11
	for hashLen <= len(b) && b[hashLen-1] == '_' {
		hashLen++
	}
	return string(b)[:hashLen]
}

func getURL(content string) string {
	s := snippet {
		Body:	[]byte(content),
	}
	return s.ID()
	
}

func main() {
	mymath.Seed(time.Now().UnixNano()) // seed the math crypto
	targetURLPrefix := "abc" // this is the target URL prefix you want to find
	startTime := time.Now() // its fun to see how long it took to find cool url
	caseSensitiveSearch := true // true: you care about lower/uppercases, false: any mix is fine (faster)
	
	for {
		randomString := getRandomStringFast()
		// ensure automatically generate comment has newline appended cuz that gets added if someone runs the code on the website (if you share after running it should not change the url)
		content := readCodeFromFile() + randomString + "\n"

		url := getURL(content)
		var success bool
		if caseSensitiveSearch {
			success = startsWithSensitive(targetURLPrefix, url)
		} else {
			success = startsWith(targetURLPrefix, url)
		}
		if success {
			fmt.Printf("Random string: %v\n", randomString)
			fmt.Printf("URL: %v\n", getURL(content))
			endTime := time.Now()
			fmt.Printf("Execution time: %.2f seconds\n", endTime.Sub(startTime).Seconds())
			break
		}
	}
	
}

// uses math rand and is faster than the crypto one i tried
func getRandomStringFast() string {
    // Create a byte slice of the required length
    bytes := make([]byte, length)

    // Generate random bytes using math/rand
    for i := range bytes {
        bytes[i] = alphabet[mymath.Intn(len(alphabet))]
    }

    // Convert byte slice to string
    return string(bytes)
}

// case insensitive way to figure out if string starts with certain prefix (much faster than startsWithSensitive)
func startsWith(substring, str string) bool {
	// Convert both substring and str to lowercase
	substringLower := strings.ToLower(substring)
	strLower := strings.ToLower(str)
	// Use strings.HasPrefix with the lowercase versions
	return strings.HasPrefix(strLower, substringLower)
}

// case sensitive way to figure out if string starts with certain prefix
func startsWithSensitive(substring, str string) bool {
	return strings.HasPrefix(str, substring)
}

// read the code base construct from file (this code will generate a comment that will be the reason for the specified output url with prefix of choice)
func readCodeFromFile() string {
	filePath := "code.txt"
	content, err := ioutil.ReadFile(filePath)
    if err != nil {
        panic(err)
    }
	contentString := string(content)
	
	// line below is crucial on windows to replace the windows style newlines \r\n with the \n newlines the playground uses
	contentCleaned := strings.ReplaceAll(contentString, "\r\n", "\n")
	
	return string(contentCleaned)

}
