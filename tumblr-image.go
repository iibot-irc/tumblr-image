package main
import (
  "encoding/xml"
  "net/http"
  "log"
  "flag"
  "fmt"
  "regexp"
)

type Item struct {
  Title string `xml:"title"`
  Description string `xml:"description"`
  Link string `xml:"link"`
  Guid string `xml:"guid"`
  PubDate string `xml:"pubDate"`
  Category []string `xml:"category"`
}

type NoSnippet struct {}

func (ns NoSnippet) Error() string {
  return "No snippet found"
}

var target = flag.String("rss", "", "tumblr rss feed url to poll")
var use_link = flag.Bool("link", false, "whether to emit a link to the post, or directly link the images within")

func GetSnippet(htmlblob string) (string, error) {
  re := regexp.MustCompile("<img src=\"(.*)\"/>");
  results := re.FindStringSubmatch(htmlblob)
  if results == nil || len(results) < 2 {
    return "", NoSnippet{}
  }
  return results[1], nil
}

func main() {
  flag.Parse()
  rss, err := http.Get(*target)
  if err != nil {
    log.Fatal("Error retrieving rss feed: ", err)
  }
  defer rss.Body.Close()

  decoder := xml.NewDecoder(rss.Body)
  for {
    t, _ := decoder.Token()
    if t == nil {
      break
    }
    switch se := t.(type) {
      case xml.StartElement:
        var decoded Item
        if se.Name.Local == "item" {
          err = decoder.DecodeElement(&decoded, &se)
          if err != nil {
            log.Fatal("Error parsing rss xml: ", err)
          }
          if *use_link {
            fmt.Println(decoded.Link);
          } else {
            img, err := GetSnippet(decoded.Description)
            if err != nil {
              log.Fatal("Error grabbing image out of xml blob: ", err)
            }
            fmt.Println(img)
          }
          return;
        }
    }
  }
}
