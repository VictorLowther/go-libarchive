package archive

import (
	"bytes"
	"os"
	"testing"
)

func TestSimpleTarFile(t *testing.T) {
	test_file, err := os.Open("./fixtures/test.tar")
	if err != nil {
		t.Fatalf("Error while reading fixture file %s ", err)
	}

	reader, err := NewReader(test_file)
	if err != nil {
		t.Fatalf("Error on creating Archive from a io.Reader:\n %s", err)
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatalf("Error on reader Close:\n %s", err)
		}
	}()
	entry, err := reader.Next()
	if err != nil {
		t.Fatalf("got error on reader.Next():\n%s", err)
	}
	name := entry.PathName()
	if name != "a" {
		t.Fatalf("got %s expected %s as Name of the first entry", name, "a")
	}

	size := entry.Size()

	b := make([]byte, 512)
	readsize, err := reader.Read(b)
	if err != nil {
		t.Fatalf("got error on reader.Read():\n%s", err)
	}
	if size != 14 {
		t.Fatalf("got %d as expected size of the read but wanted %d",size,14)
	}
	if readsize != 14 {
		t.Fatalf("got %d as size of the read but expected %d", readsize, 14)
	}

	expectedContent := []byte("Sha lalal lal\n")
	if !bytes.Equal((b[:size]), expectedContent) {
		t.Fatalf("The contents:\n [%s] are not the expectedContent:\n [%s]", b[:size], expectedContent)
	}

	_, err = reader.Next()
	if err != ErrArchiveEOF {
		t.Fatalf("Expected EOF on second reader.Next() got err :\n %s", err)
	}
}

func TestComplexTar(t *testing.T) {
	test_file,err := os.Open("./fixtures/test2.tar")
	var expectedContents = map[string]bool{
		"a/": true,
		"c/": true,
		"a/a_file": false,
		"b_file": false,
		"c/c_file":false,
	}
	actualContents := make(map[string]bool)
	if err != nil {
		t.Fatalf("Error while reading fixture file %s ",err)
	}

	archive,err := NewReader(test_file)

	if err != nil {
		t.Fatalf("Error creating an Archive from an io.Reader:\n %s",err)
	}

	defer func() {
		err := archive.Close()
		if err != nil {
			t.Fatalf("Error closing Archvie: %s",err)
		}
	}()
	
	for {
		entry, err := archive.Next()
		if err != nil {
			if err == ErrArchiveEOF {
				break
			}
			t.Fatalf("Error reading archive entry: %s",err)
			return
		}
		actualContents[entry.PathName()]=entry.IsDir()
	}

	for name,dir := range expectedContents {
		if actualContents[name] == dir {
			continue
		}
		t.Fatalf("Unexpected inconsistency in archive!\nExpected:\n%v\nActual:\n%v",
			expectedContents,
			actualContents)
	}
}

func TestTwoReaders(t *testing.T) {
	test_file, err := os.Open("./fixtures/test.tar")
	if err != nil {
		t.Fatalf("Error while reading fixture file %s ", err)
	}

	_, err = NewReader(test_file)

	test_file2, err := os.Open("./fixtures/test2.tar")
	if err != nil {
		t.Fatalf("Error while reading fixture file %s ", err)
	}

	_, err = NewReader(test_file2)
	if err != nil {
		t.Fatalf("Error on creating Archive from a io.Reader:\n %s", err)
	}
}
