package zip

import (
	"archive/zip"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"io"
	"os"
	"path/filepath"
)

type archiveToFile struct {
	config             *configuration.Config
	filePath           actions.Property
	archiveFileName    actions.Property
	createdArchivePath actions.Property
}

func CreateArchiveToFile(config *configuration.Config) actions.Action {
	return &archiveToFile{
		config: config,
		filePath: actions.Property{
			Name:        "filePath",
			Type:        actions.Text,
			Description: "The path of the directory to be archived",
			DisplayName: "Directory Path",
			Validation:  "required",
		},
		archiveFileName: actions.Property{
			Name:        "archiveFileName",
			Type:        actions.Text,
			Description: "The name of the archive file",
			DisplayName: "Archive File Name",
			Validation:  "required",
		},
		createdArchivePath: actions.Property{
			Name:        "createdArchivePath",
			Type:        actions.Text,
			Description: "The path where the new archive was created",
			DisplayName: "Created Archive Path",
			Validation:  "",
		},
	}
}
func (action *archiveToFile) GetCategoryName() string {
	return "zip"
}

func (action *archiveToFile) Inputs() []actions.Property {
	return []actions.Property{
		action.filePath,
		action.archiveFileName,
	}
}

func (action *archiveToFile) Outputs() []actions.Property {
	return []actions.Property{
		action.createdArchivePath,
	}
}

// Execute performs the archive to file action.
// It takes a Process message as input, extracts the file path and archive file name properties from the message,
// creates a new archive file at the specified path, and adds the files from the specified file path to the archive.
// The resulting archive file path is set as a property in the message.
//
// Parameters:
//   - process: The input Process message containing the file path and archive file name properties.
//
// Returns:
//   - types.Process: The modified Process message with the archive file path property set.
//   - error: An error if any occurred during the execution.
func (action *archiveToFile) Execute(process types.Process) (types.Process, error) {
	filePath, err := action.filePath.GetStringFrom(&process)
	if err != nil {
		return types.Process{}, err
	}
	archiveFileName, err := action.archiveFileName.GetStringFrom(&process)
	if err != nil {
		return types.Process{}, err
	}
	archiveFullPath := filepath.Join(action.config.Folder.TmpFolder, archiveFileName)
	archiveFile, err := os.Create(archiveFullPath)
	if err != nil {
		return types.Process{}, err
	}
	defer archiveFile.Close()

	zipWriter := zip.NewWriter(archiveFile)
	defer zipWriter.Close()

	err = filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only add files to the archive, skip directories
		//if info.IsDir() {
		//	return nil
		//}

		fileToZip, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, fileToZip)
		return err
	})
	if err != nil {
		return process, err
	}

	process.SetString(action.createdArchivePath.Name, archiveFullPath)
	return process, nil
}
