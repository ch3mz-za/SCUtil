package scu

import (
	"os"
	"testing"

	"github.com/ch3mz-za/SCUtil/pkg/common"
	"github.com/stretchr/testify/assert"
)

// TODO: Contructor of test directory to add here
func creeatTestingSUDirectory() {
	os.MkdirAll("LIVE/USER/Client/0/Controls", 0666)
	createFiles("LIVE/USER/Client/0/Controls/mappings.xml", "LIVE/USER/Random1.txt",
		"LIVE/Data.p4k", "LIVE/Random1.txt", "LIVE/Random2.txt",
	)
}

func creeatTestingSUAppdateDirectory() {
	os.MkdirAll("SU/USER/Client/0/Controls", 0666)
	createFiles("LIVE/USER/Client/0/Controls/mappings.xml", "LIVE/Data.p4k", "LIVE/Random1.txt", "LIVE/Random2.txt")
}

func createFiles(filePaths ...string) {
	for _, filePath := range filePaths {
		f, err := os.Create(filePath)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func removeTestingDirectory() {
	os.RemoveAll("LIVE")
}

func Test_ClearUserFolder(t *testing.T) {
	defer removeTestingDirectory()
	creeatTestingSUDirectory()

	// with exclusions
	err := ClearUserFolder("LIVE", true)
	assert.NoError(t, err)
	assert.True(t, common.Exists("LIVE/USER/Client/0/Controls/mappings.xml"))
	assert.False(t, common.Exists("LIVE/USER/Random1.txt"))

	// without exclusions
	err = ClearUserFolder("LIVE", false)
	assert.NoError(t, err)
	assert.False(t, common.Exists("LIVE/USER/Client/0/Controls/mappings.xml"))
}

func Test_ClearP4kData(t *testing.T) {
	defer removeTestingDirectory()
	creeatTestingSUDirectory()
	err := ClearP4kData("LIVE")
	assert.NoError(t, err)
	assert.False(t, common.Exists("LIVE/Data.p4k"))
}

func Test_ClearAllButP4kData(t *testing.T) {
	defer removeTestingDirectory()
	creeatTestingSUDirectory()
	err := ClearAllDataExceptP4k("LIVE")
	assert.NoError(t, err)
	assert.True(t, common.Exists("LIVE/Data.p4k"))
}

func Test_ClearStarCitizenAppData(t *testing.T) {
	t.Skip()
	creeatTestingSUAppdateDirectory()
	err := ClearStarCitizenAppData(true)
	assert.NoError(t, err)
}
