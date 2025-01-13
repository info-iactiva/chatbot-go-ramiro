package tools

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"langtools/utils"
)

func UploadFileToSharePoint(filePath string, fileName string) error {
	siteURL := os.Getenv("SPAAuth_SITEURL")

	if siteURL == "" {
		return fmt.Errorf("SharePoint site URL is missing. Please set SPAAuth_SITEURL in your environment")
	}

	accessToken, err := GetAccessToken()
	if err != nil {
		return fmt.Errorf("error getting access token: %w", err)
	}

	// Leer el archivo que queremos subir
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Construir la URL para la subida del archivo
	relativeURL := "/sites/DESSAMUEBLES-iActiva/Reportes%20chatbot"
	uploadURL := fmt.Sprintf("%s/_api/web/GetFolderByServerRelativeUrl('%s')/Files/add(url='%s',overwrite=true)", siteURL, relativeURL, fileName)

	// Crear la solicitud HTTP
	req, err := http.NewRequest("POST", uploadURL, bytes.NewReader(fileData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	req.Header.Set("Accept", "application/json;odata=verbose")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error uploading file to SharePoint: %w", err)
	}
	defer resp.Body.Close()

	utils.Info("File uploaded successfully to SharePoint")

	return nil
}
