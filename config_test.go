package gnotes

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	configFile := "testdata/config.ini"

	os.Setenv("MY_ENV", "my-dir")

	c, err := LoadConfig(configFile)
	require.NoError(t, err)

	expected := &Config{
		App: appSettings{
			Editor:  "vim",
			NoteDir: "/home/westley/my-dir/.config/gnotes",
		},
		S3: S3Config{
			Active:    true,
			Bucket:    "gnotes",
			Endpoint:  "https://objects-us-east-1.dream.io",
			Region:    "us-east-1",
			AccessKey: "ACCESS_KEY",
			SecretKey: "SECRET_KEY",
			UserID:    "a8085892-7bf4-11ed-bbd6-a74217c9099d",
			CryptKey:  "ie02kwj1mkslao2jdifie",
		},
	}

	assert.Equal(t, expected, c, "did not parse config correctly")

	// Test opening a invalid file
	c, err = LoadConfig("/foo/bar")
	assert.Error(t, err)
	assert.Nil(t, c)
}
