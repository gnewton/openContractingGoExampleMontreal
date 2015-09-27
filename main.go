// Glen Newton
// 2015 09 27
// Copyright Glen Newton 2015
// See License (MIT)

package main

import "fmt"
import "net/http"
import "io/ioutil"
import "encoding/json"

type OpenConMtl struct {
	License string `json:"license"`
	Meta    struct {
		Count      int     `json:"count"`
		MaxValue   float64 `json:"max_value"`
		MinValue   float64 `json:"min_value"`
		Pagination struct {
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		} `json:"pagination"`
		TotalValue float64 `json:"total_value"`
	} `json:"meta"`
	PublicationPolicy string `json:"publicationPolicy"`
	PublishedDate     string `json:"publishedDate"`
	Publisher         struct {
		Address struct {
			CountryName   string `json:"countryName"`
			Locality      string `json:"locality"`
			PostalCode    string `json:"postalCode"`
			Region        string `json:"region"`
			StreetAddress string `json:"streetAddress"`
		} `json:"address"`
		ContactPoint struct {
			Email     string `json:"email"`
			FaxNumber string `json:"faxNumber"`
			Name      string `json:"name"`
			Telephone string `json:"telephone"`
			URL       string `json:"url"`
		} `json:"contactPoint"`
		Name string `json:"name"`
	} `json:"publisher"`
	Releases []struct {
		Awards []struct {
			Date  string `json:"date"`
			ID    string `json:"id"`
			Items []struct {
				Description string `json:"description"`
				ID          string `json:"id"`
				Quantity    int    `json:"quantity"`
			} `json:"items"`
			Repartition interface{} `json:"repartition"`
			Suppliers   []struct {
				Address    struct{} `json:"address"`
				Identifier struct {
					ID string `json:"id"`
				} `json:"identifier"`
				Name string `json:"name"`
			} `json:"suppliers"`
			Value struct {
				Amount   float64 `json:"amount"`
				Currency string  `json:"currency"`
			} `json:"value"`
		} `json:"awards"`
		Buyer struct {
			Address struct {
				CountryName   string `json:"countryName"`
				Locality      string `json:"locality"`
				PostalCode    string `json:"postalCode"`
				Region        string `json:"region"`
				StreetAddress string `json:"streetAddress"`
			} `json:"address"`
			Identifier struct {
				ID string `json:"id"`
			} `json:"identifier"`
			Name              string `json:"name"`
			SubOrganisationOf struct {
				Name interface{} `json:"name"`
			} `json:"subOrganisationOf"`
		} `json:"buyer"`
		Date     string   `json:"date"`
		ID       string   `json:"id"`
		Language string   `json:"language"`
		Ocid     string   `json:"ocid"`
		Subject  []string `json:"subject"`
		Tag      string   `json:"tag"`
		Tender   struct {
			ProcurementMethodRationale string `json:"procurementMethodRationale"`
			ProcuringEntity            struct {
				Identifier struct {
					ID string `json:"id"`
				} `json:"identifier"`
				Name string `json:"name"`
			} `json:"procuringEntity"`
			Status string `json:"status"`
		} `json:"tender"`
	} `json:"releases"`
	URI string `json:"uri"`
}

func perror(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	url := "http://ville.montreal.qc.ca/vuesurlescontrats/api/releases?q=mecano&format=json"

	res, err := http.Get(url)
	perror(err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	perror(err)

	var data OpenConMtl
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		switch v := err.(type) {
		case *json.SyntaxError:
			fmt.Println(string(body[v.Offset-40 : v.Offset]))
		}
		perror(err)
	}
	for _, release := range data.Releases {
		fmt.Println("--------------------------------")
		fmt.Printf("Date: %s\n", release.Date)
		fmt.Printf("Subject: %s\n", release.Subject[0])
		fmt.Printf("Language: %s\n", release.Language)

		// Tender
		fmt.Printf("Procurement Method Rationale: %s\n", release.Tender.ProcurementMethodRationale)
		fmt.Printf("Status: %s\n", release.Tender.Status)

		for _, award := range release.Awards {
			fmt.Printf("Value: %f\n", award.Value.Amount)

			// award.Repartition is always null in example URL, so do not know how to deal with...
			//fmt.Printf("Répartition: %f\n", award.Repartition)
			fmt.Printf("No de dossier: %s\n", award.ID)
			for _, supplier := range award.Suppliers {
				fmt.Printf("Fournisseur: %s\n", supplier.Name)
			}
		}
		fmt.Printf("Buyer: %s\n", release.Buyer.Name)
	}
}
