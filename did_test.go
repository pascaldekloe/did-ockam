package did

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestIsURL(t *testing.T) {
	t.Run("returns false if no Path or Fragment", func(t *testing.T) {
		d := &DID{Method: "example", ID: "123"}
		assert(t, false, d.IsURL())
	})

	t.Run("returns true if Path", func(t *testing.T) {
		d := &DID{Method: "example", ID: "123", Path: "a/b"}
		assert(t, true, d.IsURL())
	})

	t.Run("returns true if PathSegements", func(t *testing.T) {
		d := &DID{Method: "example", ID: "123", PathSegments: []string{"a", "b"}}
		assert(t, true, d.IsURL())
	})

	t.Run("returns true if Query", func(t *testing.T) {
		d := &DID{Method: "example", ID: "123", Query: "abc"}
		assert(t, true, d.IsURL())
	})

	t.Run("returns true if Fragment", func(t *testing.T) {
		d := &DID{Method: "example", ID: "123", Fragment: "00000"}
		assert(t, true, d.IsURL())
	})

	t.Run("returns true if Path and Fragment", func(t *testing.T) {
		d := &DID{Method: "example", ID: "123", Path: "a/b", Fragment: "00000"}
		assert(t, true, d.IsURL())
	})
}

func TestString(t *testing.T) {
	t.Run("assembles a DID", func(t *testing.T) {
		d := &DID{Method: "example", ID: "123"}
		assert(t, "did:example:123", d.String())
	})

	t.Run("assembles a DID from IDStrings", func(t *testing.T) {
		d := &DID{Method: "example", IDStrings: []string{"123", "456"}}
		assert(t, "did:example:123%3A456", d.String())
	})

	t.Run("returns empty string if no method", func(t *testing.T) {
		d := &DID{ID: "123"}
		assert(t, "", d.String())
	})

	t.Run("returns empty string in no ID or IDStrings", func(t *testing.T) {
		d := &DID{Method: "example"}
		assert(t, "", d.String())
	})

	t.Run("includes Path", func(t *testing.T) {
		d := &DID{Method: "example", ID: "123", Path: "a/b"}
		assert(t, "did:example:123/a/b", d.String())
	})

	t.Run("includes Path assembled from PathSegements", func(t *testing.T) {
		d := &DID{Method: "example", ID: "123", PathSegments: []string{"a", "b"}}
		assert(t, "did:example:123/a/b", d.String())
	})

	t.Run("includes Query after IDString", func(t *testing.T) {
		d := &DID{Method: "example", ID: "123", Query: "abc"}
		assert(t, "did:example:123?abc", d.String())
	})

	t.Run("includes Query after Path", func(t *testing.T) {
		d := &DID{Method: "example", ID: "123", Path: "x/y", Query: "abc"}
		assert(t, "did:example:123/x/y?abc", d.String())
	})

	t.Run("includes Query after before Fragment", func(t *testing.T) {
		d := &DID{Method: "example", ID: "123", Fragment: "zyx", Query: "abc"}
		assert(t, "did:example:123?abc#zyx", d.String())
	})

	t.Run("includes Query", func(t *testing.T) {
		d := &DID{Method: "example", ID: "123", Query: "abc"}
		assert(t, "did:example:123?abc", d.String())
	})

	t.Run("includes Fragment", func(t *testing.T) {
		d := &DID{Method: "example", ID: "123", Fragment: "00000"}
		assert(t, "did:example:123#00000", d.String())
	})

	t.Run("includes Fragment after Param", func(t *testing.T) {
		d := &DID{Method: "example", ID: "123", Fragment: "00000"}
		assert(t, "did:example:123#00000", d.String())
	})
}

func TestParse(t *testing.T) {

	t.Run("returns error if input is empty", func(t *testing.T) {
		_, err := Parse("")
		assert(t, false, err == nil)
	})

	t.Run("returns error if input length is less than length 7", func(t *testing.T) {
		_, err := Parse("did:")
		assert(t, false, err == nil)

		_, err = Parse("did:a")
		assert(t, false, err == nil)

		_, err = Parse("did:a:")
		assert(t, false, err == nil)
	})

	t.Run("returns error if input does not have a second : to mark end of method", func(t *testing.T) {
		_, err := Parse("did:aaaaaaaaaaa")
		assert(t, false, err == nil)
	})

	t.Run("returns error if method is empty", func(t *testing.T) {
		_, err := Parse("did::aaaaaaaaaaa")
		assert(t, false, err == nil)
	})

	t.Run("returns error if input does not begin with did: scheme", func(t *testing.T) {
		_, err := Parse("a:12345")
		assert(t, false, err == nil)
	})

	t.Run("returned value is nil if input does not begin with did: scheme", func(t *testing.T) {
		d, _ := Parse("a:12345")
		assert(t, true, d == nil)
	})

	t.Run("succeeds if it has did prefix and length is greater than 7", func(t *testing.T) {
		d, err := Parse("did:a:1")
		assert(t, nil, err)
		assert(t, true, d != nil)
	})

	t.Run("succeeds to extract method", func(t *testing.T) {
		d, err := Parse("did:a:1")
		assert(t, nil, err)
		assert(t, "a", d.Method)

		d, err = Parse("did:abcdef:11111")
		assert(t, nil, err)
		assert(t, "abcdef", d.Method)
	})

	t.Run("returns error if method has any other char than 0-9 or a-z", func(t *testing.T) {
		_, err := Parse("did:aA:1")
		assert(t, false, err == nil)

		_, err = Parse("did:aa-aa:1")
		assert(t, false, err == nil)
	})

	t.Run("succeeds to extract id", func(t *testing.T) {
		d, err := Parse("did:a:1")
		assert(t, nil, err)
		assert(t, "1", d.ID)
	})

	t.Run("succeeds to extract id parts", func(t *testing.T) {
		d, err := Parse("did:a:123:456")
		assert(t, nil, err)

		parts := d.IDStrings
		assert(t, "123", parts[0])
		assert(t, "456", parts[1])
	})

	t.Run("returns error if ID has an invalid char", func(t *testing.T) {
		_, err := Parse("did:a:1&&111")
		assert(t, false, err == nil)
	})

	t.Run("succeeds to extract path", func(t *testing.T) {
		d, err := Parse("did:a:123:456/someService")
		assert(t, nil, err)
		assert(t, "someService", d.Path)
	})

	t.Run("succeeds to extract path segements", func(t *testing.T) {
		d, err := Parse("did:a:123:456/a/b")
		assert(t, nil, err)

		segments := d.PathSegments
		assert(t, "a", segments[0])
		assert(t, "b", segments[1])
	})

	t.Run("succeeds with percent encoded chars in path", func(t *testing.T) {
		d, err := Parse("did:a:123:456/a/%20a")
		assert(t, nil, err)
		assert(t, "a/%20a", d.Path)
	})

	t.Run("returns error if % in path is not followed by 2 hex chars", func(t *testing.T) {
		dids := []string{
			"did:a:123:456/%",
			"did:a:123:456/%a",
			"did:a:123:456/%!*",
			"did:a:123:456/%A!",
			"did:xyz:pqr#%A!",
			"did:a:123:456/%A%",
		}
		for _, did := range dids {
			_, err := Parse(did)
			assert(t, false, err == nil, "Input: %s", did)
		}
	})

	t.Run("does not fail if second path segment is empty", func(t *testing.T) {
		_, err := Parse("did:a:123:456/abc//pqr")
		assert(t, nil, err)
	})

	t.Run("returns error  if path has invalid char", func(t *testing.T) {
		_, err := Parse("did:a:123:456/ssss^sss")
		assert(t, false, err == nil)
	})

	t.Run("does not fail if path has atleast one segment and a trailing slash", func(t *testing.T) {
		_, err := Parse("did:a:123:456/a/b/")
		assert(t, nil, err)
	})

	t.Run("succeeds to extract query after idstring", func(t *testing.T) {
		d, err := Parse("did:a:123?abc")
		assert(t, nil, err)
		assert(t, "a", d.Method)
		assert(t, "123", d.ID)
		assert(t, "abc", d.Query)
	})

	t.Run("succeeds to extract query after path", func(t *testing.T) {
		d, err := Parse("did:a:123/a/b/c?abc")
		assert(t, nil, err)
		assert(t, "a", d.Method)
		assert(t, "123", d.ID)
		assert(t, "a/b/c", d.Path)
		assert(t, "abc", d.Query)
	})

	t.Run("succeeds to extract fragment after query", func(t *testing.T) {
		d, err := Parse("did:a:123?abc#xyz")
		assert(t, nil, err)
		assert(t, "abc", d.Query)
		assert(t, "xyz", d.Fragment)
	})

	t.Run("succeeds with percent encoded chars in query", func(t *testing.T) {
		d, err := Parse("did:a:123?ab%20c")
		assert(t, nil, err)
		assert(t, "ab%20c", d.Query)
	})

	t.Run("returns error if % in query is not followed by 2 hex chars", func(t *testing.T) {
		dids := []string{
			"did:a:123:456?%",
			"did:a:123:456?%a",
			"did:a:123:456?%!*",
			"did:a:123:456?%A!",
			"did:xyz:pqr?%A!",
			"did:a:123:456?%A%",
		}
		for _, did := range dids {
			_, err := Parse(did)
			assert(t, false, err == nil, "Input: %s", did)
		}
	})

	t.Run("returns error if query has invalid char", func(t *testing.T) {
		_, err := Parse("did:a:123:456?ssss^sss")
		assert(t, false, err == nil)
	})

	t.Run("succeeds to extract fragment", func(t *testing.T) {
		d, err := Parse("did:a:123:456#keys-1")
		assert(t, nil, err)
		assert(t, "keys-1", d.Fragment)
	})

	t.Run("succeeds with percent encoded chars in fragment", func(t *testing.T) {
		d, err := Parse("did:a:123:456#aaaaaa%20a")
		assert(t, nil, err)
		assert(t, "aaaaaa%20a", d.Fragment)
	})

	t.Run("returns error if % in fragment is not followed by 2 hex chars", func(t *testing.T) {
		dids := []string{
			"did:xyz:pqr#%",
			"did:xyz:pqr#%a",
			"did:xyz:pqr#%!*",
			"did:xyz:pqr#%!A",
			"did:xyz:pqr#%A!",
			"did:xyz:pqr#%A%",
		}
		for _, did := range dids {
			_, err := Parse(did)
			assert(t, false, err == nil, "Input: %s", did)
		}
	})

	t.Run("fails if fragment has invalid char", func(t *testing.T) {
		_, err := Parse("did:a:123:456#ssss^sss")
		assert(t, false, err == nil)
	})
}

func assert(t *testing.T, expected interface{}, actual interface{}, args ...interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		argsLength := len(args)
		var message string

		// if only one arg is present, treat it as the message
		if argsLength == 1 {
			message = args[0].(string)
		}

		// if more than one arg is present, treat it as format, args (like Printf)
		if argsLength > 1 {
			message = fmt.Sprintf(args[0].(string), args[1:]...)
		}

		// is message is not empty add some spacing
		if message != "" {
			message = "\t" + message + "\n\n"
		}

		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\tExpected: %#v\n\tActual: %#v\n%s", filepath.Base(file), line, expected, actual, message)
		t.FailNow()
	}
}
