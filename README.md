# 一、项目介绍

采用gin+gorm框架实现了抖音各个功能接口的项目。

项目服务地址：本机wsl子系统，采用端口映射的方式得以让外面的客户端通过宿主机访问到wsl上的服务。

github地址： [github仓库地址](https://github.com/wrxhardworking/douyin-project)


# 二、项目实现

### 2.1 技术选型与相关开发文档

#### GORM

数据库ORM框架，目前最好用最稳定的框架，

[gorm官方文档](https://gorm.io/zh_CN/docs/index.html)

[gorm](http://gorm.io/gorm)

#### GIN

经典的http框架，响应和接收前端请求。

[gin中文文档](https://gin-gonic.com/zh-cn/docs/)

[gin](http://github.com/gin-gonic/gin)

#### REDIS

redis内存数据库，用来做缓存。计数、存储等。

[redis中文网](https://redis.com.cn/documentation.html)

[redigo](http://github.com/gomodule/redigo/redis)

#### MYSQL

常用的db数据库，用来存储关系型数据。

[mysql官方文档](https://dev.mysql.com/doc/)

#### YAML

配置文件，书写配置属性。

[yaml基础语法](https://www.runoob.com/w3cnote/yaml-intro.html)

[yaml.v3](http://gopkg.in/yaml.v3)

#### JWT

Json web token用来用户验证，用户权限等操作。

[jwt官方文档](https://jwt.io/introduction)

[jwt-go](http://github.com/dgrijalva/jwt-go)

#### FFMPEG

用来截取视频的第一帧，生成对应的图片。

[ffmpeg官方文档](https://ffmpeg.org/)

[ffmpeg-go](http://github.com/u2takey/ffmpeg-go)

#### DOCKER

用来承载镜像的mysql和redis

[docker官网](https://www.docker.com/)

### 2.2 架构设计

架构设计我觉得分为两部分，一部分是整体项目的设计，另一部分是数据库的设计。

#### 整体架构：

本项目借鉴mvc三层架构的思想，整体扩散到局部，每个层各司其职，有些承接上下文，利于维护和更新。

项目架构如图所示：

![在这里插入图片描述](https://img-blog.csdnimg.cn/48a283a8e3504cbd9e907432570db4d7.png#pic_center)


##### cache

cache层是关于缓存设计的实现，缓存采用的是redis，主要缓存了一些热数据及经常更新更改的数据。在普通的业务请求上添加一层缓存，能显著提高请求的响应速度。

- ##### common

-  common层就是一些业务常用的相同的结构体和方法。

- ##### config

-  config层就是一些基本的配置选项，mysql、redis、静态文件路径等。这里使用conf.yaml文件来统一管理。

##### dao

dao层就是对db数据库的增删查改操作了，它无关业务代码，只面向底层数据。

##### handle

handle层就是面向前端请求提供的对外接口了，它负责整理好数据，将数据规范化返回给前端。

##### middleware

middleware就是字面意思了，中间件层，用来检验用户token、来判断前端注册账号密码是否符合规范。

##### public

public层就是用来存放静态资源的，把前端传上来的视频以及切割好的图片存于此层，用作路径映射，将资源呈现给前端。

##### service

service层就是用来实现主要的业务逻辑，核心代码所在。将处理完后的数据传给顶层的handle

##### utls

utls层就是工具层，这个层主要是提供一些字符串变换、解析，及外部工具的使用的方法。

#### 数据库架构：

数据库采用的是关系型数据库mysql,数据模型中存在着一对一，一对多，多对多的关系。

数据对象关系如下图：

![在这里插入图片描述](https://img-blog.csdnimg.cn/121e894301fc4226b1687b3566479ee5.jpeg#pic_center)


从上如就可以看出，数据库总共有设计了六张表，其中两张关系表，四张实体表。user->video，video->user在喜欢这个关系中是多对多关系，而user<->user在关注关系中是多对多，usesr->comment、user->message、video->comment则都是一对多关系。

### 2.3 项目代码介绍

#### dao

![在这里插入图片描述](https://img-blog.csdnimg.cn/dbee03f738834cc584d7d92d6eed25c5.png#pic_center)


首先从dao层的代码说起，dao层主要是用来设计对db数据库的访问并把数据封装对象的过程。代码设计全部采用的是单例模式，我的是gorm框架来对数据库进行访问，直接反射到对应的结构体中。在使用gorm的过程中，通过使用外键，使用gorm框架的的has one、has many 、many to many 、belongs to、通过设置标签进行联合查询，有些特殊需求是通过自写sql来解决的。下面通过看几个例子

关联查询：

```Go
// QueryLikeByUserid DeleteLike 查找映射 并且返回lists
func (LikeDao) QueryLikeByUserid(userid int64) ([]Video, error) {
   user := &User{}
   err := db.Preload("VideoLieLists").Where("user_id = ?", userid).Preload("VideoLieLists.Author").Find(user).Error
   if err != nil {
      return nil, err
   }
   return user.VideoLieLists, nil
}
```

手写sql:

```Go
// QueryMessageLists 消息记录
func (MessageDao) QueryMessageLists(userid, ToUserId int64, preMsgTime int64) ([]Message, error) {
   var MessageLists []Message
   MessageLists = make([]Message, 0, 10)
   err := db.Raw("SELECT * FROM `message` WHERE message.create_time>? and ((to_user_id = ? and message.from_user_id  = ?) or (to_user_id = ? and message.from_user_id = ?)) ORDER BY  create_time", preMsgTime, userid, ToUserId, ToUserId, userid).Scan(&MessageLists).Error
   if err != nil {
      return nil, err
   }
   return MessageLists, nil
}
```

总的来说dao层就是对数据的增删改查，**由于时间原因没有设置事务操作，后续会添加。**

#### Service

Service层就是用来实现基本的服务逻辑，为handle层提供服务，通过调用dao层来调用数据实现对应的逻辑。用户登陆的例子来说，拿到从handle层传来的参数，首先判断用户是否存在，**然后对password进行md5加密，**去对比数据库中的密码（数据库中的密码都是md5形式存储，不对db管理者开放），对比完成后，若成功匹配，则登陆成功，返回给handle相应的token信息和用户id。以上就是Service的基本流程，拿到参数，走逻辑，然后返回结果。

```Go
func UserLogin(username string, password string) (*UserRegisInfo, error) {
   var err error
   var token string
   var user *dao.User
   //进行md5加密
   password = utls.Md5Encryption(password)
   user, err = dao.GetUserInstance().QueryUserByName(username)
   //判断用户是否存在
   if user.ID == 0 {
      err = errors.New("user not exists")
      return nil, err
   }
   if err != nil {
      return nil, err
   }
   //验证密码是否正确
   if password != user.Password {
      err = errors.New("password is wrong")
      return nil, err
   }
   //生成token
   token, err = utls.GenerateToken(username, user.ID)
   //成功返回
   return &UserRegisInfo{Token: token, UserID: user.ID}, nil
}
```

#### handle

handle是与前端请求打交道的模块，拿到前端请求的参数，对其进行初步处理，然后调用service层的方法实现对应的功能。拿登陆这个例子来说，首先拿到前端传来的query参数，然后调用service的登陆方法，进而进行错误处理，将返回的数据进行包装返回给前端。可以这么理解就是一个只负责接收数据和传回数据的模块。

```Go
func Login(c *gin.Context) {
   username := c.Query("username")
   password := c.Query("password")

   info, err := service.UserLogin(username, password)
   if err != nil {
      log.Println(err.Error())
      c.JSON(http.StatusOK, UserRegisterResponse{
         Response: common.Response{StatusCode: 1, StatusMsg: err.Error()},
      })
      return
   } else {
      c.JSON(http.StatusOK, UserRegisterResponse{
         Response: common.Response{StatusCode: 0},
         Token:    info.Token,
         UserId:   info.UserID,
      })
      return
   }
}
```

#### middleware

中间件层，前端进入后端的一道门槛。首先我设计了注册时**检验账号和密码的规范性**的中间件，利用正则表达式，检测账号密码的合理性。其次为了保证用户的使用流畅度 ，我找了一个中间件限制了post请求大小，也就是上传文件的请求大小。

```Go
package middleware

import (
   "douyin/common"
   "github.com/gin-gonic/gin"
   "net/http"
   "regexp"
)

func Check() gin.HandlerFunc {

   return func(c *gin.Context) {
      username := c.Query("username")
      password := c.Query("username")
      //传给handle层
      //5～16字节，允许字母、数字、下划线，以字母开头
      matchString := "^[a-zA-Z][a-zA-Z0-9_]{4,15}$"
      usernameMatch, _ := regexp.MatchString(matchString, username)
      passwordMatch, _ := regexp.MatchString(matchString, password)

      if usernameMatch && passwordMatch != true {
         c.JSON(http.StatusOK, common.Response{
            StatusCode: 1, StatusMsg: "Account or password is illegal",
         })
         c.Abort()
         return
      }
      c.Set("username", username)
      c.Set("password", password)
      //挂起
      c.Next()
   }
}
```

当前端注册请求过来时，会先经过中间件结构，中间件中会判断账号密码的合理性，若不合理，直接调用gin框架的abort（）进行终止，终止掉后面所有的该请求下的函数。

其次我设计了**用户权限**，对token的检查，检查其是否有效，是否正确。

```Go
package middleware

import (
   "douyin/cache"
   "douyin/common"
   "douyin/dao"
   "douyin/handle"
   "douyin/utls"
   "fmt"
   "github.com/gin-gonic/gin"
   "net/http"
)

func UserAuth() gin.HandlerFunc {

   return func(c *gin.Context) {
      //得到token字段
      //1.get请求
      token := c.Query("token")
      if token == "" {
         //2.post请求
         token = c.PostForm("token")
      }

      // 两种情况下来，判断是否有token
      if token == "" {
         c.JSON(http.StatusOK, handle.UserResponse{
            Response: common.Response{StatusCode: 1, StatusMsg: "token 不存在"},
         })
         c.Abort()
      }
      //解析
      t, claim, err := utls.ParseToken(token)
      //判断是否有效
      if !t.Valid || err != nil {
         c.JSON(http.StatusOK, handle.UserResponse{
            Response: common.Response{StatusCode: 1, StatusMsg: "token有效期过了或者" + err.Error()},
         })
         c.Abort()
         return
      }
      //1.首先到redis中查找，没有的话去mysql中查找
      //2.mysql中没有说明token失败
      var isExists = true

      err = cache.UserIsExists(claim.UserId)
      if err != nil {
         fmt.Println(err)
         //在redis中不存在
         isExists = false
      }
      if !isExists {
         //进行db查找
         var user *dao.User
         user, err = dao.GetUserInstance().QueryUserByID(claim.UserId)
         if err != nil {
            c.JSON(http.StatusOK, handle.UserResponse{
               Response: common.Response{StatusCode: 1, StatusMsg: "token find failed"},
            })
            c.Abort()
            return
         }
         if user.ID == 0 {
            c.JSON(http.StatusOK, handle.UserResponse{
               Response: common.Response{StatusCode: 1, StatusMsg: "id is not exists"},
            })
            c.Abort()
            return
         }
      }
      //传给handle层
      c.Set("userid", claim.UserId)
      //挂起
      c.Next()
   }
}
```

其中token是jwt生成的。在utls包中封装了token的生成与解析的方法，在生成token的时候在，中在token字段中加入了userid，通过检测userid是否存在token是否存在。这其中还用到了redis缓存，这在cache中再继续介绍。

#### cache

1.在实现业务逻辑的时候，发现对于登陆用户，请求视频流接口的时候，对于该用户对视频的状态和对视频作者的状态确定起来十分棘手，因为想追求性能，又要做到不业务不能出错，首先想到的是查询出所有like表和follow中的映射关系，进行循环查找，后面一想，当用户多了起来，时间复杂度直接on2，太影响性能了。最终决定用redis进行缓存。利用redis中的set集合，用行为加userid作为主键，将关注的userid或喜欢的videoid分别放入对应的集合，在进行状态判断的时候首先从redis中查找该videoid或者userid是否在“喜欢集合”或者“关注集合中”，若不存在就进行db查找。（其实还是有待改进）下面是实现缓存的一些接口，对外提供了添加、删除、和判断的功能，以供使用。

```Go
func getUserRelationKey(userid int64) string {
   return fmt.Sprintf("%s_%d", "UserRelationKey", userid)
}

func getUserVideoRelation(userid int64) string {
   return fmt.Sprintf("%s_%d", "UserVideoRelationKey", userid)
}

// SetUserRelation 建立用户和用户的关系集合
func SetUserRelation(userid, touserId int64) error {
   conn := RedisPool.Get()
   defer func(conn redis.Conn) {
      err := conn.Close()
      if err != nil {
      }
   }(conn)

   key := getUserRelationKey(userid)
   //往集合中加关注的人
   _, err := conn.Do("SADD", key, touserId)
   if err != nil {
      return err
   }
   return nil
}

// SetUserVideoRelation 建立用户和视频的关系集合
func SetUserVideoRelation(userid, videoId int64) error {
   conn := RedisPool.Get()
   defer func(conn redis.Conn) {
      err := conn.Close()
      if err != nil {
      }
   }(conn)
   key := getUserVideoRelation(userid)
   //往集合中加关注的人
   _, err := conn.Do("SADD", key, videoId)
   if err != nil {
      return err
   }
   return nil
}

// IsUserRelation 判断是否在集合中
func IsUserRelation(userid, touserId int64) bool {
   conn := RedisPool.Get()
   defer func(conn redis.Conn) {
      err := conn.Close()
      if err != nil {
      }
   }(conn)
   key := getUserRelationKey(userid)
   res, err := conn.Do("SISMEMBER", key, touserId)
   if err != nil {
      log.Println(err.Error())
      return false
   }
   if res == 0 {
      return false
   }
   return true
}

// IsUserVideoRelation 判断是否在集合里面
func IsUserVideoRelation(userid, videoId int64) bool {
   conn := RedisPool.Get()
   defer func(conn redis.Conn) {
      err := conn.Close()
      if err != nil {
      }
   }(conn)

   key := getUserVideoRelation(userid)

   res, err := conn.Do("SISMEMBER", key, videoId)
   if err != nil {
      log.Println(err.Error())
      return false
   }
   if res == 0 {
      return false
   }
   return true
}
func DeleteUserRelation(userid, touserId int64) error {
   conn := RedisPool.Get()
   defer func(conn redis.Conn) {
      err := conn.Close()
      if err != nil {
      }
   }(conn)

   key := getUserRelationKey(userid)
   //往集合中加关注的人
   _, err := conn.Do("SREM", key, touserId)
   if err != nil {
      return err
   }
   return nil
}

// DeleteUserVideoRelation 删除关系
func DeleteUserVideoRelation(userid, videoId int64) error {
   conn := RedisPool.Get()
   defer func(conn redis.Conn) {
      err := conn.Close()
      if err != nil {
      }
   }(conn)
   key := getUserVideoRelation(userid)
   //往集合中加关注的人
   _, err := conn.Do("SREM", key, videoId)
   if err != nil {
      return err
   }
   return nil
}
```

2.继续写到后面看到这三个接口：

  optional int64 total_favorited = 9;   //获赞数量

  optional int64 work_count = 10;     //作品数量

  optional int64 favorite_count = 11; //点赞数量

在数据库中若加入了这三个字段，就需要频繁的进行更新；若不加入这三个字段利用sql语句提供的count进行计数，你会发现，count操作比较耗费性能，当数量一多那更是“动弹不得”。所以对于计数操作，我也进行了redis缓存，看了字节内部的使用redis那门课学了点。就对外提供了对上述三个字段的incr、decr、和get的接口。

```Go
// SetUserCount 设置user计数
func SetUserCount(userid int64) error {
   conn := RedisPool.Get()

   defer func(conn redis.Conn) {
      err := conn.Close()
      if err != nil {
      }
   }(conn)

   key := getUserCountKey(userid)
   _, err := conn.Do("hmset", redis.Args{key}.AddFlat(map[string]int64{
      "followCount":   0,
      "followerCount": 0,
      "workCount":     0,
      "favoriteCount": 0,
      "totalFavorite": 0,
   })...)
   if err != nil {
      return err
   }
   return nil
}

func GetUserFollowCount(userID int64) (int64, error) {
   res, err := get(userID, "followCount")
   return res, err
}
func GetUserFollowerCount(userID int64) (int64, error) {
   res, err := get(userID, "followerCount")
   return res, err
}
func GetUserWorkCount(userID int64) (int64, error) {
   res, err := get(userID, "workCount")
   return res, err
}
func GetUserFavoriteCount(userID int64) (int64, error) {
   res, err := get(userID, "favoriteCount")
   return res, err
}
func GetUserTotalFavoriteCount(userID int64) (int64, error) {
   res, err := get(userID, "totalFavorite")
   return res, err
}

func DecrByUserFollowCount(userID int64) error {
   err := change(userID, "followCount", -1)
   return err
}
func IncrByUserFollowCount(userID int64) error {
   err := change(userID, "followCount", 1)
   return err
}

func DecrByUserFollowerCount(userID int64) error {
   err := change(userID, "followerCount", -1)
   return err
}
func IncrByUserFollowerCount(userID int64) error {
   err := change(userID, "followerCount", 1)
   return err
}
func DecrByUserWorkCount(userID int64) error {
   err := change(userID, "workCount", -1)
   return err
}
func IncrByUserWorkCount(userID int64) error {
   err := change(userID, "workCount", 1)
   return err
}

func DecrByUserTotalFavorite(userID int64) error {
   err := change(userID, "totalFavorite", -1)
   return err
}
func IncrByUserTotalFavorite(userID int64) error {
   err := change(userID, "totalFavorite", 1)
   return err
}

func DecrByUserFavoriteCount(userID int64) error {
   err := change(userID, "favoriteCount", -1)
   return err
}
func IncrByUserFavoriteCount(userID int64) error {
   err := change(userID, "favoriteCount", 1)
   return err
}

func change(userid int64, field string, incr int64) error {

   key := getUserCountKey(userid)
   conn := RedisPool.Get()

   //释放资源
   defer func(conn redis.Conn) {
      err := conn.Close()
      if err != nil {
      }
   }(conn)

   isExists, _ := conn.Do("exists", key)
   //判断键值是否存在
   if isExists == 0 {
      return errors.New("is not exists")
   }

   before, err := redis.Int64(conn.Do("hget", key, field))
   if err != nil {
      return err
   }
   if before+incr < 0 {
      //此时已经小于0了
      fmt.Println(errors.New("already <0"))
      log.Println(errors.New("already <0"))
      return nil
   }
   _, err = conn.Do("HIncrBy", key, field, incr)
   if err != nil {
      return err
   }
   return nil
}

func get(userid int64, field string) (int64, error) {
   key := getUserCountKey(userid)

   conn := RedisPool.Get()
   //释放资源
   defer func(conn redis.Conn) {
      err := conn.Close()
      if err != nil {
      }
   }(conn)

   isExists, _ := conn.Do("exists", key)

   //判断键值是否存在
   if isExists == 0 {
      return 0, errors.New("is not exists")
   }

   res, err := redis.Int64(conn.Do("hget", key, field))

   if err != nil {
      return 0, err
   }
   return res, nil
}
```

3.缓存中我是用的是redis的连接池的技术。

```Go
func RedisPoolInit() error {
   RedisPool = &redis.Pool{
      MaxIdle:     config.C.Redis.Maxidle,   //最大空闲数
      MaxActive:   config.C.Redis.Maxactive, //最大连接数，0不设上
      Wait:        true,
      IdleTimeout: time.Duration(1) * time.Second, //空闲等待时间
      Dial: func() (redis.Conn, error) {
         c, err := redis.Dial("tcp", config.C.Redis.Ipaddress+":"+config.C.Redis.Port) //redis IP地址
         if err != nil {
            fmt.Println(err)
            return nil, err
         }
         //密码认证
         if _, err = c.Do("AUTH", config.C.Redis.Authpassword); err != nil {
            err := c.Close()
            if err != nil {
               return nil, err
            }
            return nil, err
         }
         redis.DialDatabase(0)
         return c, err
      },
   }
   return nil
}
```

#### config

对于config层，首先手写配置文件项.yaml，然后利用代码生成工具生成了对应的结构体。然后对其读取和解析。

```YAML
# mysql
mysql:
  username: root
  password: 123456
  ipaddress: 127.0.0.1
  port: 3309
  dbname: douyin

# redis
redis:
  ipaddress: 127.0.0.1
  port: 6379  
  authpassword: 123456    #redis验证密码
  maxidle: 5              #连接池最大空闲数
  maxactive: 0            #连接池最大连接数 0 表示不设上

# resource(静态图片和视频的资源)
resouece:
  ipaddress: 172.25.169.130 #宿主机的ip地址
  port: 12345               #宿主机的端口
  
type Config struct {
   Mysql    Mysql    `yaml:"mysql"`
   Redis    Redis    `yaml:"redis"`
   Resouece Resouece `yaml:"resouece"`
}

type Mysql struct {
   Username  string `yaml:"username"`
   Password  string `yaml:"password"`
   Ipaddress string `yaml:"ipaddress"`
   Port      string `yaml:"port"`
   Dbname    string `yaml:"dbname"`
}

type Redis struct {
   Ipaddress    string `yaml:"ipaddress"`
   Port         string `yaml:"port"`
   Authpassword string `yaml:"authpassword"`
   Maxidle      int    `yaml:"maxidle"`
   Maxactive    int    `yaml:"maxactive"`
}

type Resouece struct {
   Ipaddress string `yaml:"ipaddress"`
   Port      string `yaml:"port"`
}

var C Config

func ConfInit() error {
   yamlFile, err := os.ReadFile("./config/config.yaml")
   if err != nil {
      fmt.Println(err.Error())
      return err
   }
   // 将读取的yaml文件解析为响应的 struct
   err = yaml.Unmarshal(yamlFile, &C)
   if err != nil {
      fmt.Println(err.Error())
      return err
   }
   return nil
}
```

#### utils

对于工具包就是一些工具的使用了，有jwt、MD5等。重点go中是如何使用ffmpeg的，传video路径，传入图片生成路径，以及你想截取的帧序。

```Go
func GenerateSnapshot(videoPath, snapshotPath string, frameNum int) (err error) {
   buf := bytes.NewBuffer(nil)
   err = ffmpeg.Input(videoPath).
      Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
      Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
      WithOutput(buf, os.Stdout).
      Run()
   if err != nil {
      return err
   }

   img, err := imaging.Decode(buf)
   if err != nil {
      return err
   }

   err = imaging.Save(img, snapshotPath+".png")
   if err != nil {
      return err
   }
   return nil
}
```

# 三、测试结果

### 功能测试：

#### user_info：

##### TestUserRegister

```Go
tests := []struct {
   name string
   args
}{
   {
      "测试1",
      args{
         "____________",
         "123456mksjxnjancanjskandndjasnjdkasn",
      },
   },
   {
      "测试2",
      args{
         "小王",
         "123456",
      },
   },
   {
      "测试3",
      args{
         "",
         "123456",
      },
   },
}
```

##### TestUserLogin

```Go
tests := []struct {
   name string
   args
}{
   {
      "测试1",
      args{
         "小王",
         "123456",
      },
   },
   {
      "测试2",
      args{
         "_____",
         "123456",
      },
   },
   {
      "测试3",
      args{
         "hhhy",
         "123456",
      },
   },
```

#### video_info:

##### TestPublishedVideoLists

```Go
type args struct {
   userid int64
}
tests := []struct {
   name string
   args
}{
   {
      "测试1",
      args{
         1,
      },
   },
   {
      "测试2",
      args{
         100000,
      },
   },
   {
      "测试3",
      args{
         -10,
      },
   },
}
```

#### message_info：

##### TestGetMessageLists

```Go
type args struct {
   userid   int64
   toUserId int64
   preTime  int64
}
tests := []struct {
   name string
   args
}{
   {
      name: "测试1",
      args: args{
         userid: 10,
         toUserId: 11,
         preTime: 213123213,
      },
   },
   {
      name: "测试2",
      args: args{
         userid: 10,
         toUserId: -11,
         preTime: 00213123213,
      },
   },
}
```

##### TestSendMessage

```Go
type args struct {
   userid   int64
   toUserId int64
   content  string
   action   string
}
tests := []struct {
   name string
   args
}{
   {
      name: "测试1",
      args: args{
         userid:   10,
         toUserId: 11,
         content:  "你好啊",
         action:   "1",
      },
   },
   {
      name: "测试2",
      args: args{
         userid:   11,
         toUserId: 10,
         content:  "你也好啊",
         action:   "1",
      },
   },
}
```

#### like_info：

##### TestGetLikeLists

```Go
tests := []struct {
   name   string
   userid int64
}{
   {
      name:   "test1",
      userid: 10,
   },
   {
      name:   "test2",
      userid: 11,
   },
}
```

##### TestThumbUpOrCancel

```Go
tests := []struct {
   name    string
   userid  int64
   videoId int64
   action  string
}{
   {
      name:    "test1",
      userid:  10,
      videoId: 1,
      action:  "1",
   },
   {

      name:    "test2",
      userid:  11,
      videoId: 2,
      action:  "1",
   },
   {
      name:    "test3",
      userid:  10,
      videoId: 2,
      action:  "1",
   },
}
```

#### follow_info：

##### TestFollowerLists

```Go
tests := []struct {
   name   string
   userid int64
}{
   {
      name:   "test1",
      userid: 10,
   },
   {
      name:   "test2",
      userid: 11,
   },
}
```

##### TestFollowLists

```Go
tests := []struct {
   name   string
   userid int64
}{
   {
      name:   "test1",
      userid: 10,
   },
   {
      name:   "test2",
      userid: 16,
   },
}
```

#### comment_info：

##### TestCommentOrDelete

```Go
tests := []struct {
   name        string
   videoId     int64
   userId      int64
   action      string
   commentText string
   commentId   int64
}{
   {
      name:        "test1",
      videoId:     1,
      userId:      10,
      action:      "1",
      commentText: "test_comment1",
      commentId:   13,
   },
   {
      name:        "test2",
      videoId:     2,
      userId:      10,
      action:      "1",
      commentText: "test_comment2",
      commentId:   11,
   },
   {
      name:        "test3",
      videoId:     3,
      userId:      12,
      action:      "1",
      commentText: "test_comment3",
      commentId:   12,
   },
}
```

##### TestGetCommentLists

```Go
tests := []struct {
   name    string
   videoId int64
}{
   {
      name:    "test1",
      videoId: 1,
   },
   {
      name:    "test2",
      videoId: 2,
   },
   {
      name:    "test3",
      videoId: 3,
   },
}
```



# 四、项目总结与反思

### 1.仍然存在的问题：

- 数据库没有进行事务操作，可能会导致数据的不安全，在极高并发的时候还需要考虑读写锁。
- redis没有设计好合适的过期值和分布式锁。
- 测试不够完善，后续进行完善。

### 2.已经识别出来的优化项目

- 所有的有关计数的count的操作都可以用redis缓存来解决，项目中关注数和粉丝数还未实现。
- 可以考虑用go语言的泛型来处理一些冗余的代码，这等待后续学习go的泛型。
- 项目中还存在一部分拷贝，考虑优化成指针等。

### 3.架构演进的可能性

- 可以尝试进行分布式开发。
- 可以引进短信验证的功能，增强安全性。

# 五、其他补充资料

通过端口映射的方法使得外部的客户端能访问到wsl：

https://juejin.cn/post/7198169454186070071
