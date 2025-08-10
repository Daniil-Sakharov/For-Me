package hw3

import (
	"fmt"
	"github.com/mailru/easyjson"
	"io"
	"io/ioutil"
	"os"
	"strings"
	// "log"
)

//easyjson:json
type User struct {
	Browsers []string `json:"browsers"`
	Company  string   `json:"company"`
	Country  string   `json:"country"`
	Email    string   `json:"email"`
	Job      string   `json:"job"`
	Name     string   `json:"name"`
	Phone    string   `json:"phone"`
}

func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	uniqueBrowsers := make(map[string]struct{})
	foundUsers := ""

	lines := strings.Split(string(fileContents), "\n")

	var users []User
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		var user User
		//fmt.Printf("%v %v\n", err, line)
		err := easyjson.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	for i, client := range users {
		isAndroid := false
		isMSIE := false
		for _, browser := range client.Browsers {

			if len(browser) == 0 {
				//log.Println("cant cast browsers")
				continue
			}

			if strings.Contains(browser, "Android") {
				isAndroid = true
				uniqueBrowsers[browser] = struct{}{}
				//log.Printf("SLOW New client: %s, first seen: %s", client, client["name"])
			}

			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				uniqueBrowsers[browser] = struct{}{}
				//log.Printf("SLOW New client: %s, first seen: %s", client, client["name"])
			}
		}

		if !(isAndroid && isMSIE) {
			continue
		}

		//log.Println("Android and MSIE client:", client["name"], client["email"])
		email := strings.ReplaceAll(client.Email, "@", " [at] ")
		foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, client.Name, email)
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(uniqueBrowsers))
}
