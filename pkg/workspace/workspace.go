package workspace

import (
	"context"

	"signature.service/pkg/storage"
)

type WorkspaceFactoryInterface interface {
	NewWorksapce(context.Context) *Workspace
}

type WorkspaceFactory struct {
	storage storage.Storage
}

func NewWorkspaceFactory(storage storage.Storage) *WorkspaceFactory {
	return &WorkspaceFactory{
		storage: storage,
	}
}

func (wf *WorkspaceFactory) NewWorksapce(ctx context.Context) *Workspace {
	return NewWorkspace(wf.storage)
}

type Workspace struct {
	storage storage.Storage
}

func NewWorkspace(storage storage.Storage) *Workspace {
	return &Workspace{storage: storage}
}
