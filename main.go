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
  Order_state int  `json:"order_state" form:"order_state"`
  Order_contact_phone string `json:"order_contact_phone" form:"order_contact_phone"`
}


func main() {
	router := gin.Default()
  router.POST("/cancelorder", func(c *gin.Context){
     user_id := c.PostForm("user_id")
   
    order_id := c.PostForm("order_id")

    if(len(order_id)==0){
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
    checkErr(err)
    // rows, err := db.Query("SELECT * FROM order_list WHERE order_state = ?",user_id)
    // fmt.Println(rows.Next())
    stmt, err := db.Prepare("update order_list set order_state=? where order_id=? and order_state = ?")
   
    checkErr(err)
    res, err := stmt.Exec("0",order_id,user_id)
    checkErr(err)
    affect, err := res.RowsAffected()
    checkErr(err)
    
    if affect==0{
      c.JSON(200, gin.H{
        "status":  "200",
        "message": "订单不存在",
        "data":  affect,
      })   
    }else if affect== 1 {
      c.JSON(200, gin.H{
        "status":  "200",
        "message": "取消订单成功",
        "data":  affect,
      })   
    }else{
      c.JSON(200, gin.H{
        "status":  "400",
        "message": "其他错误",
        "data":  affect,
      })   
    }
     

  })
  router.POST("/register", func(c *gin.Context){

        user_name := c.PostForm("user_name")
        user_password:= c.PostForm("user_password")
        user_phone:= c.PostForm("user_phone")
        user_car_number:=c.PostForm("user_car_number")
        register_time := strconv.FormatInt( time.Now().UTC().UnixNano(), 10)[:10]
       

        if(len(user_name)==0||len(user_name)==0||len(user_name)==0||len(user_car_number)==0){
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
        checkErr(err)
        res, err := db.Exec("INSERT INTO user_list(user_name,user_password,register_time,user_phone,user_car_number) SELECT '"+user_name+"','"+user_password+"','"+register_time+"','"+user_phone+"','"+user_car_number+"' FROM DUAL WHERE NOT EXISTS(SELECT user_name FROM user_list WHERE user_name = '"+user_name+"')")
        checkErr(err)
        id, err := res.LastInsertId()  
        checkErr(err)
        fmt.Println(id);
        if(id==0){
               c.JSON(200, gin.H{
                  "status":  "400",
                  "message": "用户名已经存在",
                })
        }else{
               c.JSON(200, gin.H{
                  "status":  "200",
                  "message": "注册成功",
                  "id":    id,
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
    c.Writer.Header().Set("Content-type", "application/json;charset=utf-8")
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept,Content-type,Content-length,Accept-Encoding,X-CSRF-Token,Authorization")
    pageNo,err :=strconv.Atoi(c.PostForm("page"))
    checkErr(err)
    pageSize,err  :=strconv.Atoi(c.PostForm("size"))
    checkErr(err)
    token := c.PostForm("user_token")
    fmt.Println(token)
    fmt.Println("token---->用户传的token",token)

    conn, err := redis.Dial("tcp", "127.0.0.1:6379")

    defer conn.Close()
    user_token, err := redis.String(conn.Do("GET", "user_token"))
    fmt.Println("user_token---->redis",token)
    if(token!=user_token){
      c.JSON(200, gin.H{
        "status":  "400",
        "message": "身份验证错误",
      })    
      return
    }
    var order orderInfo
    var orderlist []interface{}
    fmt.Println(token)
   
    fmt.Println("token--->", token)
    db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/paiche?charset=utf8")
    defer db.Close()
    checkErr(err)
    rows, err := db.Query("SELECT * FROM order_list   WHERE order_state = 0 LIMIT "+strconv.Itoa((pageNo-1)*pageSize)+","+strconv.Itoa(pageSize))
    rows.Columns()
    defer rows.Close()
    for rows.Next(){
      err = rows.Scan(&order.Order_id,&order.Order_publish_name,&order.Order_publish_city,&order.Order_acceptor_city,&order.Order_price,&order.Order_publish_time,&order.Order_state,&order.Order_contact_phone)
      checkErr(err)
      // fmt.Println(order)
      orderlist = append(orderlist,order)  
      fmt.Println(order)
    }
    if len(orderlist)< 10{
      c.JSON(200, gin.H{
       
        "status":  "200",
        "message": "获取列表成功",
        "isNext":    false,
        "data": orderlist,
        
        
      })    
    }else{
      c.JSON(200, gin.H{
       
        "status":  "200",
        "message": "获取列表成功",
        "isNext":    true,
        "data": orderlist,
      
      })    
    }
   
  })
  router.POST("/publishorder", func(c *gin.Context){
        order_publish_name := c.PostForm("order_publish_name")
        order_publish_city:= c.PostForm("order_publish_city")
        order_acceptor_city:= c.PostForm("order_acceptor_city")
        order_price:=c.PostForm("order_price")
        order_state:=c.PostForm("order_state")
        order_contact_phone:=c.PostForm("order_contact_phone")
        order_publish_time := strconv.FormatInt( time.Now().UTC().UnixNano(), 10)[:10]

        c.Writer.Header().Set("Content-type", "application/json")
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept,Content-type,Content-length,Accept-Encoding,X-CSRF-Token,Authorization")
        db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/paiche?charset=utf8")
        defer db.Close()
        checkErr(err)
          var buf Buffer
          buf.WriteString("INSERT INTO order_list (order_publish_name,order_publish_city,order_acceptor_city,order_price,order_publish_time,order_state,order_contact_phone) values (")
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
          buf.WriteString(",")
          buf.WriteString("'")
          buf.WriteString(order_state)
          buf.WriteString("'")
          buf.WriteString(",")
          buf.WriteString("'")
          buf.WriteString(order_contact_phone)
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
        
  })
  router.POST("/getorder",func(c *gin.Context){
    user_id := c.PostForm("user_id")
    order_id := c.PostForm("order_id")
 
    var order orderInfo

    c.Writer.Header().Set("Content-type", "application/json")
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept,Content-type,Content-length,Accept-Encoding,X-CSRF-Token,Authorization")
    db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/paiche?charset=utf8")
    defer db.Close()
    checkErr(err)
    rows, err := db.Query("SELECT * FROM order_list WHERE order_id = ?",order_id)
    for rows.Next(){
      err = rows.Scan(&order.Order_id,&order.Order_publish_name,&order.Order_publish_city,&order.Order_acceptor_city,&order.Order_price,&order.Order_publish_time,&order.Order_state)
      checkErr(err)
      // fmt.Println(order)
      fmt.Println(order)
    }
    if(order.Order_state!=0){
          c.JSON(200, gin.H{
            "status":  "400",
            "message": "订单已不存在",
          })    
          return
    }else{
      stmt, err := db.Prepare("update order_list set order_state=? where order_id=?")
        checkErr(err)
        res, err := stmt.Exec(user_id,order_id)
        checkErr(err)
        affect, err := res.RowsAffected()
        checkErr(err)
        fmt.Println(affect)
        c.JSON(200, gin.H{
          "status":  "200",
          "message": "领取订单成功",
        })    
    }
   })
  router.Run(":8080")// listen and serve on 0.0.0.0:8080
}
func checkErr(err error) {
	if err != nil {
    
    panic(err)
    
	}
}

