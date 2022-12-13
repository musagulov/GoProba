package main

import (
        "fmt"
        "log"
        "context"
        "encoding/json"
        "net/http"
        "os"
        "golang.org/x/oauth2"
        "google.golang.org/api/option"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/sheets/v4"
        "strconv"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
    // The file token.json stores the user's access and refresh tokens, and is
    // created automatically when the authorization flow completes for the first
    // time.
    tokFile := "token.json"
    tok, err := tokenFromFile(tokFile)
    if err != nil {
            tok = getTokenFromWeb(config)
            saveToken(tokFile, tok)
    }
    return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
    authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    fmt.Printf("Go to the following link in your browser then type the "+
            "authorization code: \n%v\n", authURL)

    var authCode string
    if _, err := fmt.Scan(&authCode); err != nil {
            log.Fatalf("Unable to read authorization code: %v", err)
    }

    tok, err := config.Exchange(context.TODO(), authCode)
    if err != nil {
            log.Fatalf("Unable to retrieve token from web: %v", err)
    }
    return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    if err != nil {
            return nil, err
    }
    defer f.Close()
    tok := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(tok)
    return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
    fmt.Printf("Saving credential file to: %s\n", path)
    f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
            log.Fatalf("Unable to cache oauth token: %v", err)
    }
    defer f.Close()
    json.NewEncoder(f).Encode(token)
}

func main() {

    ctx := context.Background()
    b, err := os.ReadFile("./client_secret_452127709660-k90pvbscne34762n6oihpsaov60586bk.apps.googleusercontent.com.json")
    if err != nil {
            log.Fatalf("Unable to read client secret file: %v", err)
    }

    // If modifying these scopes, delete your previously saved token.json.
    config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
    if err != nil {
            log.Fatalf("Unable to parse client secret file to config: %v", err)
    }
    client := getClient(config)
    svc, err := sheets.NewService(ctx, option.WithHTTPClient(client))
    if err != nil {
        log.Fatal(err)
    }
    lastRow := getLastRow(svc)
    fmt.Println("Last row is %v", lastRow)
    insertRow(lastRow, svc)
}

func getLastRow(svc *sheets.Service) int {
    rowIndex := 0
   for counter := 2; true; counter ++ { 
            ranger := "C" + strconv.Itoa(counter)
            fmt.Println(ranger)
            resp, err := svc.Spreadsheets.Values.Get("1uQujCUJyzA-qjrVQmhteFvZDBE3Vw5KzmWgwEdQTb4c", ranger).Do()
            if err != nil {
                    log.Fatal(err)
            }
            fmt.Println(len(resp.Values))
            if len(resp.Values)>0 {
                fmt.Println(resp.Values[0])
            } else {
                rangerHoriz := "D" + strconv.Itoa(counter)
                resp, err := svc.Spreadsheets.Values.Get("1uQujCUJyzA-qjrVQmhteFvZDBE3Vw5KzmWgwEdQTb4c", rangerHoriz).Do()
                if err != nil {
                        log.Fatal(err)
                }
                if len(resp.Values)==0 {
                    rowIndex = counter
                    break
                } else {
                    continue
                }
            }            
        }
    return rowIndex
}

func insertRow(rowIndex int, svc *sheets.Service) {
    ranger := "sheet1!C" + strconv.Itoa(rowIndex) + ":E" + strconv.Itoa(rowIndex)
    rb := &sheets.ValueRange{
        Range: ranger,
        Values: [][]interface{}{{"Izbasar Musagulov","GO code TBD","I Love My Kids"}},
    }

    valueInputOption :=  "USER_ENTERED"
    _, err := svc.Spreadsheets.Values.Append("1uQujCUJyzA-qjrVQmhteFvZDBE3Vw5KzmWgwEdQTb4c", ranger, rb).ValueInputOption(valueInputOption).Do()
    if err != nil {
        log.Fatal(err)
    }
}