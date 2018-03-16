package main

import (
  "github.com/gin-gonic/gin"
   ."bytes"
  _ "github.com/go-sql-driver/mysql"
  "database/sql"
  "crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
  "github.com/garyburd/redigo/redis"
)
type userInfo struct {
  User_id int  `json:"user_id" form:"user_id"`
  User_name    string    `json:"user_name" form:"user_name"`
  User_password  string `json:"user_password" form:"user_password"`
  Reigister_time   int64 `json:"reigister_time" form:"reigister_time"`
  User_phone    string `json:"user_phone" form:"user_phone"`
  User_car_number  string  `json:"user_car_number" form:"user_car_number"`
  User_token  string  `json:"user_token" form:"user_token"`
}

type orderInfo struct {
  Order_id int  `json:"order_id" form:"order_id"`
  Order_publish_name   string  `json:"order_publish_name" form:"order_publish_name"`
  Order_publish_city  string  `json:"order_publish_city" form:"order_publish_city"`
  Order_acceptor_city   string `json:"order_acceptor_city" form:"order_acceptor_city"`
  Order_price    int `json:"order_price" form:"order_price"`
  Order_publish_time  string  `json:"order_publish_time" form:"order_publish_time"`
}


func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
  })
  router.POST("/register", func(c *gin.Context){

        user_name := c.PostForm("user_name")
        user_password:= c.PostForm("user_password")
        user_phone:= c.PostForm("user_phone")
        user_car_number:=c.PostForm("user_car_number")
        register_time := strconv.FormatInt( time.Now().UTC().UnixNano(), 10)[:10]
        if(user_name == ""||user_password == ""|| user_phone == ""||register_time == ""){
          c.JSON(200, gin.H{
            "status":  "400",
            "message": "参数有误",
          })
          return
        }
        c.Writer.Header().Set("Content-type", "application/json")
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept,Content-type,Content-length,Accept-Encoding,X-CSRF-Token,Authorization")
        db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/paiche?charset=utf8")
        defer db.Close()
        // db.Close()
        checkErr(err)
        rows, err := db.Query("SELECT * FROM user_list WHERE user_name = ?",user_name)
        checkErr(err)
        fmt.Println(rows.Next())
        if(!rows.Next()){
          fmt.Println("没有这个人")
          db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/paiche?charset=utf8")
          defer db.Close()
          var buf Buffer
          buf.WriteString("INSERT INTO user_list (user_name,user_password,register_time,user_phone,user_car_number) values (")
          //buf.WriteString(strconv.FormatInt(*organizationId, 10))
          buf.WriteString("'")
          buf.WriteString(user_name)
          buf.WriteString("'")
          buf.WriteString(",")
          buf.WriteString("'")
          buf.WriteString(user_password)
          buf.WriteString("'")
          buf.WriteString(",")
          buf.WriteString("'")
          buf.WriteString(register_time)
          buf.WriteString("'")
          buf.WriteString(",")
          buf.WriteString("'")
          buf.WriteString(user_phone)
          buf.WriteString("'")
          buf.WriteString(",")
          buf.WriteString("'")
          buf.WriteString(user_car_number)
          buf.WriteString("'")
          buf.WriteString(")")

          fmt.Println(buf.String())
          stmt, err := db.Prepare(buf.String())
          if err != nil{
            fmt.Println(err)
          }
          res, err := stmt.Exec()
          checkErr(err)
          id, err := res.LastInsertId()
          fmt.Println(id)
          checkErr(err)
          fmt.Println("id:",id)
          c.JSON(200, gin.H{
            "status":  "200",
            "message": "注册成功",
            "id":    id,
          })
        }else{
          fmt.Println(rows)
          for rows.Next(){
            var user userInfo
            rows.Columns()
            err = rows.Scan(&user.User_id,&user.User_name,&user.User_password,&user.Reigister_time,&user.User_phone,&user.User_car_number)
            checkErr(err)
            fmt.Println(user)
          }
          c.JSON(200, gin.H{
            "status":  "400",
            "message": "用户名已经存在",
           
          })
        }
      
  })
  router.POST("/login", func(c *gin.Context){
    user_name := c.PostForm("user_name")
    user_password:= c.PostForm("user_password")
    c.Writer.Header().Set("Content-type", "application/json")
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept,Content-type,Content-length,Accept-Encoding,X-CSRF-Token,Authorization")
    db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/paiche?charset=utf8")
    defer db.Close()
    // db.Close()
    checkErr(err)
    rows, err := db.Query("SELECT * FROM user_list WHERE user_name = ?",user_name)
    if(!rows.Next()){
      c.JSON(200, gin.H{
        "status":  "400",
        "message": "用户名不存在",
      })
    }else{
      fmt.Println("用户名存在")
      rows, err := db.Query("SELECT * FROM user_list WHERE user_name = ? ",user_name)
      checkErr(err)
      for rows.Next(){
        var user userInfo
        rows.Columns()
        err = rows.Scan(&user.User_id,&user.User_name,&user.User_password,&user.Reigister_time,&user.User_phone,&user.User_car_number)
        checkErr(err)
        fmt.Println(user)
        if(user_password!=user.User_password){
          c.JSON(200, gin.H{
            "status":  "400",
            "message": "用户密码错误",
          })
        }else{
          crutime := time.Now().Unix()
          h := md5.New()
          io.WriteString( h, strconv.FormatInt(crutime, 10))
          token := fmt.Sprintf("%x", h.Sum(nil))
          fmt.Println("token--->", token)
          conn, err := redis.Dial("tcp", "127.0.0.1:6379")
          if err != nil {
              fmt.Println("Connect to redis error", err)
              return
          }
          defer conn.Close()
          _, err = conn.Do("SET", "user_token", token)
          if err != nil {
              fmt.Println("redis set failed:", err)
          }
          user_token, err := redis.String(conn.Do("GET", "user_token"))
          if err != nil {
              fmt.Println("redis get failed:", err)
          } else {
              fmt.Printf("Get user  token: %v \n", user_token)
          }

          c.JSON(200, gin.H{
            "status":  "200",
            "message": "登陆成功",
            "token":  token,
          })
          
        }
      }
    
    }
  })
  router.POST("/orderlist", func(c *gin.Context){
    token := c.PostForm("user_token")
    conn, err := redis.Dial("tcp", "127.0.0.1:6379")

    defer conn.Close()
    user_token, err := redis.String(conn.Do("GET", "user_token"))
    if(token!=user_token){
      c.JSON(200, gin.H{
        "status":  "400",
        "message": "身份验证错误",
      })    
      return
    }
    var order orderInfo
    // fmt.Println(order)
    var orderlist []interface{}
    fmt.Println(token)
    c.Writer.Header().Set("Content-type", "application/json")
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept,Content-type,Content-length,Accept-Encoding,X-CSRF-Token,Authorization")
    fmt.Println("token--->", token)
    db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/paiche?charset=utf8")
    defer db.Close()
    checkErr(err)
    rows, err := db.Query("SELECT * FROM order_list")
    
    rows.Columns()
    defer rows.Close()
    for rows.Next(){
      err = rows.Scan(&order.Order_id,&order.Order_publish_name,&order.Order_publish_city,&order.Order_acceptor_city,&order.Order_price,&order.Order_publish_time)
      checkErr(err)
      // fmt.Println(order)
      orderlist = append(orderlist,order)  
      fmt.Println(order)
    }
    c.JSON(200, gin.H{
      "status":  "200",
      "message": "获取列表成功",
      "data": orderlist,
    })    
  })
  router.POST("/publishorder", func(c *gin.Context){
        order_publish_name := c.PostForm("order_publish_name")
        order_publish_city:= c.PostForm("order_publish_city")
        order_acceptor_city:= c.PostForm("order_acceptor_city")
        order_price:=c.PostForm("order_price")
        
        order_publish_time := strconv.FormatInt( time.Now().UTC().UnixNano(), 10)[:10]

        c.Writer.Header().Set("Content-type", "application/json")
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept,Content-type,Content-length,Accept-Encoding,X-CSRF-Token,Authorization")
        db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/paiche?charset=utf8")
        defer db.Close()
        checkErr(err)
          var buf Buffer
          buf.WriteString("INSERT INTO order_list (order_publish_name,order_publish_city,order_acceptor_city,order_price,order_publish_time) values (")
          //buf.WriteString(strconv.FormatInt(*organizationId, 10))
          buf.WriteString("'")
          buf.WriteString(order_publish_name)
          buf.WriteString("'")
          buf.WriteString(",")
          buf.WriteString("'")
          buf.WriteString(order_publish_city)
          buf.WriteString("'")
          buf.WriteString(",")
          buf.WriteString("'")
          buf.WriteString(order_acceptor_city)
          buf.WriteString("'")
          buf.WriteString(",")
          buf.WriteString("'")
          buf.WriteString(order_price)
          buf.WriteString("'")
          buf.WriteString(",")
          buf.WriteString("'")
          buf.WriteString(order_publish_time)
          buf.WriteString("'")
          buf.WriteString(")")

          fmt.Println(buf.String())
          stmt, err := db.Prepare(buf.String())
          checkErr(err)
          res, err := stmt.Exec()
          checkErr(err)
          id, err := res.LastInsertId()
          fmt.Println(id)
          checkErr(err)
          fmt.Println("id:",id)
          c.JSON(200, gin.H{
            "status":  "200",
            "message": "发布成功",
            "id":    id,
          })
        
        // else{
        //   fmt.Println(rows)
        //   for rows.Next(){
        //     var user userInfo
        //     rows.Columns()
        //     err = rows.Scan(&user.User_id,&user.User_name,&user.User_password,&user.Reigister_time,&user.User_phone,&user.User_car_number)
        //     checkErr(err)
        //     fmt.Println(user)
        //   }
        //   c.JSON(200, gin.H{
        //     "status":  "400",
        //     "message": "用户名已经存在",
           
        //   })
        // }
       

  })
  router.Run(":8080")// listen and serve on 0.0.0.0:8080
}
func checkErr(err error) {
	if err != nil {
    
    panic(err)
    
	}
}

// func Reigister(c *gin.Context){

//   
// }
// func Form_post(c *gin.Context){
 
// }
// func main() {
//  router := gin.Default()

// //router.POST("/register", Reigister)
//  router.POST("/form_post", func(c *gin.Context){
//   message := c.PostForm("message")
//   nick := c.DefaultPostForm("nick", "anonymous")

//   c.JSON(200, gin.H{
//     "status":  "posted",
//     "message": message,
//     "nick":    nick,
//   }
//  } )
// })
//  router.Run(":8080")

// }
