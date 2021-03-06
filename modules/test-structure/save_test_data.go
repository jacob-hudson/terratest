package test_structure

import (
	"testing"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/packer"
	"github.com/gruntwork-io/terratest/modules/logger"
	"fmt"
	"path/filepath"
	"encoding/json"
	"os"
	"io/ioutil"
	"strings"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/aws"
)

// Serialize and save TerraformOptions into the given folder. This allows you to create TerraformOptions during setup
// and to reuse that TerraformOptions later during validation and teardown.
func SaveTerraformOptions(t *testing.T, testFolder string, terraformOptions *terraform.Options) {
	SaveTestData(t, formatTerraformOptionsPath(testFolder), terraformOptions)
}

// Load and unserialize TerraformOptions from the given folder. This allows you to reuse a TerraformOptions that was
// created during an earlier setup step in later validation and teardown steps.
func LoadTerraformOptions(t *testing.T, testFolder string) *terraform.Options {
	var terraformOptions terraform.Options
	LoadTestData(t, formatTerraformOptionsPath(testFolder), &terraformOptions)
	return &terraformOptions
}

// Format a path to save TerraformOptions in the given folder
func formatTerraformOptionsPath(testFolder string) string {
	return FormatTestDataPath(testFolder, "TerraformOptions.json")
}

// Serialize and save PackerOptions into the given folder. This allows you to create PackerOptions during setup
// and to reuse that PackerOptions later during validation and teardown.
func SavePackerOptions(t *testing.T, testFolder string, packerOptions *packer.Options) {
	SaveTestData(t, formatPackerOptionsPath(testFolder), packerOptions)
}

// Load and unserialize PackerOptions from the given folder. This allows you to reuse a PackerOptions that was
// created during an earlier setup step in later validation and teardown steps.
func LoadPackerOptions(t *testing.T, testFolder string) *packer.Options {
	var packerOptions packer.Options
	LoadTestData(t, formatPackerOptionsPath(testFolder), &packerOptions)
	return &packerOptions
}

// Format a path to save PackerOptions in the given folder
func formatPackerOptionsPath(testFolder string) string {
	return FormatTestDataPath(testFolder, "PackerOptions.json")
}

// Serialize and save an Ec2KeyPair into the given folder. This allows you to create an Ec2KeyPair during setup
// and to reuse that Ec2KeyPair later during validation and teardown.
func SaveEc2KeyPair(t *testing.T, testFolder string, keyPair *aws.Ec2Keypair) {
	SaveTestData(t, formatEc2KeyPairPath(testFolder), keyPair)
}

// Load and unserialize an Ec2KeyPair from the given folder. This allows you to reuse an Ec2KeyPair that was
// created during an earlier setup step in later validation and teardown steps.
func LoadEc2KeyPair(t *testing.T, testFolder string) *aws.Ec2Keypair {
	var keyPair aws.Ec2Keypair
	LoadTestData(t, formatEc2KeyPairPath(testFolder), &keyPair)
	return &keyPair
}

// Format a path to save an Ec2KeyPair in the given folder
func formatEc2KeyPairPath(testFolder string) string {
	return FormatTestDataPath(testFolder, "Ec2KeyPair.json")
}

// Serialize and save a uniquely named string value into the given folder. This allows you to create one or more string
// values during one stage -- each with a unique name -- and to reuse those values during later stages.
func SaveString(t *testing.T, testFolder string, name string, val string) {
	path := formatNamedTestDataPath(testFolder, name)
	SaveTestData(t, path, val)
}

// Load and unserialize a uniquely named string value from the given folder. This allows you to reuse one or more string
// values that were created during an earlier setup step in later steps.
func LoadString(t *testing.T, testFolder string, name string) string {
	var val string
	LoadTestData(t, formatNamedTestDataPath(testFolder, name), &val)
	return val
}

// Save a uniquely named int value into the given folder. This allows you to create one or more int
// values during one stage -- each with a unique name -- and to reuse those values during later stages.
func SaveInt(t *testing.T, testFolder string, name string, val int) {
	path := formatNamedTestDataPath(testFolder, name)
	SaveTestData(t, path, val)
}

// Load a uniquely named int value from the given folder. This allows you to reuse one or more int
// values that were created during an earlier setup step in later steps.
func LoadInt(t *testing.T, testFolder string, name string) int {
	var val int
	LoadTestData(t, formatNamedTestDataPath(testFolder, name), &val)
	return val
}

// Serialize and save an AMI ID into the given folder. This allows you to build an AMI during setup and to reuse that
// AMI later during validation and teardown.
func SaveAmiId(t *testing.T, testFolder string, amiId string) {
	SaveString(t, testFolder, "AMI", amiId)
}

// Load and unserialize an AMI ID from the given folder. This allows you to reuse an AMI  that was created during an
// earlier setup step in later validation and teardown steps.
func LoadAmiId(t *testing.T, testFolder string) string {
	return LoadString(t, testFolder, "AMI")
}

// Format a path to save an arbitrary named value in the given folder
func formatNamedTestDataPath(testFolder string, name string) string {
	filename := fmt.Sprintf("%s.json", name)
	return FormatTestDataPath(testFolder, filename)
}

// Format a path to save test data
func FormatTestDataPath(testFolder string, filename string) string {
	return filepath.Join(testFolder, ".test-data", filename)
}

// Serialize and save a value used at test time to the given path. This allows you to create some sort of test data
// (e.g., TerraformOptions) during setup and to reuse this data later during validation and teardown.
func SaveTestData(t *testing.T, path string, value interface{}) {
	logger.Logf(t, "Storing test data in %s so it can be reused later", path)

	if IsTestDataPresent(t, path) {
		logger.Logf(t, "[WARNING] The named test data at path %s is non-empty. Save operation will overwrite existing value with \"%v\".\n.", path, value)
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("Failed to convert value %s to JSON: %v", path, err)
	}

	t.Logf("Marshalled JSON: %s", string(bytes))

	parentDir := filepath.Dir(path)
	if err := os.MkdirAll(parentDir, 0777); err != nil {
		t.Fatalf("Failed to create folder %s: %v", parentDir, err)
	}

	if err := ioutil.WriteFile(path, bytes, 0644); err != nil {
		t.Fatalf("Failed to save value %s: %v", path, err)
	}
}

// Load and unserialize a value stored at the given path. The value should be a pointer to a struct into which the
// value will be deserialized. This allows you to reuse some sort of test data (e.g., TerraformOptions) from earlier
// setup steps in later validation and teardown steps.
func LoadTestData(t *testing.T, path string, value interface{}) {
	logger.Logf(t, "Loading test data from %s", path)

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to load value from %s: %v", path, err)
	}

	if err := json.Unmarshal(bytes, value); err != nil {
		t.Fatalf("Failed to parse JSON for value %s: %v", path, err)
	}
}

// Return true if a file exists at $path and the test data there is non-empty.
func IsTestDataPresent(t *testing.T, path string) bool {
	bytes, err := ioutil.ReadFile(path)
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		return false
	} else if err != nil {
		t.Fatalf("Failed to load test data from %s due to unexpected error: %v", path, err)
	}

	if isEmptyJson(t, bytes) {
		return false
	}

	return true
}

// Return true if the given bytes are empty, or in a valid JSON format that can reasonably be considered empty.
// The types used are based on the type possibilities listed at https://golang.org/src/encoding/json/decode.go?s=4062:4110#L51
func isEmptyJson(t *testing.T, bytes []byte) bool {
	var value interface{}

	if len(bytes) == 0 {
		return true
	}

	if err := json.Unmarshal(bytes, &value); err != nil {
		t.Fatalf("Failed to parse JSON while testing whether it is empty: %v", err)
	}

	if value == nil {
		return true
	}

	valueBool, ok := value.(bool)
	if ok && ! valueBool {
		return true
	}

	valueFloat64, ok := value.(float64)
	if ok && valueFloat64 == 0 {
		return true
	}

	valueString, ok := value.(string)
	if ok && valueString == "" {
		return true
	}

	valueSlice, ok := value.([]interface{})
	if ok && len(valueSlice) == 0 {
		return true
	}

	valueMap, ok := value.(map[string]interface{})
	if ok && len(valueMap) == 0 {
		return true
	}

	return false
}

// Clean up the test data at the given path
func CleanupTestData(t *testing.T, path string) {
	if files.FileExists(path) {
		logger.Logf(t, "Cleaning up test data from %s", path)
		if err := os.Remove(path); err != nil {
			t.Fatalf("Failed to clean up file at %s: %v", path, err)
		}
	} else {
		logger.Logf(t, "%s does not exist. Nothing to cleanup.", path)
	}
}
