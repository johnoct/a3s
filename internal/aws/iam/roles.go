package iam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/johnoct/a3s/internal/aws/client"
)

type RoleService struct {
	client *iam.Client
}

func NewRoleService(awsClient *client.AWSClient) *RoleService {
	return &RoleService{
		client: iam.NewFromConfig(awsClient.Config),
	}
}

type Role struct {
	Name               string
	ARN                string
	CreateDate         time.Time
	Description        string
	MaxSessionDuration int32
	Path               string
	RoleID             string
	Tags               []Tag
	TrustPolicy        string
	LastUsed           *time.Time
	ManagedPolicies    []string
	InlinePolicies     []string
}

type Tag struct {
	Key   string
	Value string
}

func (s *RoleService) ListRoles(ctx context.Context) ([]Role, error) {
	var roles []Role
	paginator := iam.NewListRolesPaginator(s.client, &iam.ListRolesInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list roles: %w", err)
		}

		for _, r := range output.Roles {
			role := Role{
				Name:               *r.RoleName,
				ARN:                *r.Arn,
				CreateDate:         *r.CreateDate,
				Path:               *r.Path,
				RoleID:             *r.RoleId,
				MaxSessionDuration: aws.ToInt32(r.MaxSessionDuration),
			}

			if r.Description != nil {
				role.Description = *r.Description
			}

			if r.AssumeRolePolicyDocument != nil {
				decoded, _ := url.QueryUnescape(*r.AssumeRolePolicyDocument)
				role.TrustPolicy = formatJSON(decoded)
			}

			if r.RoleLastUsed != nil && r.RoleLastUsed.LastUsedDate != nil {
				role.LastUsed = r.RoleLastUsed.LastUsedDate
			}

			roles = append(roles, role)
		}
	}

	return roles, nil
}

func (s *RoleService) GetRoleDetails(ctx context.Context, roleName string) (*Role, error) {
	getRoleOutput, err := s.client.GetRole(ctx, &iam.GetRoleInput{
		RoleName: &roleName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	r := getRoleOutput.Role
	role := &Role{
		Name:               *r.RoleName,
		ARN:                *r.Arn,
		CreateDate:         *r.CreateDate,
		Path:               *r.Path,
		RoleID:             *r.RoleId,
		MaxSessionDuration: aws.ToInt32(r.MaxSessionDuration),
	}

	if r.Description != nil {
		role.Description = *r.Description
	}

	if r.AssumeRolePolicyDocument != nil {
		decoded, _ := url.QueryUnescape(*r.AssumeRolePolicyDocument)
		role.TrustPolicy = formatJSON(decoded)
	}

	if r.RoleLastUsed != nil && r.RoleLastUsed.LastUsedDate != nil {
		role.LastUsed = r.RoleLastUsed.LastUsedDate
	}

	// Get tags
	tags, err := s.client.ListRoleTags(ctx, &iam.ListRoleTagsInput{
		RoleName: &roleName,
	})
	if err == nil {
		for _, t := range tags.Tags {
			role.Tags = append(role.Tags, Tag{
				Key:   *t.Key,
				Value: *t.Value,
			})
		}
	}

	// Get attached managed policies
	managedPolicies, err := s.client.ListAttachedRolePolicies(ctx, &iam.ListAttachedRolePoliciesInput{
		RoleName: &roleName,
	})
	if err == nil {
		for _, p := range managedPolicies.AttachedPolicies {
			role.ManagedPolicies = append(role.ManagedPolicies, *p.PolicyName)
		}
	}

	// Get inline policies
	inlinePolicies, err := s.client.ListRolePolicies(ctx, &iam.ListRolePoliciesInput{
		RoleName: &roleName,
	})
	if err == nil {
		role.InlinePolicies = inlinePolicies.PolicyNames
	}

	return role, nil
}

func (s *RoleService) GetInlinePolicy(ctx context.Context, roleName, policyName string) (string, error) {
	output, err := s.client.GetRolePolicy(ctx, &iam.GetRolePolicyInput{
		RoleName:   &roleName,
		PolicyName: &policyName,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get inline policy: %w", err)
	}

	decoded, _ := url.QueryUnescape(*output.PolicyDocument)
	return formatJSON(decoded), nil
}

func formatJSON(jsonStr string) string {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return jsonStr
	}

	formatted, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return jsonStr
	}

	return string(formatted)
}
