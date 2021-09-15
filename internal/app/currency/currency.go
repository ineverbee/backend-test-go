package currency

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// ChangeCurrency
func ChangeCurrency(s string, n int) (float64, error) {
	if n == 0 {
		return 0, nil
	}

	url := "https://currency-converter5.p.rapidapi.com/currency/convert?format=json&from=RUB&to=" + s + "&amount=" + strconv.Itoa(n)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-host", "currency-converter5.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "9407f19621msh75a4b4562cf2b30p11d087jsnbb4e25148fab")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return 0, err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var unmarsh func(raw []byte) error
	var objmap map[string]json.RawMessage
	unmarsh = func(raw []byte) error {
		err := json.Unmarshal(raw, &objmap)

		if err != nil {
			return err
		}
		return nil
	}
	err = unmarsh(body)
	_ = unmarsh(objmap["rates"])
	_ = unmarsh(objmap[s])
	str := string(objmap["rate_for_amount"])

	if err != nil {
		return 0, err
	}

	f, err := strconv.ParseFloat(strings.Trim(str, "\""), 64)
	return f, err
}
