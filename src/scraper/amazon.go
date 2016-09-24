package scraper
import(
  "github.com/PuerkitoBio/goquery"
  //"fmt"
//  "github.com/Jeffail/gabs"
)
type Products struct{
  Name string `json:name`
  Image string `json:img`
  Url string `json:url`
}
func GetProductList(query string) []Products {
doc,_:=goquery.NewDocument("https://www.snapdeal.com/search?keyword="+query)
var Name,Image,Url [4]string
doc.Find(".product-title").Each(func (i int,Selection *goquery.Selection)  {
  if i<4{
    Name[i]=Selection.Text();
  }
 //fmt.Println(Selection.Text())
})
doc.Find("img.product-image").Each(func (i int,Selection *goquery.Selection)  {
  if i<4 {
    d,_:=Selection.Attr("src")
    Image[i]=d
  }
})
doc.Find(".dp-widget-link.noUdLine").Each(func (i int,Selection *goquery.Selection)  {
  if i<4{
    d,_:=Selection.Attr("href")
    Url[i]=d
  }
})
 Result:=make([]Products, 0)
for i := 0;i < 4;i++ {
Result=append(Result,Products{Name[i],Image[i],Url[i]})
}
return Result
}
