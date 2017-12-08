package fortnox

import (
	"context"
	"fmt"
	pth "github.com/byrnedo/apibase/helpers/pointerhelp"
	"github.com/byrnedo/apibase/helpers/stringhelp"
	"gopkg.in/jarcoal/httpmock.v1"
	"os"
	"testing"
)

var (
	accessToken = os.Getenv("FORTNOX_ACCESS_TOKEN")
	secret      = os.Getenv("FORTNOX_CLIENT_SECRET")
)

func init() {
	if accessToken == "" {
		panic("must give FORTNOX_ACCESS_TOKEN env")
	}
	if secret == "" {
		panic("must give FORTNOX_CLIENT_SECRET env")
	}
}

func addTestOpts() []OptionsFunc {
	return []OptionsFunc{WithAuthOpts(accessToken, secret)}
}

func TestGetAccessToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", DefaultURL,
		httpmock.NewStringResponder(200, `{"Authorization": {"AccessToken": "test"}}`))

	token, err := GetAccessToken(context.Background(), "test", secret, func(opts *AccessTokenOptions) {
		httpmock.ActivateNonDefault(opts.HTTPClient)
	})
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("Token empty")
	}
}

func TestNewFortnoxClient(t *testing.T) {

	c := NewFortnoxClient(WithAuthOpts("token", "secret"), WithURLOpts("url"))

	if c.clientOptions.BaseURL != "url" {
		t.Fatal("Incorrect url")
	}

	if c.clientOptions.AccessToken != "token" {
		t.Fatal("Incorrect token")
	}

	if c.clientOptions.ClientSecret != "secret" {
		t.Fatal("Incorrect secret")
	}

	if c.clientOptions.ContentType != "application/json" {
		t.Fatal("Incorrect content type")
	}
}

func TestGetOrders(t *testing.T) {
	c := NewFortnoxClient(addTestOpts()...)

	r, err := c.ListOrders(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.MetaInformation == nil {
		t.Fatal("Meta was nil")
	}

	if r.Orders == nil {
		t.Fatal("Response was nil")
	}
	//pretty.Print(r)
}

func TestGetOrder(t *testing.T) {
	c := NewFortnoxClient(addTestOpts()...)
	for i := 1; i < 10; i++ {
		_, err := c.GetOrder(context.Background(), fmt.Sprintf("%d", i))
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestGetInvoices(t *testing.T) {
	c := NewFortnoxClient(addTestOpts()...)

	r, err := c.ListInvoices(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.MetaInformation == nil {
		t.Fatal("Meta was nil")
	}

	if r.Invoices == nil {
		t.Fatal("Response was nil")
	}
	//pretty.Print(r)
}

func TestGetInvoice(t *testing.T) {
	c := NewFortnoxClient(addTestOpts()...)
	for i := 1; i < 10; i++ {
		r, err := c.GetInvoice(context.Background(), fmt.Sprintf("%d", i))
		if err != nil {
			t.Fatal(err)
		}
		if r == nil {
			t.Fatal("Response was nil")
		}
	}

}

func TestFortnoxClient_GetCompanySettings(t *testing.T) {
	c := NewFortnoxClient(addTestOpts()...)
	r, err := c.GetCompanySettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if r == nil {
		t.Fatal("Response was nil")
	}

}

func TestFortnoxClient_GetArticles(t *testing.T) {

	c := NewFortnoxClient(addTestOpts()...)

	r, err := c.ListArticles(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.MetaInformation == nil {
		t.Fatal("Meta was nil")
	}

	if r.Articles == nil {
		t.Fatal("Response was nil")
	}

}

func TestFortnoxClient_GetArticle(t *testing.T) {

	c := NewFortnoxClient(addTestOpts()...)

	r, err := c.GetArticle(context.Background(), "10")
	if err != nil {
		t.Fatal(err)
	}
	if r == nil {
		t.Fatal("Response was nil")
	}

}

func TestFortnoxClient_GetLabels(t *testing.T) {

	c := NewFortnoxClient(addTestOpts()...)

	r, err := c.ListLabels(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if r == nil {
		t.Fatal("Response was nil")
	}

}

func TestFortnoxClient_CreateLabel(t *testing.T) {

	c := NewFortnoxClient(addTestOpts()...)

	r, err := c.CreateLabel(context.Background(), "test"+stringhelp.RandString(4))
	if err != nil {
		t.Fatal(err)
	}

	if r == nil {
		t.Fatal("Response was nil")
	}
}

func TestFortnoxClient_CreateOrder(t *testing.T) {
	c := NewFortnoxClient(addTestOpts()...)

	order := &CreateOrder{
		CustomerNumber: pth.StringPtr("1"),
		OrderRows: []*CreateOrderRow{
			{Description: pth.StringPtr("Desc Text")},
		},
	}
	r, err := c.CreateOrder(context.Background(), order)
	if err != nil {
		t.Fatal(err)
	}

	if r == nil {
		t.Fatal("Response was nil")
	}

	if len(r.OrderRows) != 1 {
		t.Fatalf("unexpected number of order rows, expected 1, got %d", len(r.OrderRows))
	}

	row := r.OrderRows[0]
	checkTextRow(row, t)

}

func checkTextRow(row OrderRow, t *testing.T) {
	if row.Description != "Desc Text" {
		t.Fatalf("unexpected description: %s", row.Description)
	}
	// if no article
	if row.AccountNumber != 0 {
		t.Fatalf("unexpected account number: %d", row.AccountNumber)
	}

	if row.CostCenter != "" {
		t.Fatalf("unexpected cost center: %s", row.CostCenter)
	}

}
