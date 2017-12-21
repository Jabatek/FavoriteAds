package main

import (
	"net/http"
	"io/ioutil"
	"strings"
)



func request(s string) string{
	resp, err := http.Get("http://localhost:8080/graphql?query="+s)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body)
}

func main() {


	println()
	println("Select all favorite of user 1: must get empty favorite")
	println(request("query{User(id:1){id,favorites{id,uid,url,savetype}}}"))


	println()
	println("insert url1 : must get one favorite with url=url1 and  savetype=auto")
	request("mutation{addfads(uid:1,url:\"url1\",savetype:\"auto\"){url}}")
	println(request("query{User(id:1){id,favorites{id,uid,url,savetype}}}"))


	println()
	println("insert new url1 to  generate conflict replace: we must still get only one favorite")
	var sel =request("mutation{addfads(uid:1,url:\"url1\",savetype:\"auto\"){id}}")
	var idUrl1=strings.Replace(strings.Replace(sel, "{\"data\":{\"addfads\":{\"id\":", "", -1),"}}}","",-1)
	println(request("query{User(id:1){id,favorites{id,uid,url,savetype}}}"))



	println()
	println("insert url2: we must get two favorites url1 and url2 with url2.savetype=manual")
	sel = request("mutation{addfads(uid:1,url:\"url2\",savetype:\"manual\"){id}}")
	var idUrl2=strings.Replace(strings.Replace(sel, "{\"data\":{\"addfads\":{\"id\":", "", -1),"}}}","",-1)
	println(request("query{User(id:1){id,favorites{uid,url,savetype}}}"))


	println()
	println("delet url2 and url1:we must get empty favorite ")
		request("mutation{deletfads(id:"+idUrl2+"){id,}}")
		request("mutation{deletfads(id:"+idUrl1+"){id,}}")
	println(request("query{User(id:1){id,favorites{uid,url,savetype}}}"))


}