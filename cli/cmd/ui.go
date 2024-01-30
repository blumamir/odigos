package cmd

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/keyval-dev/odigos/cli/cmd/resources"
	"github.com/keyval-dev/odigos/cli/pkg/kube"
	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
)

const (
	defaultPort   = 3000
	uiDownloadUrl = "https://github.com/keyval-dev/odigos/releases/download/v%s/ui_%s_%s_%s.tar.gz"
)

// uiCmd represents the ui command
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Start the Odigos UI",
	Long:  `Start the Odigos UI. This will start a web server that will serve the UI`,
	Run: func(cmd *cobra.Command, args []string) {
		// get all flags as slice of strings
		var flags []string
		cmd.Flags().Visit(func(f *pflag.Flag) {
			flags = append(flags, fmt.Sprintf("--%s=%s", f.Name, f.Value))
		})

		ctx := cmd.Context()
		client, err := kube.CreateClient(cmd)
		if err != nil {
			kube.PrintClientErrorAndExit(err)
		}

		ns, err := resources.GetOdigosNamespace(client, ctx)
		if err != nil {
			if !resources.IsErrNoOdigosNamespaceFound(err) {
				fmt.Printf("\033[31mERROR\033[0m Cannot start Odigos UI. Failed to check if Odigos is already installed: %s\n", err)
			} else {
				fmt.Printf("\033[31mERROR\033[0m Odigos is not installed in your kubernetes cluster. Run 'odigos install' or switch your k8s context to use a different cluster \n")
			}
			os.Exit(1)
		}
		flags = append(flags, fmt.Sprintf("--namespace=%s", ns))

		clusterVersion, _ := GetOdigosVersionInCluster(ctx, client, ns)
		binaryPath, binaryDir := getOdigosUiBinaryPath()
		shouldReplaceBinary := shouldDownloadNewUiBinary(binaryPath, clusterVersion)

		if shouldReplaceBinary {
			err := downloadOdigosUIVersion(runtime.GOARCH, runtime.GOOS, binaryDir, clusterVersion)
			if err != nil {
				fmt.Printf("Error downloading UI binary: %v\n", err)
				os.Exit(1)
			}
		}

		// execute UI binary with all flags and stream output
		process := exec.Command(binaryPath, flags...)
		process.Stdout = os.Stdout
		process.Stderr = os.Stderr
		err = process.Run()
		if err != nil {
			fmt.Printf("Error starting UI: %v\n", err)
			os.Exit(1)
		}
	},
}

func shouldDownloadNewUiBinary(binaryPath string, odigosClusterVersion string) bool {

	_, err := os.Stat(binaryPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Could not find UI binary, downloading it\n")
			return true
		} else {
			fmt.Printf("Error checking for UI binary: %v\n", err)
			os.Exit(1)
		}
	}

	// check if the binary matches the version of odigos in the cluster
	// run the binary with --version flag and check if the version matches
	cmd := exec.Command(binaryPath, "--version")
	outputBytes, err := cmd.Output()
	if err != nil {
		fmt.Printf("Unable to extract current odigos ui version: %v, installing new version\n", err)
		return true
	}
	output := string(outputBytes)
	re := regexp.MustCompile(`v\d+\.\d+\.\d+`)
	version := re.FindString(output)
	if version != odigosClusterVersion {
		fmt.Printf("UI binary version (%s) does not match Odigos version (%s), downloading new UI binary\n", version, odigosClusterVersion)
		return true
	}

	return false
}

func getOdigosUiBinaryPath() (binaryPath, binaryDir string) {
	// Look for binary named odigos-ui in the same directory as the current binary
	// and execute it.
	// currentBinaryPath, err := os.Executable()
	currentBinaryPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting current binary path: %v\n", err)
		os.Exit(1)
	}

	if strings.HasPrefix(currentBinaryPath, "/ko-app") {
		// we are running in a docker container, no permission to write in the current directory
		// use /tmp as the binary directory
		binaryDir = "/tmp"
		binaryPath = filepath.Join(binaryDir, "odigos-ui")
		return
	}

	binaryDir = filepath.Dir(currentBinaryPath)
	binaryPath = filepath.Join(binaryDir, "odigos-ui")
	return
}

func downloadOdigosUIVersion(arch string, os string, currentDir string, odigosVersion string) error {

	if odigosVersion == "" {
		latestVersion, err := GetLatestReleaseVersion()
		if err != nil {
			return err
		}
		odigosVersion = latestVersion
	}

	fmt.Printf("Downloading version %s of Odigos UI ...\n", odigosVersion)
	// if the version starts with "v", remove it
	odigosVersion = strings.TrimPrefix(odigosVersion, "v")
	url := getDownloadUrl(os, arch, odigosVersion)
	return downloadAndExtractTarGz(url, currentDir)
}

func downloadAndExtractTarGz(url string, dir string) error {
	// Step 1: Download the tar.gz file
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer response.Body.Close()

	// Step 2: Create a new gzip reader
	gzipReader, err := gzip.NewReader(response.Body)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()

	// Step 3: Create a new tar reader
	tarReader := tar.NewReader(gzipReader)

	// Step 4: Extract files one by one
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // Reached the end of the tar archive
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %v", err)
		}

		// Step 5: Create the file or directory
		targetPath := filepath.Join(dir, header.Name)
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
			continue
		}

		file, err := os.Create(targetPath)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}

		// Step 6: Copy file contents from tar to the target file
		if _, err := io.Copy(file, tarReader); err != nil {
			file.Close()
			return fmt.Errorf("failed to extract file: %v", err)
		}

		file.Close()
	}

	return os.Chmod(filepath.Join(dir, "odigos-ui"), 0755)
}

func getDownloadUrl(os string, arch string, version string) string {
	return fmt.Sprintf(uiDownloadUrl, version, version, os, arch)
}

type Release struct {
	TagName string `json:"tag_name"`
}

func GetLatestReleaseVersion() (string, error) {
	url := "https://api.github.com/repos/keyval-dev/odigos/releases/latest"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch latest release: %s", resp.Status)
	}

	var release Release
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return "", err
	}

	ver, _ := strings.CutPrefix(release.TagName, "v")
	return ver, nil
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().String("address", "localhost", "address to listen on")
	uiCmd.Flags().Int("port", defaultPort, "port to listen on")
	uiCmd.Flags().Bool("debug", false, "enable debug mode")
}
