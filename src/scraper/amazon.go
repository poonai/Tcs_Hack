package scraper
import(
  "github.com/PuerkitoBio/goquery"
)
/*
Products structure
Name :-name of the product
Image:-image url of the Product
*/
type Products struct{
  Name string `json:name`
  Image string `json:img`
  Url string `json:url`
}
/*
get the product struct slices based on the query
@params query search string
*/
func GetProductList(query string) []Products {
doc,_:=goquery.NewDocument("https://www.snapdeal.com/search?keyword="+query)
var Name,Image,Url [4]string
doc.Find(".product-title").Each(func (i int,Selection *goquery.Selection)  {
  if i<4{
    Name[i]=Selection.Text();
  }

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
