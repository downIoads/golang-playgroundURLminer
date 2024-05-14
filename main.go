package main

import (
	mymath "math/rand"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	// This salt is not meant to be kept secret (it’s checked in after all). It’s
	// a tiny bit of paranoia to avoid whatever problems a collision may cause.
	salt = "Go playground salt\n"
	maxSnippetSize = 64 * 1024
	
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // chars of randomly generated comment
	length = 32 // length of randomly generated comment
)

type snippet struct {
	Body []byte `datastore:",noindex"` // golang.org/issues/23253
}

func (s *snippet) ID() string {
	h := sha256.New()
	io.WriteString(h, salt)
	h.Write(s.Body)
	sum := h.Sum(nil)
	b := make([]byte, base64.URLEncoding.EncodedLen(len(sum)))
	base64.URLEncoding.Encode(b, sum)
	// Web sites don’t always linkify a trailing underscore, making it seem like
	// the link is broken. If there is an underscore at the end of the substring,
	// extend it until there is not.
	hashLen := 11
	for hashLen <= len(b) && b[hashLen-1] == '_' {
		hashLen++
	}
	return string(b)[:hashLen]
}

func GetURL(content string) string {
	s := snippet {
		Body:	[]byte(content),
	}
	return s.ID()
	
}

func main() {
	mymath.Seed(time.Now().UnixNano()) // seed the math crypto
	targetURLPrefix := "godev" // 5 chars takes around 1 minute on my pc (case insensitive search)
	startTime := time.Now() // its fun to see how long it took to find cool url
	caseSensitiveSearch := true
	
	for {
		//randomString := getRandomString()
		randomString := getRandomStringFast()

		// ctrl+A ctrl+C from playground and ensure that here you use multiline string
		content := `// You can edit this code!
// Click here and start typing.
package main

import "fmt"

func main() {
	fmt.Println("godev")
}

// ` + randomString + "\n"
		
		// ensure automatically generate comment has newline appended cuz that gets added if someone runs the code on the website (if you share after running it should not change the url)
		url := GetURL(content)
		var success bool
		if caseSensitiveSearch {
			success = startsWithSensitive(targetURLPrefix, url)
		} else {
			success = startsWith(targetURLPrefix, url)
		}
		if success {
			fmt.Printf("Random string: %v\n", randomString)
			fmt.Printf("URL: %v\n", GetURL(content))
			endTime := time.Now()
			fmt.Printf("Execution time: %.2f seconds\n", endTime.Sub(startTime).Seconds())
			break
		}
	}
	
}

// this uses crypto rand and is too slow, no point in this
func getRandomString() string {
    // Create a byte slice of the required length
    bytes := make([]byte, length)

    // Generate random bytes using crypto/rand
    _, err := rand.Read(bytes)
    if err != nil {
        panic(err)
    }

    // Map random bytes to the alphabet
    for i, b := range bytes {
        bytes[i] = alphabet[b%byte(len(alphabet))]
    }

    // Convert byte slice to string
    return string(bytes)
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

// made it case insensitive so that its easier to get hits
func startsWith(substring, str string) bool {
	// Convert both substring and str to lowercase
	substringLower := strings.ToLower(substring)
	strLower := strings.ToLower(str)
	// Use strings.HasPrefix with the lowercase versions
	return strings.HasPrefix(strLower, substringLower)
}

func startsWithSensitive(substring, str string) bool {
	return strings.HasPrefix(str, substring)
}

// cool urls (url: randomString):
	// godev: 	dcldvjfxjcxybjrujmxu
	