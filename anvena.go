
package main
import (
        "errors"
        "fmt"
        "net/http"
        "github.com/PuerkitoBio/goquery"
        "encoding/json"
        "github.com/julienschmidt/httprouter"
        "strings"
        "github.com/Jeffail/gabs"
        "io/ioutil"
)
type Weather struct{
  Main string `json:main`
}
/*
@params lat latitude
@params lng longitude
returns humidty,current weather location
*/
func GetWeather(lat,lng string)(float64,string){
 res,err:=http.Get("http://api.openweathermap.org/data/2.5/weather?lat="+lat+"&lon="+lng+"&APPID=f06e779ae27b1a8269f171f8372fc352")
 if err!=nil{
 }
 var result float64;
 defer res.Body.Close()
 body, _ := ioutil.ReadAll(res.Body)
 jsonParsed, _ := gabs.ParseJSON(body)
 result=jsonParsed.Path("main.humidity").Data().(float64);
 main:=jsonParsed.Path("weather.main").Data().([]interface{})
  return result,main[0].(string)
}

type hotel struct {
  Name string `json:"name"`
  Url string `json:url`
}
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

func HotelByFood(food string,place string)([]string,[]string,error){
  doc,err := goquery.NewDocument("https://www.zomato.com/"+place+"/restaurestaurants?q="+food)
   if err != nil {
    return nil,nil,errors.New("something went wrong")
   }
   result:=make([]string,0)
   url:=make([]string,0)
   resultch:=make(chan []string, 0)
  urlch:=make(chan []string, 0)
  length:=doc.Find(".result-title").Length()
  fmt.Print(length)
  go doc.Find(".result-title").Each(func (i int , s *goquery.Selection)  {
       fmt.Print(s.Text())
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
 func demo(food string)([]hotel,error)  {
   doc,err := goquery.NewDocument("https://www.foodpanda.in/restaurants?user_search="+food+"&sort=&sort=&page=1")
    if err != nil {
     //return nil,nil,errors.New("something went wrong")
     return nil, errors.New("something went wrong")
    }
    var result []hotel
    doc.Find("article").Each(func (i int, s *goquery.Selection)  {
      name:=strings.Trim(s.Find("span.vendor__name").Text(),"\n")
      name=strings.Trim(name," ")
      name=strings.Trim(name,"\n")
      url,_:=s.Find("img").Attr("data-src")
      url="https:"+url
      result=append(result,hotel{name,url})
    })
    return result,nil
 }
 func foodAdHandler(w http.ResponseWriter,r *http.Request, ps httprouter.Params)  {
   hotels,_:=demo(ps.ByName("food"))
   w.Header().Set("Content-Type", "application/json")
   jData,_:=json.Marshal(hotels)
   w.Write(jData)
 }

 func weatherAdHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
   queryValues:=r.URL.Query()
   lat:=queryValues.Get("lat")
   lng:=queryValues.Get("lng")
   humidty,weather:=GetWeather(lat,lng)
   fmt.Print(weather)
   fmt.Print(humidty)
   if weather=="Rain"{
      w.Header().Set("Content-Type", "application/json")
      jData,_:=json.Marshal(GetProductList("rain"))
      w.Write(jData)
     }else  if humidty>=70{
        w.Header().Set("Content-Type", "application/json")
        jData,_:=json.Marshal(GetProductList("cotton"))
        w.Write(jData)
     }else if humidty<70{
        w.Header().Set("Content-Type", "application/json")
        jData,_:=json.Marshal(GetProductList("trekking"))
        w.Write(jData)
   }
 }
 func workoutAdHandler(w http.ResponseWriter,r *http.Request,ps httprouter.Params)  {
   queryValues:=r.URL.Query()
   switch queryValues.Get("q") {
      case "running":
        w.Header().Set("Content-Type", "application/json")
        jData,_:=json.Marshal(GetProductList("running"))
        w.Write(jData)
      case "weight":
        w.Header().Set("Content-Type", "application/json")
        jData,_:=json.Marshal(GetProductList("protein"))
        w.Write(jData)
      default :
      w.Header().Set("Content-Type", "application/json")
      jData,_:=json.Marshal(GetProductList("gyms"))
      w.Write(jData)
   }
 }
func main()  {
  router:=httprouter.New()
  router.GET("/foodad/:food",foodAdHandler)
  router.GET("/weatherad",weatherAdHandler)
  router.GET("/workout",workoutAdHandler)
  http.ListenAndServe(":80",router)
}
