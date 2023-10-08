package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type SMSLists struct {
	XMLName        string `xml:"SMSList"`
	SenderName     string `xml:"SenderName"`
	ReceiverMSISDN int    `xml:"ReceiverMSISDN"`
	SMSText        string `xml:"SMSText"`
}

type SubmitSMSRequest struct {
	XMLName    string     `xml:"SubmitSMSRequest"`
	AccountId  int        `xml:"AccountId"`
	Password   string     `xml:"Password"`
	SecureHash []byte     `xml:"SecureHash"`
	SMSList    []SMSLists `xml:"SMSList"`
	XMLNS      string     `xml:"xmlns,attr"`
}
type SubmitSMSResponse struct {
	XMLName      xml.Name `xml:"SubmitSMSResponse"`
	ResultStatus string   `xml:"ResultStatus"`
	Description  string   `xml:"Description"`
	SMSStatus    []string `xml:"SMSStatus"`
}

func main() {
	funcsms(550049024, "Vodafone.1", "NTGEGYPT", 201010984336, "how are yoy?")
}

func funcsms(Accountid int, Pass string, Name string, Receiver int, Text string) {
	message := string(Accountid) + Pass + Name + string(Receiver) + Text
	key := "F5B4064ABB0646F9986E154C5AFF0FD7"

	// Convert the key and message to byte arrays
	keyBytes := []byte(key)
	messageBytes := []byte(message)

	// Create a new HMAC hasher using SHA-256 and the key
	hasher := hmac.New(sha256.New, keyBytes)

	// Write the message to the hasher
	hasher.Write(messageBytes)

	// Compute the HMAC
	result := hasher.Sum(nil)

	// Convert the result to a hexadecimal string
	hexResult := hex.EncodeToString(result)
	bk := SubmitSMSRequest{
		XMLNS:      "http://www.edafa.com/web2sms/sms/model/",
		AccountId:  Accountid,
		Password:   Pass,
		SecureHash: []byte(hexResult),
		SMSList: []SMSLists{{
			SenderName:     Name,
			ReceiverMSISDN: Receiver,
			SMSText:        Text,
		}},
	}
	os.Stdout.WriteString(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "\t")

	err := enc.Encode(&bk)
	fmt.Println(os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	out, err := xml.Marshal(&bk)

	resp, err := http.Post("https://e3len.vodafone.com.eg/web2sms/sms/submit/", "application/xml", bytes.NewBuffer(out))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response body in xml
	//fmt.Println(string(body))
	xmlData := string(body)
	var response SubmitSMSResponse
	err = xml.Unmarshal([]byte(xmlData), &response)
	if err != nil {
		fmt.Println("Failed to unmarshal XML response:", err)
		return
	}

	fmt.Println("ResultStatus:", response.ResultStatus)
	fmt.Println("Description:", response.Description)
	fmt.Println("SMSStatus:", response.SMSStatus)
}
