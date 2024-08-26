package service

import (
	"github.com/malamsyah/go-skele/internal/model"
	"github.com/malamsyah/go-skele/internal/repository"
)

type ResourceService interface {
	CreateResource(resourceType string, payload string) (*model.Resource, error)
	GetResource(id uint) (*model.Resource, error)
	UpdateResource(id uint, res model.Resource) (*model.Resource, error)
	DeleteResource(id uint) error
	GetResources() ([]model.Resource, error)
}

type resourceService struct {
	repo repository.Repository[model.Resource]
}

func NewResourceService(repo repository.Repository[model.Resource]) ResourceService {
	return &resourceService{
		repo: repo,
	}
}

func (s *resourceService) CreateResource(resourceType string, payload string) (*model.Resource, error) {
	resource := model.Resource{
		ResourceType: resourceType,
		Payload:      payload,
	}

	err := s.repo.Create(&resource)
	if err != nil {
		return nil, err
	}

	return &resource, nil
}

func (s *resourceService) GetResource(id uint) (*model.Resource, error) {
	resource, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *resourceService) UpdateResource(id uint, res model.Resource) (*model.Resource, error) {
	resource, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	resource.ResourceType = res.ResourceType
	resource.Payload = res.Payload
	return s.repo.Update(resource)
}

func (s *resourceService) DeleteResource(id uint) error {
	return s.repo.Delete(id)
}

func (s *resourceService) GetResources() ([]model.Resource, error) {
	return s.repo.FindAll()
}
