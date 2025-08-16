package identity

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/johnoct/a3s/internal/aws/client"
)

type Identity struct {
	Account     string
	UserID      string
	ARN         string
	DisplayName string
}

func GetCallerIdentity(ctx context.Context, awsClient *client.AWSClient) (*Identity, error) {
	stsClient := sts.NewFromConfig(awsClient.Config)
	
	result, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to get caller identity: %w", err)
	}

	identity := &Identity{
		Account: *result.Account,
		UserID:  *result.UserId,
		ARN:     *result.Arn,
	}

	// Extract display name from ARN
	// Examples:
	// arn:aws:iam::123456789012:user/john -> john
	// arn:aws:sts::123456789012:assumed-role/RoleName/SessionName -> RoleName/SessionName
	arnParts := strings.Split(identity.ARN, "/")
	if len(arnParts) > 1 {
		identity.DisplayName = strings.Join(arnParts[1:], "/")
	} else {
		// Fallback to last part of ARN after colons
		arnParts = strings.Split(identity.ARN, ":")
		if len(arnParts) > 5 {
			identity.DisplayName = arnParts[5]
		} else {
			identity.DisplayName = identity.UserID
		}
	}

	return identity, nil
}