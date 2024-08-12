package file

import "os"

func GetWorkingDirectory() (wd string, err error) {
	// Get the absolute path of the current working directory
	wd, err = os.Getwd()
	if err != nil {
		return "", err
	}
	return wd, nil
}
