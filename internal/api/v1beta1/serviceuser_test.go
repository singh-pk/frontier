package v1beta1

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/raystack/frontier/core/organization"
	"github.com/raystack/frontier/core/permission"
	"github.com/raystack/frontier/core/project"
	"github.com/raystack/frontier/core/relation"
	"github.com/raystack/frontier/core/resource"
	"github.com/raystack/frontier/core/serviceuser"
	"github.com/raystack/frontier/internal/api/v1beta1/mocks"
	"github.com/raystack/frontier/internal/bootstrap/schema"
	"github.com/raystack/frontier/pkg/metadata"
	frontierv1beta1 "github.com/raystack/frontier/proto/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	su1 = serviceuser.ServiceUser{
		ID:        "1",
		Title:     "1",
		OrgID:     "1",
		State:     "1",
		Metadata:  metadata.Metadata{},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	su2 = serviceuser.ServiceUser{
		ID:    "2",
		OrgID: "2",
		Title: "2",
		State: "2",
		Metadata: metadata.Metadata{
			"key": "value",
		},
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	su1PB = &frontierv1beta1.ServiceUser{
		Id:    "1",
		Title: "1",
		OrgId: "1",
		State: "1",
		Metadata: &structpb.Struct{
			Fields: map[string]*structpb.Value{},
		},
		CreatedAt: timestamppb.New(time.Time{}),
		UpdatedAt: timestamppb.New(time.Time{}),
	}
	su2PB = &frontierv1beta1.ServiceUser{
		Id:    "2",
		Title: "2",
		OrgId: "2",
		State: "2",
		Metadata: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"key": {
					Kind: &structpb.Value_StringValue{
						StringValue: "value",
					},
				},
			},
		},
		CreatedAt: timestamppb.New(time.Time{}),
		UpdatedAt: timestamppb.New(time.Time{}),
	}
)

func TestHandler_ListServiveUsers(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(su *mocks.ServiceUserService)
		request *frontierv1beta1.ListServiceUsersRequest
		want    *frontierv1beta1.ListServiceUsersResponse
		wantErr error
	}{
		{
			name:    "should return internal server error when list service user service returns error",
			request: &frontierv1beta1.ListServiceUsersRequest{},
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().List(mock.AnythingOfType("context.backgroundCtx"), serviceuser.Filter{
					OrgID: "",
					State: "",
				}).Return(nil, errors.New("test error"))
			},
			want:    nil,
			wantErr: errors.New("test error"),
		},
		{
			name: "Test List Service Users",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().List(mock.AnythingOfType("context.backgroundCtx"), serviceuser.Filter{
					OrgID: "",
					State: "",
				}).Return([]serviceuser.ServiceUser{su1, su2}, nil)
			},
			request: &frontierv1beta1.ListServiceUsersRequest{},
			want: &frontierv1beta1.ListServiceUsersResponse{
				Serviceusers: []*frontierv1beta1.ServiceUser{
					su1PB,
					su2PB,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServiveUserSvc := new(mocks.ServiceUserService)
			if tt.setup != nil {
				tt.setup(mockServiveUserSvc)
			}
			h := Handler{
				serviceUserService: mockServiveUserSvc,
			}
			got, err := h.ListServiceUsers(context.Background(), &frontierv1beta1.ListServiceUsersRequest{})
			assert.EqualValues(t, tt.want, got)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_GetServiceUser(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(su *mocks.ServiceUserService)
		request *frontierv1beta1.GetServiceUserRequest
		want    *frontierv1beta1.GetServiceUserResponse
		wantErr error
	}{
		{
			name: "should return internal server error when get service user service returns error",
			request: &frontierv1beta1.GetServiceUserRequest{
				Id: "1",
			},
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().Get(mock.AnythingOfType("context.backgroundCtx"), "1").Return(serviceuser.ServiceUser{}, errors.New("test error"))
			},
			want:    nil,
			wantErr: errors.New("test error"),
		},
		{
			name: "should return not found error when service user is not found",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().Get(mock.AnythingOfType("context.backgroundCtx"), "1").Return(serviceuser.ServiceUser{}, serviceuser.ErrNotExist)
			},
			request: &frontierv1beta1.GetServiceUserRequest{
				Id: "1",
			},
			want:    nil,
			wantErr: grpcServiceUserNotFound,
		},
		{
			name: "should return service user",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().Get(mock.AnythingOfType("context.backgroundCtx"), "1").Return(su1, nil)
			},
			request: &frontierv1beta1.GetServiceUserRequest{
				Id: "1",
			},
			want: &frontierv1beta1.GetServiceUserResponse{
				Serviceuser: su1PB,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServiveUserSvc := new(mocks.ServiceUserService)
			if tt.setup != nil {
				tt.setup(mockServiveUserSvc)
			}
			h := Handler{
				serviceUserService: mockServiveUserSvc,
			}
			got, err := h.GetServiceUser(context.Background(), tt.request)
			assert.EqualValues(t, tt.want, got)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_CreateServiceUser(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(su *mocks.ServiceUserService)
		request *frontierv1beta1.CreateServiceUserRequest
		want    *frontierv1beta1.CreateServiceUserResponse
		wantErr error
	}{
		{
			name: "should return internal server error when create service user service returns error",
			request: &frontierv1beta1.CreateServiceUserRequest{
				Body: &frontierv1beta1.ServiceUserRequestBody{
					Title:    su1PB.GetTitle(),
					Metadata: su1PB.GetMetadata(),
				},
				OrgId: su1PB.GetOrgId(),
			},
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().Create(mock.AnythingOfType("context.backgroundCtx"), serviceuser.ServiceUser{
					Title:    su1.Title,
					Metadata: su1.Metadata,
					OrgID:    su1.OrgID,
				}).Return(serviceuser.ServiceUser{}, errors.New("test error"))
			},
			want:    nil,
			wantErr: errors.New("test error"),
		},
		{
			name: "should return service user",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().Create(mock.AnythingOfType("context.backgroundCtx"), serviceuser.ServiceUser{
					Title:    su1.Title,
					Metadata: su1.Metadata,
					OrgID:    su1.OrgID,
				}).Return(su1, nil)
			},
			request: &frontierv1beta1.CreateServiceUserRequest{
				Body: &frontierv1beta1.ServiceUserRequestBody{
					Title:    su1PB.GetTitle(),
					Metadata: su1PB.GetMetadata(),
				},
				OrgId: su1PB.GetOrgId(),
			},
			want: &frontierv1beta1.CreateServiceUserResponse{
				Serviceuser: su1PB,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServiveUserSvc := new(mocks.ServiceUserService)
			if tt.setup != nil {
				tt.setup(mockServiveUserSvc)
			}
			h := Handler{
				serviceUserService: mockServiveUserSvc,
			}
			got, err := h.CreateServiceUser(context.Background(), tt.request)
			assert.EqualValues(t, tt.want, got)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_DeleteServiceUser(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(su *mocks.ServiceUserService)
		request *frontierv1beta1.DeleteServiceUserRequest
		want    *frontierv1beta1.DeleteServiceUserResponse
		wantErr error
	}{
		{
			name: "should return internal server error when delete service user service returns error",
			request: &frontierv1beta1.DeleteServiceUserRequest{
				Id: "1",
			},
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().Delete(mock.AnythingOfType("context.backgroundCtx"), "1").Return(errors.New("test error"))
			},
			want:    nil,
			wantErr: errors.New("test error"),
		},
		{
			name: "should return not found error when service user is not found",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().Delete(mock.AnythingOfType("context.backgroundCtx"), "1").Return(serviceuser.ErrNotExist)
			},
			request: &frontierv1beta1.DeleteServiceUserRequest{
				Id: "1",
			},
			want:    nil,
			wantErr: grpcServiceUserNotFound,
		},
		{
			name: "should return service user",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().Delete(mock.AnythingOfType("context.backgroundCtx"), "1").Return(nil)
			},
			request: &frontierv1beta1.DeleteServiceUserRequest{
				Id: "1",
			},
			want:    &frontierv1beta1.DeleteServiceUserResponse{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServiveUserSvc := new(mocks.ServiceUserService)
			if tt.setup != nil {
				tt.setup(mockServiveUserSvc)
			}
			h := Handler{
				serviceUserService: mockServiveUserSvc,
			}
			got, err := h.DeleteServiceUser(context.Background(), tt.request)
			assert.EqualValues(t, tt.want, got)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_CreateServiceUserJWK(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(su *mocks.ServiceUserService)
		request *frontierv1beta1.CreateServiceUserJWKRequest
		want    *frontierv1beta1.CreateServiceUserJWKResponse
		wantErr error
	}{
		{
			name: "should return internal server error when create service user key service returns error",
			request: &frontierv1beta1.CreateServiceUserJWKRequest{
				Id:    "1",
				Title: "title",
			},
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().CreateKey(mock.AnythingOfType("context.backgroundCtx"), serviceuser.Credential{
					Title:         "title",
					ServiceUserID: "1",
				}).Return(serviceuser.Credential{}, errors.New("test error"))
			},
			want:    nil,
			wantErr: errors.New("test error"),
		},
		{
			name: "should return not found error when service user is not found",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().CreateKey(mock.AnythingOfType("context.backgroundCtx"), serviceuser.Credential{
					Title:         "title",
					ServiceUserID: "1",
				}).Return(serviceuser.Credential{}, serviceuser.ErrNotExist)
			},
			request: &frontierv1beta1.CreateServiceUserJWKRequest{
				Id:    "1",
				Title: "title",
			},
			want:    nil,
			wantErr: grpcServiceUserNotFound,
		},
		{
			name: "should return service user key",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().CreateKey(mock.AnythingOfType("context.backgroundCtx"), serviceuser.Credential{
					ServiceUserID: "1",
					Title:         "title",
				}).Return(suKey1PB, nil)
			},
			request: &frontierv1beta1.CreateServiceUserJWKRequest{
				Id:    "1",
				Title: "title",
			},
			want: &frontierv1beta1.CreateServiceUserJWKResponse{
				Key: &Key1PB,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServiveUserSvc := new(mocks.ServiceUserService)
			if tt.setup != nil {
				tt.setup(mockServiveUserSvc)
			}
			h := Handler{
				serviceUserService: mockServiveUserSvc,
			}
			got, err := h.CreateServiceUserJWK(context.Background(), tt.request)
			assert.EqualValues(t, tt.want, got)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

// Todo
var suKey1PB = serviceuser.Credential{
	ID:            "1",
	ServiceUserID: "1",
	Title:         "title",
	SecretHash:    "hash",
	PublicKey:     jwk.NewSet(),
	PrivateKey:    []byte("private"),
}
var Key1PB = frontierv1beta1.KeyCredential{
	Type:        "sv_rsa",
	PrincipalId: "1",
	PrivateKey:  "private",
	Kid:         "1",
}

func TestHandler_ListServiceUserJWKs(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(su *mocks.ServiceUserService)
		request *frontierv1beta1.ListServiceUserJWKsRequest
		want    *frontierv1beta1.ListServiceUserJWKsResponse
		wantErr error
	}{
		{
			name: "should return internal server error when list service user keys service returns error",
			request: &frontierv1beta1.ListServiceUserJWKsRequest{
				Id: "1",
			},
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().ListKeys(mock.AnythingOfType("context.backgroundCtx"), "1").Return(nil, errors.New("test error"))
			},
			want:    nil,
			wantErr: errors.New("test error"),
		},
		{
			name: "should return not found error when service user is not found",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().ListKeys(mock.AnythingOfType("context.backgroundCtx"), "1").Return(nil, serviceuser.ErrNotExist)
			},
			request: &frontierv1beta1.ListServiceUserJWKsRequest{
				Id: "1",
			},
			want:    nil,
			wantErr: grpcServiceUserNotFound,
		},
		{
			name: "should return service user keys",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().ListKeys(mock.AnythingOfType("context.backgroundCtx"), "1").Return([]serviceuser.Credential{suKey1PB}, nil)
			},
			request: &frontierv1beta1.ListServiceUserJWKsRequest{
				Id: "1",
			},
			want: &frontierv1beta1.ListServiceUserJWKsResponse{
				Keys: []*frontierv1beta1.ServiceUserJWK{
					{
						Id:          "1",
						Title:       "title",
						PrincipalId: "1",
						PublicKey:   "{\"keys\":[]}",
						CreatedAt:   timestamppb.New(time.Time{}),
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServiveUserSvc := new(mocks.ServiceUserService)
			if tt.setup != nil {
				tt.setup(mockServiveUserSvc)
			}
			h := Handler{
				serviceUserService: mockServiveUserSvc,
			}
			got, err := h.ListServiceUserJWKs(context.Background(), tt.request)
			assert.EqualValues(t, tt.want, got)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_GetServiceUserJWK(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(su *mocks.ServiceUserService)
		request *frontierv1beta1.GetServiceUserJWKRequest
		want    *frontierv1beta1.GetServiceUserJWKResponse
		wantErr error
	}{
		{
			name: "should return internal server error when get service user key service returns error",
			request: &frontierv1beta1.GetServiceUserJWKRequest{
				Id:    "1",
				KeyId: "1",
			},
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().GetKey(mock.AnythingOfType("context.backgroundCtx"), "1").Return(serviceuser.Credential{}, errors.New("test error"))
			},
			want:    nil,
			wantErr: errors.New("test error"),
		},
		{
			name: "should return not found error when service user is not found",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().GetKey(mock.AnythingOfType("context.backgroundCtx"), "1").Return(serviceuser.Credential{}, serviceuser.ErrCredNotExist)
			},
			request: &frontierv1beta1.GetServiceUserJWKRequest{
				Id:    "1",
				KeyId: "1",
			},
			want:    nil,
			wantErr: grpcSvcUserCredNotFound,
		},
		{
			name: "should return service user key",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().GetKey(mock.AnythingOfType("context.backgroundCtx"), "1").Return(suKey1PB, nil)
			},
			request: &frontierv1beta1.GetServiceUserJWKRequest{
				Id:    "1",
				KeyId: "1",
			},
			want: &frontierv1beta1.GetServiceUserJWKResponse{
				Keys: []*frontierv1beta1.JSONWebKey{
					// {
					// 	// Kid: "1",
					// 	// Kty: "RSA",
					// 	// N:   "null",
					// 	// E:   "null",
					// 	// Alg: "RS256",
					// },
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServiveUserSvc := new(mocks.ServiceUserService)
			if tt.setup != nil {
				tt.setup(mockServiveUserSvc)
			}
			h := Handler{
				serviceUserService: mockServiveUserSvc,
			}
			got, err := h.GetServiceUserJWK(context.Background(), tt.request)
			assert.EqualValues(t, tt.want, got)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_DeleteServiceUserJWK(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(su *mocks.ServiceUserService)
		request *frontierv1beta1.DeleteServiceUserJWKRequest
		want    *frontierv1beta1.DeleteServiceUserJWKResponse
		wantErr error
	}{
		{
			name: "should return internal server error when delete service user key service returns error",
			request: &frontierv1beta1.DeleteServiceUserJWKRequest{
				Id:    "1",
				KeyId: "1",
			},
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().DeleteKey(mock.AnythingOfType("context.backgroundCtx"), "1").Return(errors.New("test error"))
			},
			want:    nil,
			wantErr: errors.New("test error"),
		},
		{
			name: "should return not found error when service user is not found",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().DeleteKey(mock.AnythingOfType("context.backgroundCtx"), "1").Return(serviceuser.ErrCredNotExist)
			},
			request: &frontierv1beta1.DeleteServiceUserJWKRequest{
				Id:    "1",
				KeyId: "1",
			},
			want:    nil,
			wantErr: grpcSvcUserCredNotFound,
		},
		{
			name: "should return service user key",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().DeleteKey(mock.AnythingOfType("context.backgroundCtx"), "1").Return(nil)
			},
			request: &frontierv1beta1.DeleteServiceUserJWKRequest{
				Id:    "1",
				KeyId: "1",
			},
			want:    &frontierv1beta1.DeleteServiceUserJWKResponse{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServiveUserSvc := new(mocks.ServiceUserService)
			if tt.setup != nil {
				tt.setup(mockServiveUserSvc)
			}
			h := Handler{
				serviceUserService: mockServiveUserSvc,
			}
			got, err := h.DeleteServiceUserJWK(context.Background(), tt.request)
			assert.EqualValues(t, tt.want, got)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_DeleteServiceUserCredential(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(su *mocks.ServiceUserService)
		request *frontierv1beta1.DeleteServiceUserCredentialRequest
		want    *frontierv1beta1.DeleteServiceUserCredentialResponse
		wantErr error
	}{
		{
			name: "should return internal server error when delete service user secret service returns error",
			request: &frontierv1beta1.DeleteServiceUserCredentialRequest{
				Id:       "1",
				SecretId: "1",
			},
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().DeleteSecret(mock.AnythingOfType("context.backgroundCtx"), "1").Return(errors.New("test error"))
			},
			want:    nil,
			wantErr: errors.New("test error"),
		},
		{
			name: "should return service user secret",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().DeleteSecret(mock.AnythingOfType("context.backgroundCtx"), "1").Return(nil)
			},
			request: &frontierv1beta1.DeleteServiceUserCredentialRequest{
				Id:       "1",
				SecretId: "1",
			},
			want:    &frontierv1beta1.DeleteServiceUserCredentialResponse{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServiveUserSvc := new(mocks.ServiceUserService)
			if tt.setup != nil {
				tt.setup(mockServiveUserSvc)
			}
			h := Handler{
				serviceUserService: mockServiveUserSvc,
			}
			got, err := h.DeleteServiceUserCredential(context.Background(), tt.request)
			assert.EqualValues(t, tt.want, got)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_CreateServiceUserCredential(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(su *mocks.ServiceUserService)
		request *frontierv1beta1.CreateServiceUserCredentialRequest
		want    *frontierv1beta1.CreateServiceUserCredentialResponse
		wantErr error
	}{
		{
			name: "should return internal server error when create service user secret service returns error",
			request: &frontierv1beta1.CreateServiceUserCredentialRequest{
				Id:    "1",
				Title: "title",
			},
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().CreateSecret(mock.AnythingOfType("context.backgroundCtx"), serviceuser.Credential{
					// ID:            "1",
					Title:         "title",
					ServiceUserID: "1",
				}).Return(serviceuser.Secret{
					ID:        "1",
					Value:     "value",
					CreatedAt: time.Now(),
				}, errors.New("test error"))
			},
			want:    nil,
			wantErr: errors.New("test error"),
		},
		{
			name: "should return service user secret",
			setup: func(su *mocks.ServiceUserService) {
				su.EXPECT().CreateSecret(mock.AnythingOfType("context.backgroundCtx"), serviceuser.Credential{
					Title:         "title",
					ServiceUserID: "1",
				}).Return(serviceuser.Secret{
					ID:        "1",
					Value:     "value",
					CreatedAt: time.Time{},
				}, nil)
			},
			request: &frontierv1beta1.CreateServiceUserCredentialRequest{
				Id:    "1",
				Title: "title",
			},
			want: &frontierv1beta1.CreateServiceUserCredentialResponse{
				Secret: &frontierv1beta1.SecretCredential{
					Id:        "1",
					Secret:    "value",
					CreatedAt: timestamppb.New(time.Time{}),
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServiveUserSvc := new(mocks.ServiceUserService)
			if tt.setup != nil {
				tt.setup(mockServiveUserSvc)
			}
			h := Handler{
				serviceUserService: mockServiveUserSvc,
			}
			got, err := h.CreateServiceUserCredential(context.Background(), tt.request)
			assert.EqualValues(t, tt.want, got)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestHandler_ListServiceUserProjects(t *testing.T) {
	testProjectMap := map[string]project.Project{
		"ab657ae7-8c9e-45eb-9862-dd9ceb6d5c71": {
			ID:   "ab657ae7-8c9e-45eb-9862-dd9ceb6d5c71",
			Name: "prj-1",
			Metadata: metadata.Metadata{
				"email": "org1@org1.com",
			},
			Organization: organization.Organization{
				ID: testOrgID,
			},
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		"c7772c63-fca4-4c7c-bf93-c8f85115de4b": {
			ID:   "c7772c63-fca4-4c7c-bf93-c8f85115de4b",
			Name: "prj-2",
			Metadata: metadata.Metadata{
				"email": "org1@org2.com",
			},
			Organization: organization.Organization{
				ID: testOrgID,
			},
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
	}

	tests := []struct {
		name    string
		setup   func(projSvc *mocks.ProjectService, permSvc *mocks.PermissionService, resourceSvc *mocks.ResourceService)
		request *frontierv1beta1.ListServiceUserProjectsRequest
		want    *frontierv1beta1.ListServiceUserProjectsResponse
		wantErr error
	}{
		{
			name: "should return internal server error when list service user project returns error",
			request: &frontierv1beta1.ListServiceUserProjectsRequest{
				Id: "1",
			},
			setup: func(projSvc *mocks.ProjectService, permSvc *mocks.PermissionService, resourceSvc *mocks.ResourceService) {
				projSvc.EXPECT().ListByUser(mock.AnythingOfType("context.backgroundCtx"), "1", schema.ServiceUserPrincipal, project.Filter{}).Return(nil, errors.New("test error"))
			},
			want:    nil,
			wantErr: errors.New("test error"),
		},
		{
			name: "should return project list when there is no error",
			request: &frontierv1beta1.ListServiceUserProjectsRequest{
				Id: "1",
			},
			setup: func(projSvc *mocks.ProjectService, permSvc *mocks.PermissionService, resourceSvc *mocks.ResourceService) {
				var projects []project.Project

				for _, projectID := range testProjectIDList {
					projects = append(projects, testProjectMap[projectID])
				}
				projSvc.EXPECT().ListByUser(mock.AnythingOfType("context.backgroundCtx"), "1", schema.ServiceUserPrincipal, project.Filter{}).Return(projects, nil)
			},
			want: &frontierv1beta1.ListServiceUserProjectsResponse{
				Projects: []*frontierv1beta1.Project{{
					Id:   "ab657ae7-8c9e-45eb-9862-dd9ceb6d5c71",
					Name: "prj-1",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"email": structpb.NewStringValue("org1@org1.com"),
						},
					},
					OrgId:     "9f256f86-31a3-11ec-8d3d-0242ac130003",
					CreatedAt: timestamppb.New(time.Time{}),
					UpdatedAt: timestamppb.New(time.Time{}),
				},
					{
						Id:   "c7772c63-fca4-4c7c-bf93-c8f85115de4b",
						Name: "prj-2",
						Metadata: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"email": structpb.NewStringValue("org1@org2.com"),
							},
						},
						OrgId:     "9f256f86-31a3-11ec-8d3d-0242ac130003",
						CreatedAt: timestamppb.New(time.Time{}),
						UpdatedAt: timestamppb.New(time.Time{}),
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "should return project list with access pairs if withPermission is passed",
			request: &frontierv1beta1.ListServiceUserProjectsRequest{
				Id:              "1",
				WithPermissions: []string{"get"},
			},
			setup: func(projSvc *mocks.ProjectService, permSvc *mocks.PermissionService, resourceSvc *mocks.ResourceService) {
				var projects []project.Project

				for _, projectID := range testProjectIDList {
					projects = append(projects, testProjectMap[projectID])
				}

				ctx := mock.AnythingOfType("context.backgroundCtx")
				projSvc.EXPECT().ListByUser(ctx, "1", schema.ServiceUserPrincipal, project.Filter{}).Return(projects, nil)

				permSvc.EXPECT().Get(ctx, "app/project:get").Return(
					permission.Permission{
						ID:          uuid.New().String(),
						Name:        "get",
						NamespaceID: "app/project",
						Metadata:    map[string]any{},
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
					}, nil)

				resourceSvc.EXPECT().BatchCheck(ctx, []resource.Check{
					{
						Object: relation.Object{
							ID:        "ab657ae7-8c9e-45eb-9862-dd9ceb6d5c71",
							Namespace: "app/project",
						},
						Permission: "get",
					},
					{
						Object: relation.Object{
							ID:        "c7772c63-fca4-4c7c-bf93-c8f85115de4b",
							Namespace: "app/project",
						},
						Permission: "get",
					},
				}).Return([]relation.CheckPair{
					{
						Relation: relation.Relation{
							Object: relation.Object{
								ID:        "ab657ae7-8c9e-45eb-9862-dd9ceb6d5c71",
								Namespace: "app/project",
							},
							RelationName: "get",
						},
						Status: true,
					},
					{
						Relation: relation.Relation{
							Object: relation.Object{
								ID:        "c7772c63-fca4-4c7c-bf93-c8f85115de4b",
								Namespace: "app/project",
							},
							RelationName: "get",
						},
						Status: true,
					},
				}, nil)
			},
			want: &frontierv1beta1.ListServiceUserProjectsResponse{
				Projects: []*frontierv1beta1.Project{{
					Id:   "ab657ae7-8c9e-45eb-9862-dd9ceb6d5c71",
					Name: "prj-1",
					Metadata: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"email": structpb.NewStringValue("org1@org1.com"),
						},
					},
					OrgId:     "9f256f86-31a3-11ec-8d3d-0242ac130003",
					CreatedAt: timestamppb.New(time.Time{}),
					UpdatedAt: timestamppb.New(time.Time{}),
				},
					{
						Id:   "c7772c63-fca4-4c7c-bf93-c8f85115de4b",
						Name: "prj-2",
						Metadata: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"email": structpb.NewStringValue("org1@org2.com"),
							},
						},
						OrgId:     "9f256f86-31a3-11ec-8d3d-0242ac130003",
						CreatedAt: timestamppb.New(time.Time{}),
						UpdatedAt: timestamppb.New(time.Time{}),
					},
				},
				AccessPairs: []*frontierv1beta1.ListServiceUserProjectsResponse_AccessPair{
					{
						ProjectId:   "ab657ae7-8c9e-45eb-9862-dd9ceb6d5c71",
						Permissions: []string{"get"},
					},
					{
						ProjectId:   "c7772c63-fca4-4c7c-bf93-c8f85115de4b",
						Permissions: []string{"get"},
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProjectSvc := new(mocks.ProjectService)
			mockPermssionSvc := new(mocks.PermissionService)
			mockResourceSvc := new(mocks.ResourceService)

			if tt.setup != nil {
				tt.setup(mockProjectSvc, mockPermssionSvc, mockResourceSvc)
			}
			h := Handler{
				projectService:    mockProjectSvc,
				permissionService: mockPermssionSvc,
				resourceService:   mockResourceSvc,
			}
			got, err := h.ListServiceUserProjects(context.Background(), tt.request)
			assert.EqualValues(t, tt.want, got)
			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}
