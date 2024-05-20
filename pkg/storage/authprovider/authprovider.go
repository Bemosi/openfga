package authprovider

import (
	"fmt"

	"github.com/openfga/openfga/pkg/storage/authprovider/azure"
	"github.com/openfga/openfga/pkg/storage/authprovider/common"
	"github.com/openfga/openfga/pkg/storage/authprovider/password"
)

func New(method string) (common.StorageAuthProvider, error) {
	switch method {
	case "azure_managed_identity":
		return &azure.ManagedIdentityAuthProvider{}, nil
	case "password":
		fallthrough
	case "":
		return &password.PasswordAuthProvider{}, nil
	default:
		return nil, fmt.Errorf("storage auth provider '%s' is unsupported", method)
	}
}
