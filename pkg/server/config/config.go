// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/daytonaio/daytona/pkg/types"

	log "github.com/sirupsen/logrus"
)

func GetConfig() (*types.ServerConfig, error) {
	configFilePath, err := configFilePath()
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		return nil, errors.New("config file does not exist")
	}

	if err != nil {
		return nil, err
	}

	var c types.ServerConfig
	configContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(configContent, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func configFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

func Save(c *types.ServerConfig) error {
	configFilePath, err := configFilePath()
	if err != nil {
		return err
	}

	configContent, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(configFilePath), 0700)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, configContent, 0600)
	if err != nil {
		return err
	}

	return nil
}

func GetConfigDir() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userConfigDir, "daytona", "server"), nil
}

func GetWorkspaceLogsDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "logs"), nil
}

func GetWorkspaceLogFilePath(workspaceId string) (string, error) {
	projectLogsDir, err := GetWorkspaceLogsDir()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(projectLogsDir, workspaceId, "log")

	err = os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func DeleteWorkspaceLogs(workspaceId string) error {
	logsDir, err := GetWorkspaceLogsDir()
	if err != nil {
		return err
	}

	workspaceLogsDir := filepath.Join(logsDir, workspaceId)

	_, err = os.Stat(workspaceLogsDir)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	return os.RemoveAll(workspaceLogsDir)
}

func GetProjectLogFilePath(workspaceId string, projectId string) (string, error) {
	projectLogsDir, err := GetWorkspaceLogsDir()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(projectLogsDir, workspaceId, projectId, "log")

	err = os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func init() {
	_, err := GetConfig()
	if err == nil {
		return
	}

	c, err := getDefaultConfig()
	if err != nil {
		log.Fatal("failed to get default config")
	}

	err = Save(c)
	if err != nil {
		log.Fatal("failed to save default config file")
	}
}
