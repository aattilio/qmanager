package core

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
)

type Catalog struct {
	OperatingSystems []OperatingSystemMetadata
}

type XMLCatalog struct {
	XMLName xml.Name                  `xml:"catalog"`
	Systems []OperatingSystemMetadata `xml:"os"`
}

type OperatingSystemMetadata struct {
	ID                  string   `xml:"id,attr"`
	Name                string   `xml:"name,attr"`
	Version             string   `xml:"version,attr"`
	Family              string   `xml:"family,attr"`
	Architecture        string   `xml:"arch,attr"`
	Mirrors             []string `xml:"mirrors>mirror"`
	MinRAM              int      `xml:"min_ram_mb"`
	MinVCPUs            int      `xml:"min_vcpus"`
	MinDiskGB           int      `xml:"min_disk_gb"`
	RecommendedDiskBus  string   `xml:"recommended_disk_bus"`
	RecommendedNetModel string   `xml:"recommended_net_model"`
}

func LoadConfigurationCatalogFromDirectory(
	directoryPath string,
) (
	*Catalog,
	error,
) {
	files, err := os.ReadDir(
		directoryPath,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed_to_read_catalog_directory: %w",
			err,
		)
	}

	mergedCatalog := &Catalog{
		OperatingSystems: make(
			[]OperatingSystemMetadata,
			0,
		),
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".xml" {
			continue
		}

		fullPath := filepath.Join(
			directoryPath,
			file.Name(),
		)

		fileData, err := os.ReadFile(
			fullPath,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed_to_read_catalog_file_%s: %w",
				file.Name(),
				err,
			)
		}

		var xmlContent XMLCatalog
		err = xml.Unmarshal(
			fileData,
			&xmlContent,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed_to_unmarshal_xml_%s: %w",
				file.Name(),
				err,
			)
		}

		mergedCatalog.OperatingSystems = append(
			mergedCatalog.OperatingSystems,
			xmlContent.Systems...,
		)
	}

	return mergedCatalog, nil
}
