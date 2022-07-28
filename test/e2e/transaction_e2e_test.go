package e2e_test

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"text/template"

	httpexpect "github.com/gavv/httpexpect/v2"
	"github.com/tunaiku/mobilebanking/test/e2e/setup"
)

type testCase struct {
	desc                 string
	payload              map[string]interface{}
	pathVariables        map[string]interface{}
	responseHTTPStatus   int
	responseBodyExpecter func(*httpexpect.Response)
}

type transactionEndpointTestTable struct {
	testCases  []testCase
	endpoint   string
	httpExpect *httpexpect.Expect
	httpMethod string
}

func runTestsCreateTransaction(t *testing.T, endpoint string, httpMethod string, httpExpect *httpexpect.Expect, desc string,
	payload map[string]interface{}, responseHTTPStatus int, responseBodyExpecter func(*httpexpect.Response)) {

	var pathVariables map[string]interface{}

	tmpl, err := template.New("request").Parse(endpoint)
	if err != nil {
		t.Error(err)
	}
	var output bytes.Buffer
	err = tmpl.Execute(&output, pathVariables)
	if err != nil {
		panic(err)
	}
	requestEndpoint := string(output.Bytes())
	method := strings.ToUpper(httpMethod)
	testCaseName := fmt.Sprintln(method, " ", requestEndpoint, " ", desc)
	t.Run(testCaseName, func(t *testing.T) {

		accessToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmYzU1ZTNhOC1jMGZiLTQwYzctYWI4YS05Y2RhM2ZjYTQwZDQifQ.2oM9B0sTpIlgN-zvDGyrnaNJDiIIU6eIgiko7NxZj2s"
		r := httpExpect.Request(method, requestEndpoint)
		r = r.WithJSON(payload)
		r = r.WithHeader("Authorization", accessToken)

		resp := r.Expect()

		resp.Status(responseHTTPStatus)

		expect := responseBodyExpecter
		if expect != nil {
			expect(resp)
		}
	})
}

func Test_should_created_transaction_when_auth_method_sets_to_pin_and_the_parameter_valid(t *testing.T) {
	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
		endpoint := "/transaction"
		httpMethod := "post"
		httpExpect := e
		desc := " should created transaction when auth_method sets to `pin` and the parameter valid"
		payload := map[string]interface{}{
			"auth_method":         "pin",
			"amount":              3000,
			"transaction_code":    "T001",
			"destination_account": "10002",
		}
		responseHTTPStatus := http.StatusCreated
		responseBodyExpecter := func(resp *httpexpect.Response) {
			resp.JSON().Object().ContainsKey("transaction_id").NotEmpty()
		}
		runTestsCreateTransaction(t, endpoint, httpMethod, httpExpect, desc, payload, responseHTTPStatus, responseBodyExpecter)
	})
}

func Test_should_be_failed_when_the_amount_not_match_the_minimum_transaction_amount(t *testing.T) {
	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
		endpoint := "/transaction"
		httpMethod := "post"
		httpExpect := e
		desc := " should be failed with '400' as http status code and {\"message\":\"amount does not reach the minimum transaction amount\"} when the amount not match the minimum transaction amount"
		payload := map[string]interface{}{
			"auth_method":         "pin",
			"amount":              2000,
			"transaction_code":    "T001",
			"destination_account": "10001",
		}
		responseHTTPStatus := http.StatusBadRequest
		responseBodyExpecter := func(resp *httpexpect.Response) {
			resp.JSON().Object().ValueEqual("message", "amount does not reach the minimum transaction amount")
		}
		runTestsCreateTransaction(t, endpoint, httpMethod, httpExpect, desc, payload, responseHTTPStatus, responseBodyExpecter)
	})
}

func Test_should_be_failed_when_auth_method_not_found(t *testing.T) {
	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
		endpoint := "/transaction"
		httpMethod := "post"
		httpExpect := e
		desc := " should be failed with '400' as http status code and {\"message\":\"unsupported authorization method\"} when auth_method sets to 'password'"
		payload := map[string]interface{}{
			"auth_method":         "password",
			"amount":              3000,
			"transaction_code":    "T001",
			"destination_account": "10001",
		}
		responseHTTPStatus := http.StatusBadRequest
		responseBodyExpecter := func(resp *httpexpect.Response) {
			resp.JSON().Object().ValueEqual("message", "unsupported authorization method")
		}
		runTestsCreateTransaction(t, endpoint, httpMethod, httpExpect, desc, payload, responseHTTPStatus, responseBodyExpecter)
	})
}

func Test_should_be_failed_when_auth_method_sets_to_otp_but_the_user_not_configure_it(t *testing.T) {
	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
		endpoint := "/transaction"
		httpMethod := "post"
		httpExpect := e
		desc := " should be failed with '400' as http status code and {\"message\":\"authorization method not configured\"} when auth_method sets to 'otp' but the user not configure it"
		payload := map[string]interface{}{
			"auth_method":         "otp",
			"amount":              3000,
			"transaction_code":    "T001",
			"destination_account": "10001",
		}
		responseHTTPStatus := http.StatusBadRequest
		responseBodyExpecter := func(resp *httpexpect.Response) {
			resp.JSON().Object().ValueEqual("message", "authorization method not configured")
		}
		runTestsCreateTransaction(t, endpoint, httpMethod, httpExpect, desc, payload, responseHTTPStatus, responseBodyExpecter)
	})
}

func Test_should_be_failed_when_transaction_code_is_not_found(t *testing.T) {
	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
		endpoint := "/transaction"
		httpMethod := "post"
		httpExpect := e
		desc := " should be failed with '400' as http status code and {\"message\":\"transaction code not found\"} when transaction_code sets to 'T003' but that transaction is not found"
		payload := map[string]interface{}{
			"auth_method":         "pin",
			"amount":              3000,
			"transaction_code":    "T003",
			"destination_account": "10001",
		}
		responseHTTPStatus := http.StatusBadRequest
		responseBodyExpecter := func(resp *httpexpect.Response) {
			resp.JSON().Object().ValueEqual("message", "transaction code not found")
		}
		runTestsCreateTransaction(t, endpoint, httpMethod, httpExpect, desc, payload, responseHTTPStatus, responseBodyExpecter)
	})
}

func Test_should_be_failed_when_destination_account_it_not_found(t *testing.T) {
	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
		endpoint := "/transaction"
		httpMethod := "post"
		httpExpect := e
		desc := " should be failed with '400' as http status code and {\"message\":\"destination account not found\"} when destination_account sets to '10003' but the that account not found"
		payload := map[string]interface{}{
			"auth_method":         "pin",
			"amount":              3000,
			"transaction_code":    "T001",
			"destination_account": "10003",
		}
		responseHTTPStatus := http.StatusBadRequest
		responseBodyExpecter := func(resp *httpexpect.Response) {
			resp.JSON().Object().ValueEqual("message", "destination account not found")
		}
		runTestsCreateTransaction(t, endpoint, httpMethod, httpExpect, desc, payload, responseHTTPStatus, responseBodyExpecter)
	})
}

func runTestsVerifyTransaction(t *testing.T, endpoint string, httpMethod string, httpExpect *httpexpect.Expect, desc string,
	pathVariables map[string]interface{}, payload map[string]interface{}, responseHTTPStatus int, responseBodyExpecter func(*httpexpect.Response)) {

	tmpl, err := template.New("request").Parse(endpoint)
	if err != nil {
		t.Error(err)
	}
	var output bytes.Buffer
	err = tmpl.Execute(&output, pathVariables)
	if err != nil {
		panic(err)
	}
	requestEndpoint := string(output.Bytes())
	method := strings.ToUpper(httpMethod)
	testCaseName := fmt.Sprintln(method, " ", requestEndpoint, " ", desc)
	t.Run(testCaseName, func(t *testing.T) {

		accessToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmYzU1ZTNhOC1jMGZiLTQwYzctYWI4YS05Y2RhM2ZjYTQwZDQifQ.2oM9B0sTpIlgN-zvDGyrnaNJDiIIU6eIgiko7NxZj2s"
		r := httpExpect.Request(method, requestEndpoint)
		r = r.WithJSON(payload)
		r = r.WithHeader("Authorization", accessToken)

		resp := r.Expect()

		resp.Status(responseHTTPStatus)

		expect := responseBodyExpecter
		if expect != nil {
			expect(resp)
		}
	})
}

func Test_transaction_should_verified(t *testing.T) {
	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
		endpoint := "/transaction/{{.ID}}/verify"
		httpMethod := "put"
		httpExpect := e
		desc := " transaction should verified with `201-Accepted` when transaction id is valid and the credential is matches "
		pathVariables := map[string]interface{}{
			"ID": "a3289ce9-0c83-4d2f-854f-3a4668c70a71",
		}
		payload := map[string]interface{}{
			"credential": "123456",
		}
		responseHTTPStatus := http.StatusAccepted
		responseBodyExpecter := func(resp *httpexpect.Response) {
			resp.JSON().Object().ContainsKey("transaction_id")
		}

		runTestsVerifyTransaction(t, endpoint, httpMethod, httpExpect, desc, pathVariables,
			payload, responseHTTPStatus, responseBodyExpecter)
	})
}

func Test_transaction_should_be_failed_when_valid_credential_is_invalid(t *testing.T) {
	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
		endpoint := "/transaction/{{.ID}}/verify"
		httpMethod := "put"
		httpExpect := e
		desc := " transaction should be failed with `400-Bad Request` and `{\"message\":\"invalid credential\"}`  when valid credential is invalid "
		pathVariables := map[string]interface{}{
			"ID": "a3289ce9-0c83-4d2f-854f-3a4668c70a71",
		}
		payload := map[string]interface{}{
			"credential": "1234",
		}
		responseHTTPStatus := http.StatusBadRequest
		responseBodyExpecter := func(resp *httpexpect.Response) {
			resp.JSON().Object().ValueEqual("message", "invalid credential")
		}

		runTestsVerifyTransaction(t, endpoint, httpMethod, httpExpect, desc, pathVariables,
			payload, responseHTTPStatus, responseBodyExpecter)
	})
}

func Test_transaction_should_be_failed_when_there_is_no_transaction_with_belong_to_path_variable_id(t *testing.T) {
	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
		endpoint := "/transaction/{{.ID}}/verify"
		httpMethod := "put"
		httpExpect := e
		desc := " transaction should be failed with `404-Not Found` and `{\"message\":\"transaction not found\"}`  when there is no transaction with belong to path variable id "
		pathVariables := map[string]interface{}{
			"ID": "1112",
		}
		payload := map[string]interface{}{
			"credential": "1234",
		}
		responseHTTPStatus := http.StatusNotFound
		responseBodyExpecter := func(resp *httpexpect.Response) {
			resp.JSON().Object().ValueEqual("message", "transaction not found")
		}

		runTestsVerifyTransaction(t, endpoint, httpMethod, httpExpect, desc, pathVariables,
			payload, responseHTTPStatus, responseBodyExpecter)
	})
}

//func (tbl *transactionEndpointTestTable) runTests(t *testing.T) {
//
//	for _, tC := range tbl.testCases {
//		tmpl, err := template.New("request").Parse(tbl.endpoint)
//		if err != nil {
//			t.Error(err)
//		}
//		var output bytes.Buffer
//		err = tmpl.Execute(&output, tC.pathVariables)
//		if err != nil {
//			panic(err)
//		}
//		requestEndpoint := string(output.Bytes())
//		method := strings.ToUpper(tbl.httpMethod)
//		testCaseName := fmt.Sprintln(method, " ", requestEndpoint, " ", tC.desc)
//		t.Run(testCaseName, func(t *testing.T) {
//
//			accessToken := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJmYzU1ZTNhOC1jMGZiLTQwYzctYWI4YS05Y2RhM2ZjYTQwZDQifQ.2oM9B0sTpIlgN-zvDGyrnaNJDiIIU6eIgiko7NxZj2s"
//
//			resp := tbl.httpExpect.Request(method, requestEndpoint).WithJSON(tC.payload).WithHeader("Authorization", accessToken).Expect()
//			resp.Status(tC.responseHTTPStatus)
//			expect := tC.responseBodyExpecter
//			if expect != nil {
//				expect(resp)
//			}
//		})
//	}
//}

//func TestCreateTransaction(t *testing.T) {
//	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
//		testTable := transactionEndpointTestTable{
//			endpoint:   "/transaction",
//			httpMethod: "post",
//			httpExpect: e,
//			testCases: []testCase{
//				{
//					desc: " should created transaction when auth_method sets to `pin` and the parameter valid",
//					payload: map[string]interface{}{
//						"auth_method":         "pin",
//						"amount":              3000,
//						"transaction_code":    "T001",
//						"destination_account": "10002",
//					},
//					responseHTTPStatus: http.StatusCreated,
//					responseBodyExpecter: func(resp *httpexpect.Response) {
//						resp.JSON().Object().ContainsKey("transaction_id").NotEmpty()
//					},
//				},
//				{
//					desc: " should be failed with '400' as http status code and {\"message\":\"amount does not reach the minimum transaction amount\"} when the amount not match the minimum transaction amount",
//					payload: map[string]interface{}{
//						"auth_method":         "pin",
//						"amount":              2000,
//						"transaction_code":    "T001",
//						"destination_account": "10001",
//					},
//					responseHTTPStatus: http.StatusBadRequest,
//					responseBodyExpecter: func(resp *httpexpect.Response) {
//						resp.JSON().Object().ValueEqual("message", "amount does not reach the minimum transaction amount")
//					},
//				},
//				{
//					desc: " should be failed with '400' as http status code and {\"message\":\"unsupported authorization method\"} when auth_method sets to 'password'",
//					payload: map[string]interface{}{
//						"auth_method":         "password",
//						"amount":              3000,
//						"transaction_code":    "T001",
//						"destination_account": "10001",
//					},
//					responseHTTPStatus: http.StatusBadRequest,
//					responseBodyExpecter: func(resp *httpexpect.Response) {
//						resp.JSON().Object().ValueEqual("message", "unsupported authorization method")
//					},
//				},
//				{
//					desc: " should be failed with '400' as http status code and {\"message\":\"authorization method not configured\"} when auth_method sets to 'otp' but the user not configure it",
//					payload: map[string]interface{}{
//						"auth_method":         "otp",
//						"amount":              3000,
//						"transaction_code":    "T001",
//						"destination_account": "10001",
//					},
//					responseHTTPStatus: http.StatusBadRequest,
//					responseBodyExpecter: func(resp *httpexpect.Response) {
//						resp.JSON().Object().ValueEqual("message", "authorization method not configured")
//					},
//				},
//				{
//					desc: " should be failed with '400' as http status code and {\"message\":\"transaction code not found\"} when transaction_code sets to 'T003' but that transaction is not found",
//					payload: map[string]interface{}{
//						"auth_method":         "pin",
//						"amount":              3000,
//						"transaction_code":    "T003",
//						"destination_account": "10001",
//					},
//					responseHTTPStatus: http.StatusBadRequest,
//					responseBodyExpecter: func(resp *httpexpect.Response) {
//						resp.JSON().Object().ValueEqual("message", "transaction code not found")
//					},
//				},
//				{
//					desc: " should be failed with '400' as http status code and {\"message\":\"destination account not found\"} when destination_account sets to '10003' but the that account not found",
//					payload: map[string]interface{}{
//						"auth_method":         "pin",
//						"amount":              3000,
//						"transaction_code":    "T001",
//						"destination_account": "10003",
//					},
//					responseHTTPStatus: http.StatusBadRequest,
//					responseBodyExpecter: func(resp *httpexpect.Response) {
//						resp.JSON().Object().ValueEqual("message", "destination account not found")
//					},
//				},
//			}}
//		testTable.runTests(t)
//	})
//}
//
//func TestVerifyTransaction(t *testing.T) {
//	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
//		testTable := transactionEndpointTestTable{
//			endpoint:   "/transaction/{{.ID}}/verify",
//			httpMethod: "put",
//			httpExpect: e,
//			testCases: []testCase{
//				{
//					desc: " transaction should verified with `201-Accepted` when transaction id is valid and the credential is matches ",
//					pathVariables: map[string]interface{}{
//						"ID": "a3289ce9-0c83-4d2f-854f-3a4668c70a71",
//					},
//					payload: map[string]interface{}{
//						"credential": "123456",
//					},
//					responseHTTPStatus: http.StatusAccepted,
//					responseBodyExpecter: func(resp *httpexpect.Response) {
//						resp.JSON().Object().ContainsKey("transaction_id")
//					},
//				},
//				{
//					desc: " transaction should be failed with `400-Bad Request` and `{\"message\":\"invalid credential\"}`  when valid credential is invalid ",
//					pathVariables: map[string]interface{}{
//						"ID": "a3289ce9-0c83-4d2f-854f-3a4668c70a71",
//					},
//					payload: map[string]interface{}{
//						"credential": "1234",
//					},
//					responseHTTPStatus: http.StatusBadRequest,
//					responseBodyExpecter: func(resp *httpexpect.Response) {
//						resp.JSON().Object().ValueEqual("message", "invalid credential")
//					},
//				},
//				{
//					desc: " transaction should be failed with `400-Bad Request` and `{\"message\":\"verification process already happened\"}`  when the transaction state is not `WaitAuthorization` ",
//					pathVariables: map[string]interface{}{
//						"ID": "a3289ce9-0c83-4d2f-854f-3a4668c70a71",
//					},
//					payload: map[string]interface{}{
//						"credential": "123456",
//					},
//					responseHTTPStatus: http.StatusBadRequest,
//					responseBodyExpecter: func(resp *httpexpect.Response) {
//						resp.JSON().Object().ValueEqual("message", "verification process already happened")
//					},
//				},
//				{
//					desc: " transaction should be failed with `404-Not Found` and `{\"message\":\"transaction not found\"}`  when there is no transaction with belong to path variable id ",
//					pathVariables: map[string]interface{}{
//						"ID": "1112",
//					},
//					payload: map[string]interface{}{
//						"credential": "1234",
//					},
//					responseHTTPStatus: http.StatusNotFound,
//					responseBodyExpecter: func(resp *httpexpect.Response) {
//						resp.JSON().Object().ValueEqual("message", "transaction not found")
//					},
//				},
//			},
//		}
//		testTable.runTests(t)
//	})
//}
//
//func Test_transaction_should_be_failed_when_the_transaction_state_is_not_waitAuthorization(t *testing.T) {
//	setup.InvokeHttpTest(t, func(e *httpexpect.Expect) {
//		endpoint := "/transaction/{{.ID}}/verify"
//		httpMethod := "put"
//		httpExpect := e
//		desc := " transaction should be failed with `400-Bad Request` and `{\"message\":\"verification process already happened\"}`  when the transaction state is not `WaitAuthorization` "
//		pathVariables := map[string]interface{}{
//			"ID": "a3289ce9-0c83-4d2f-854f-3a4668c70a71",
//		}
//		payload := map[string]interface{}{
//			"credential": "123456",
//		}
//		responseHTTPStatus := http.StatusBadRequest
//		responseBodyExpecter := func(resp *httpexpect.Response) {
//			resp.JSON().Object().ValueEqual("message", "verification process already happened")
//		}
//
//		runTestsVerifyTransaction(t, endpoint, httpMethod, httpExpect, desc, pathVariables,
//			payload, responseHTTPStatus, responseBodyExpecter)
//	})
//}
