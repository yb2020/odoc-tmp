package bean

import (
	"github.com/yb2020/odoc/services/doc/model"
)

type CopyFolderTempInfo struct {
	NewFolderId   string                        `json:"newFolderId"`
	OldFolderId   string                        `json:"oldFolderId"`
	OldParentId   string                        `json:"oldParentId"`
	NewParentId   string                        `json:"newParentId"`
	Relations     []model.UserDocFolderRelation `json:"relations"`
	UserDocFolder *model.UserDocFolder          `json:"userDocFolder"`
}
