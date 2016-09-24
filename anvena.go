package main
import	"github.com/kataras/iris"
import "github.com/PuerkitoBio/goquery"
import "errors"
import "strings"
import "weather"
import "fmt"
import "scraper"
//import "strconv"
type hotel struct {
  Name string `json:"name"`
  Address string `json:"address"`
}

func HotelByFood(food string,place string)([]string,[]string,error){
  doc,err := goquery.NewDocument("https://www.zomato.com/"+place+"/restaurants?q="+food)
  if err != nil {
   return nil,nil,errors.New("something went wrong")
  }
  result:=make([]string,0)
  url:=make([]string,0)
  resultch:=make(chan []string, 0)
  urlch:=make(chan []string, 0)
  length:=doc.Find(".result-title").Length()
  go doc.Find(".result-title").Each(func (i int , s *goquery.Selection)  {

       result=append(result,s.Text())
       if length==len(result){
        resultch<-result
      }
  })
  go doc.Find(".feat-img").Each(func (i int,s *goquery.Selection)  {

  d1,_:=s.Attr("data-original")
  url=append(url,d1)
  if length==len(url){
    urlch<-url
  }
  })
  return <-urlch,<-resultch,nil
}
func main() {
iris.Get("foodad/:place/:food",func (ctx *iris.Context)  {
   hotelimg,hotelname,err:=HotelByFood(ctx.Param("food"),ctx.Param("place"))
   if err==nil {
     result:=make([]hotel, 0)
     for i := 0; i < len(hotelname); i++ {

       result=append(result,hotel{strings.Split(hotelname[i],"\n")[0],strings.Split(hotelimg[i],"?")[0]})
     }
     ctx.JSON(200,result)
   }else {
     ctx.Write("something went wrong")
   }
})
iris.Get("/weatherad",func (ctx *iris.Context)  {
  humidity,weather:=weather.GetWeather(ctx.URLParam("lat"),ctx.URLParam("lng"))
  if weather=="Rain"{
       ctx.JSON(200,scraper.GetProductList("rain"))
  }else  if humidity>70{
    ctx.JSON(200,scraper.GetProductList("cotton tshirt"))
    //ctx.Write("Asdf")
  }else if humidity<70{
    ctx.JSON(200,scraper.GetProductList("trekking"))
  }
})
iris.Get("/workout",func (ctx *iris.Context)  {
  switch ctx.URLParam("q") {
  case "running":
        ctx.JSON(200,scraper.GetProductList("shoes"))
  case "weight":
        ctx.JSON(200,scraper.GetProductList("protein"))
  default :
        ctx.JSON(200,scraper.GetProductList("gyms"))
  }
})
iris.Listen(":8081")
fmt.Println(weather.GetWeather("12.9165167","79.13249859999996"))
}
