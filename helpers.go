package nsscache

import (
	"io"
	"io/ioutil"
	"os"
	"path"
)

// WriteAtomic allows atomic updates to files by first writing to a
// temporary file, setting all parameters on the temporary file and
// renaming it to the desired name.  On most Linux systems this will
// be an atomic action.
func WriteAtomic(filename string, wt io.WriterTo, perm os.FileMode) error {
	dir, name := path.Split(filename)
	f, err := ioutil.TempFile(dir, name)
	if err != nil {
		return err
	}
	_, err = wt.WriteTo(f)
	if err == nil {
		err = f.Sync()
	}
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if permErr := os.Chmod(f.Name(), perm); err == nil {
		err = permErr
	}
	if err == nil {
		err = os.Rename(f.Name(), filename)
	}
	// Any err should result in full cleanup.
	if err != nil {
		os.Remove(f.Name())
	}
	return err
}
