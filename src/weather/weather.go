package weather
import(
  "github.com/levigross/grequests"
  "github.com/Jeffail/gabs"
)
type Weather struct{
  Main string `json:main`
}

func GetWeather(lat,lng string)float64{
 res,err:=grequests.Get("http://api.openweathermap.org/data/2.5/weather?lat="+lat+"&lon="+lng+"&APPID=f06e779ae27b1a8269f171f8372fc352",nil)
 if err!=nil{
 }
 var result float64;
 jsonParsed, _ := gabs.ParseJSON(res.Bytes())
 result=jsonParsed.Path("main.humidity").Data().(float64);

  return result
}
