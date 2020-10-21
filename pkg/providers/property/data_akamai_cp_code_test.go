package property

import (
	"regexp"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/papi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestDSCPCode(t *testing.T) {
	t.Run("match by name", func(t *testing.T) {
		client := &mockpapi{}

		// name provided by fixture is "test cpcode"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_test1", Name: "wrong CP code"},
			{ID: "cpc_test2", Name: "test cpcode", ProductIDs:[]string{"prd_test1", "prd_test2"}},
		}}

		client.On("GetCPCodes",
			mock.Anything, // ctx is irrelevant for this test
			papi.GetCPCodesRequest{ContractID: "ctr_test", GroupID: "grp_test"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSCPCode/match_by_name.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "id", "cpc_test2"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "group", "grp_test"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "contract", "ctr_test"),
					),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("match by name output products", func(t *testing.T) {
		client := &mockpapi{}

		// name provided by fixture is "test cpcode"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_test1", Name: "wrong CP code"},
			{ID: "cpc_test2", Name: "test cpcode", ProductIDs:[]string{"prd_test1", "prd_test2"}},
		}}

		client.On("GetCPCodes",
			mock.Anything, // ctx is irrelevant for this test
			papi.GetCPCodesRequest{ContractID: "ctr_test", GroupID: "grp_test"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSCPCode/match_by_name_output_products.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "id", "cpc_test2"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "group", "grp_test"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "contract", "ctr_test"),
						resource.TestCheckOutput("product1", "prd_test1"),
						resource.TestCheckOutput("product2", "prd_test2"),
					),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("match by full ID", func(t *testing.T) {
		client := &mockpapi{}

		// name provided by fixture is "cpc_test2"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_test1", Name: "wrong CP code"},
			{ID: "cpc_test2", Name: "test cpcode", ProductIDs:[]string{"prd_test1", "prd_test2"}},
		}}

		client.On("GetCPCodes",
			mock.Anything, // ctx is irrelevant for this test
			papi.GetCPCodesRequest{ContractID: "ctr_test", GroupID: "grp_test"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSCPCode/match_by_full_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "id", "cpc_test2"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "name", "test cpcode"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "group", "grp_test"),
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "contract", "ctr_test"),
					),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("match by unprefixed ID", func(t *testing.T) {
		client := &mockpapi{}

		// name provided by fixture is "test2"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_test1", Name: "wrong CP code"},
			{ID: "cpc_test2", Name: "test cpcode", ProductIDs:[]string{"prd_test1", "prd_test2"}},
		}}

		client.On("GetCPCodes",
			mock.Anything, // ctx is irrelevant for this test
			papi.GetCPCodesRequest{ContractID: "ctr_test", GroupID: "grp_test"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSCPCode/match_by_unprefixed_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "id", "cpc_test2"),
					),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("no matches", func(t *testing.T) {
		client := &mockpapi{}

		// name provided by fixture is "test cpcode"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_test1", Name: "wrong CP code"},
			{ID: "cpc_test3", Name: "Also wrong CP code"},
		}}

		client.On("GetCPCodes",
			mock.Anything, // ctx is irrelevant for this test
			papi.GetCPCodesRequest{ContractID: "ctr_test", GroupID: "grp_test"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      loadFixtureString("testdata/TestDSCPCode/match_by_unprefixed_id.tf"),
					ExpectError: regexp.MustCompile(`invalid CP Code`),
				}},
			})
		})

		client.AssertExpectations(t)
	})

	t.Run("ambiguous name", func(t *testing.T) {
		TODO(t, "Should we error out and tell the user to supply an exact CP Code ID to resolve ambiguity?")

		client := &mockpapi{}

		// name provided by fixture is "test cpcode"
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_test1", Name: "wrong CP code"},
			{ID: "cpc_test2", Name: "test cpcode"},
			{ID: "cpc_test3", Name: "test cpcode"},
			{ID: "cpc_test4", Name: "test cpcode"},
			{ID: "cpc_test5", Name: "Also wrong CP code"},
		}}

		client.On("GetCPCodes",
			mock.Anything, // ctx is irrelevant for this test
			papi.GetCPCodesRequest{ContractID: "ctr_test", GroupID: "grp_test"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config:      loadFixtureString("testdata/TestDSCPCode/ambiguous_name.tf"),
					ExpectError: regexp.MustCompile(`ambiguous CP Code`),
				}},
			})
		})
	})

	t.Run("name collides with ID", func(t *testing.T) {
		TODO(t, "Should we error out and tell the user to supply an exact CP Code ID to resolve ambiguity?")

		client := &mockpapi{}

		// name provided by fixture is "cpc_test2", which is an exact ID match but it matches name of "cpc_test1" first
		cpc := papi.CPCodeItems{Items: []papi.CPCode{
			{ID: "cpc_test1", Name: "cpc_test2"},
			{ID: "cpc_test2", Name: "correct CP Code"},
			{ID: "cpc_test3", Name: "wrong CP code"},
		}}

		client.On("GetCPCodes",
			mock.Anything, // ctx is irrelevant for this test
			papi.GetCPCodesRequest{ContractID: "ctr_test", GroupID: "grp_test"},
		).Return(&papi.GetCPCodesResponse{CPCodes: cpc}, nil)

		useClient(client, func() {
			resource.UnitTest(t, resource.TestCase{
				Providers: testAccProviders,
				Steps: []resource.TestStep{{
					Config: loadFixtureString("testdata/TestDSCPCode/name_collides_with_id.tf"),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.akamai_cp_code.test", "id", "cpc_test2"),
					),
				}},
			})
		})

		client.AssertExpectations(t)
	})
}
