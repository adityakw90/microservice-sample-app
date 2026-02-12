package mapper

import (
	"fmt"
	"google.golang.org/protobuf/types/known/structpb"

	authpb "github.com/adityakw90/service-user-proto/gen/go/auth"
	"github.com/adityakw90/microservice-sample-app/api-gateway/internal/core/domain/model"
)

// TokensFromProto converts gRPC Token to domain Tokens.
func TokensFromProto(token *authpb.Token) *model.Tokens {
	if token == nil {
		return nil
	}
	return &model.Tokens{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
}

// TokenClaimsFromProto converts gRPC validate response to domain TokenClaims.
func TokenClaimsFromProto(resp *authpb.ValidateTokenResponse) *model.TokenClaims {
	if resp == nil {
		return nil
	}
	claims := &model.TokenClaims{
		UID:            resp.Uid,
		Identifier:     resp.Identifier,
		IdentifierType: resp.IdentifierType,
	}
	if resp.Claims != nil {
		claims.Claims = resp.Claims.AsMap()
	}
	return claims
}

// AttributesToProto converts a map[string]any to protobuf Struct.
func AttributesToProto(attrs map[string]any) (*structpb.Struct, error) {
	if attrs == nil {
		return &structpb.Struct{Fields: make(map[string]*structpb.Value)}, nil
	}

	fields := make(map[string]*structpb.Value)
	for k, v := range attrs {
		val, err := convertValue(v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert attribute %s: %w", k, err)
		}
		fields[k] = val
	}
	return &structpb.Struct{Fields: fields}, nil
}

// convertValue converts any value to protobuf Value.
func convertValue(v any) (*structpb.Value, error) {
	switch val := v.(type) {
	case string:
		return &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: val}}, nil
	case bool:
		return &structpb.Value{Kind: &structpb.Value_BoolValue{BoolValue: val}}, nil
	case float64:
		return &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: val}}, nil
	case int:
		return &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: float64(val)}}, nil
	case int64:
		return &structpb.Value{Kind: &structpb.Value_NumberValue{NumberValue: float64(val)}}, nil
	case map[string]any:
		return convertMap(val)
	default:
		return &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: fmt.Sprintf("%v", val)}}, nil
	}
}

// convertMap converts a nested map to protobuf StructValue.
func convertMap(m map[string]any) (*structpb.Value, error) {
	fields := make(map[string]*structpb.Value)
	for k, v := range m {
		val, err := convertValue(v)
		if err != nil {
			return nil, err
		}
		fields[k] = val
	}
	return &structpb.Value{Kind: &structpb.Value_StructValue{StructValue: &structpb.Struct{Fields: fields}}}, nil
}
