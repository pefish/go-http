package go_http

import (
	"fmt"
	go_logger "github.com/pefish/go-logger"
	go_test_ "github.com/pefish/go-test"
	"testing"
	"time"
)

func TestHttpClass_interfaceToUrlQuery(t *testing.T) {
	type Test struct {
		A string `json:"a"`
		B string `json:"b"`
		C uint64 `json:"c"`
	}
	result, _ := interfaceToUrlQuery(Test{
		A: `11`,
		B: `22`,
		C: 123,
	})
	go_test_.Equal(t, 16, len(result))
}

func TestHttpClass_Post(t *testing.T) {
	str := `{"side": "test", "chain": "test"}`
	requester := NewHttpRequester(WithLogger(go_logger.Logger.CloneWithLevel("debug")))
	_, body, err := requester.PostForString(RequestParam{
		Url:    `http://www.baidu.com`,
		Params: str,
	})
	go_test_.Equal(t, nil, err)
	go_test_.Equal(t, true, len(body) > 0)
}

func TestHttpClass_Proxy(t *testing.T) {
	//var client http.Client
	//req, err := http.NewRequest("GET", "http://ip.me", nil)
	//test.Equal(t, nil, err)
	//req.Header.Add("User-Agent", "curl/7.64.1")
	//resp, err := client.Do(req)
	//test.Equal(t, nil, err)
	//b, err := ioutil.ReadAll(resp.Body)
	//test.Equal(t, nil, err)
	//fmt.Println(2, string(b))

	requester := NewHttpRequester()
	//requester := NewHttpRequester(WithHttpProxy("http://127.0.0.1:1087"))
	_, body, err := requester.GetForString(RequestParam{
		Url: `https://ip.me`,
		Headers: map[string]interface{}{
			"User-Agent": "curl/7.64.1",
		},
	})
	go_test_.Equal(t, nil, err)
	fmt.Println(1, body)

}

func TestHttpClass_GetForStruct(t *testing.T) {
	var userInfo struct {
		Id      string `json:"id"`
		Message string `json:"message"`
	}
	_, body, err := NewHttpRequester(WithTimeout(10*time.Second)).GetForStruct(RequestParam{
		Url: fmt.Sprintf("https://api.zealy.io/communities/%s/users", "baselendfinance"),
		Params: map[string]interface{}{
			"ethAddress": "0x8f72B6E3DF451F61b4E152de696B2Ba4748cFcc0",
		},
		Headers: map[string]interface{}{
			"x-api-key": "5ac42cpXHE3R9cq7l0Nm-ng78CA",
		},
		BasicAuth: nil,
	}, &userInfo)
	//test.Equal(t, nil, err)
	fmt.Printf("zealy users fetch error. body: %s - %#v\n", string(body), err)
	//if err != nil {
	//	fmt.Printf("zealy users fetch error - %#v\n", err)
	//	//fmt.Println()
	//}
}

func TestHttpClass_PostFormDataForStruct(t *testing.T) {
	var httpResult struct {
		Success     bool   `json:"success"`
		ChallengeTs string `json:"challenge_ts"`
		Hostname    string `json:"hostname"`
	}
	_, _, err := NewHttpRequester(
		WithLogger(go_logger.Logger.CloneWithLevel("debug")),
		WithTimeout(5*time.Second),
	).PostFormDataForStruct(
		RequestParam{
			Url: "https://www.google.com/recaptcha/api/siteverify",
			Params: map[string]interface{}{
				"secret":   "",
				"response": "03AFcWeA79MKSyp6NjxjGNs2EYso_vhr0P685ZSMfcH2V7kXp79SgVm6a6Gh1i6BslENbRI0mkEtawmhaV0_R65ylzAsHk5IL4UWee0OYJ-mT-gBrr1OIKIjWubjJ3tTrZ0uClc_xLDxkfzdd7_PK3QJv1wJZkHrdBujwKkdT8na0PgI2SkcvbRB_ead4qtwo6EyllY6R6LloBY5hxUnmOuVkgIqdLGUU7sS2-aS1chLJEo-Cg7-ldOtOSt6yN2VbZ2J6Xifufpoi1DP-lKjuvsYTe5j3k829renAh5tymwnREWFUkJaqalXLD2Qh3R-9UNfsMBGDA7qeiz81fAg-X_CrgKHU2yYSkVN_yAUuG4svYpXhU19uFLbfXuPn0NqoF5_s7pNQbOv_3hgaFbuRTWLAdh3Bvx5XRf8XOMpX3FPQM7x4gLzECENlBfxcVmbPhOMq9F-rSYIu30lyHzHvu2RWWqtfmhk-_g58_A-uUylb44yHUf1VbSBmHQ-8TA9vXj1JF391wM-OZKUf-ANuegFyBkVb1e7KQHN4MlsMn1IkLgFatURRNkEhJG4NLB6FpXUpMcweAuJJxVaRHt1Evulp0KFmrumCqM7bnQbHD2R8I0v3pZHnNXjcn-7YAT8LiDWqGMK5HZMN7ZijqEbdQzi6KA_i9-2KhTIKft51RGYnIBujSG4SQA5gwN2dDLW1rFb6kyI6Ie1pQ_tqF_hmoRFYjgBtio0OiRuFCnRQuqMD7_aA2MyHjt-TzNkATbjfn1k201ElRxbPxdZhHKv0OYjLMPoUvQ8MwX4FimKlxpR_H1pk-Hz8bAy2IQEcyh4guVjwoSu2Jbk-st3jdTN8MHLZUDZNb5MixAq92v1QN7D7lcNjjAr-2Lngp-fJARazs2jJSOaUljjq4hKmYWJv6wt_zCYgAJwF8MU154IcDI9FDzrUhfhJBibDID2l02ScIUPZXXK1Y615Ja4fekwjWLGok1EJVFmZhdDJb8Zke20uMzgd2MIoF4OPxmgfwA16qqN2h_Vdw7eZw96bYssBx099nTEfa2Pi8RsDGeq_nRCU_JQsVYzJ3XeQAbH4XT_9rlss7i99ieyW1CnEH9TAZ3iW8iwEJMcyeWoeLpcrSVRlgSyl3STkImIjsMMTDRDC2R860Je96plgwd2ynOP7WIosiPybzU5mHQm9TJ6aq5DmpC49GCBKRk-B_npMLCF1i-LY9eT_DDnF-8lNpTqG3-Nks4388HbiwDWTU4JrPrkxzoMc5zWemSGeA91gnQV5EgUFvUDx9fT5ld8A-MoWhMxDjp9VYQE1HOMoDpd6BSE11I5p9AsEL3upPTAAfcSMF3lNma_8sacZGMjct1cH3S5-qYG15VEZ3CnlaoQ9jisQB1ka3H-1YoRAcu7tNM02YFxUQ76abBDoPXTN4A7bl66OJfPnmay8KjzMAhcuGAN07j1hKbnMJv8JJRdek9JlVtNxCyFfk18nv6zh1J-jqe2oNGvvwPXSmR-9PPg-w5xfpVDD62H3zZmHQ68m19AMUQUJjnGjzVz80BlsF8T6gF8UBPPHJ7X5S0Eu04GnfFpbUF8awnLCV7F4JRjC_JxGYx3dikwfxbVm6yqaKxyOGYENLmCkxLAwwShlCB-w38YnFUq3su1_vXkifbnex137CBBGYtNWBhtoyA7w-i80AxDjpqNBzQBMkMZwlcWKF-7Ebj3dX-L06mAFkYTLztH2dKLzG_nqob_m21uebUTPLzlCV9Wmez0-3Ap1SnPcTb2xKXNUI1o3JmHQ1VRIAHu5L-V3-h7dgHsd4PDkOsvjeGba9Z5U_2UlSGg",
			},
		},
		&httpResult,
	)
	go_test_.Equal(t, nil, err)
	fmt.Println(httpResult)
}
