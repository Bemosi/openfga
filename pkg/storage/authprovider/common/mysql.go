package common

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/go-sql-driver/mysql"
)

func (i *iam) Connect(ctx context.Context) (driver.Conn, error) {
	token, isNew, err := i.getRDSAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	if !isNew {
		return i.baseConnector.Connect(ctx)
	}
	i.mysqlCfg.Passwd = token
	i.baseConnector, err = mysql.NewConnector(i.mysqlCfg)
	if err != nil {
		return nil, err
	}
	return i.baseConnector.Connect(ctx)
}

func (i *iam) Driver() driver.Driver {
	return i.baseConnector.Driver()
}

func open() (*sql.DB, error) {
	// private static readonly string __queryString = "?api-version=2019-08-01&resource=https://vault.azure.net";
	// private static readonly Uri __defaultIMDSRequestUri = new Uri(
	// 	baseUri: new Uri("http://169.254.169.254"),
	// 	relativeUri: $"metadata/identity/oauth2/token{__queryString}");
	deadline := time.Now().Add(1500 * time.Millisecond)
	ctx, cancelCtx := context.WithDeadline(nil, deadline)
	defer cancelCtx()
	var scopes = []string{"https://ossrdbms-aad.database.windows.net/.default"}

	defaultCredential, err := azidentity.NewDefaultAzureCredential(nil)
	defaultCredential.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: scopes,
	})
	return sql.OpenDB(i), nil
}
