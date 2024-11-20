package unifiaccessclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
)

func TestXX(t *testing.T) {
	client := NewWithHttpClient(
		"192.168.1.1",
		12445,
		"",
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	)

	// new, err := client.CreateUser(schema.UserRequest{FirstName: "Perico", LastName: "Palotes"})
	// if err != nil {
	// 	t.Error(err)
	// }
	// fmt.Printf("%s %s\n", new.FirstName, new.LastName)

	users, err := client.ListUsers()
	if err != nil {
		t.Error(err)
	}

	for _, u := range users {
		fmt.Printf("%s %s %s %s\n", u.Id, u.FirstName, u.LastName, u.Status)
	}

	// u, err := client.GetUser("63f654b9-9d11-4a92-9d18-cc536b74d9e8")
	// if err != nil {
	// 	t.Error(err)
	// }
	// fmt.Println(u.FirstName)

	// status := schema.StatusDeactivated
	// err = client.UpdateUser("9258da27-16e6-48a1-a629-788cd529ae42", schema.UserRequest{Status: &status})
	// if err != nil {
	// 	t.Error(err)
	// }
}
