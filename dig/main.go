package main

import (
	"fmt"
	"log"

	"github.com/jessevdk/go-flags"
	"go.uber.org/dig"
	"gopkg.in/ini.v1"
)

// -c 或者 --config 来指定配置文件
type Option struct {
	ConfigFile string `short:"c" long:"config" description:"Name of config file."`
}

func InitOption() (*Option, error) {
	var opt Option
	_, err := flags.Parse(&opt) // 解析op
	return &opt, err
}

func GetOption() *Option {
	return &Option{ConfigFile: "./demo.ini"}
}
func InitConf(opt *Option) (*ini.File, error) {
	cfg, err := ini.Load(opt.ConfigFile) // 拿到参数, 传入ini包
	return cfg, err
}

// 依赖于ini.cfg
func PrintInfo(cfg *ini.File) {
	fmt.Println("App Name:", cfg.Section("").Key("app_name").String())
	fmt.Println("Log Level:", cfg.Section("").Key("log_level").String())
}
func test1() {
	// 这里的依赖是
	// printinfo --> conf
	// conf --> initconf --> option
	// option --> initoption
	container := dig.New() // 创建注入容器

	// 对于重复的key, 会对其进行覆盖
	container.Provide(GetOption) // 注册option的提供 --> 根据initoption的返回值, 插入一个映射: 需要Option类时,就调用该函数
	container.Provide(InitOption)
	container.Provide(InitConf) // 注册conf的提供 --> 同理, 需要Conf类时就调用InitConf函数
	container.Invoke(PrintInfo) // 调用函数 --> 检查参数, 遍历函数栈树, 对于所有需要的参数尝试用映射中的方法去构造
}

// ----------------------------------------------------------------------------

type RedisConfig struct {
	IP   string
	Port int
}
type MySQLConfig struct {
	IP       string
	Port     int
	User     string
	Password string
	Database string
}
type Config struct {
	// dig.In // 用于参数注入, 即使dig中并没有用于构造Config的函数, 其依然可以自行构造, 通过查找成员的构造来实现
	dig.Out
	Redis *RedisConfig
	MySQL *MySQLConfig
}

func InitConfig(opt *Option) (*ini.File, error) {
	cfg, err := ini.Load(opt.ConfigFile)
	return cfg, err
}

func InitRedisConfig(cfg *ini.File) (*RedisConfig, error) {
	port, err := cfg.Section("redis").Key("port").Int()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &RedisConfig{
		IP:   cfg.Section("redis").Key("ip").String(),
		Port: port,
	}, nil
}
func InitMySQLConfig(cfg *ini.File) (*MySQLConfig, error) {
	port, err := cfg.Section("mysql").Key("port").Int()
	if err != nil {
		return nil, err
	}
	return &MySQLConfig{
		IP:       cfg.Section("mysql").Key("ip").String(),
		Port:     port,
		User:     cfg.Section("mysql").Key("user").String(),
		Password: cfg.Section("mysql").Key("password").String(),
		Database: cfg.Section("mysql").Key("database").String(),
	}, nil
}
func ShowInfo(config Config) {
	fmt.Println("=========== redis section ===========")
	fmt.Println("redis ip:", config.Redis.IP)
	fmt.Println("redis port:", config.Redis.Port)
	fmt.Println("=========== mysql section ===========")
	fmt.Println("mysql ip:", config.MySQL.IP)
	fmt.Println("mysql port:", config.MySQL.Port)
	fmt.Println("mysql user:", config.MySQL.User)
	fmt.Println("mysql password:", config.MySQL.Password)
	fmt.Println("mysql db:", config.MySQL.Database)
}

func test2() {

	container := dig.New()
	container.Provide(InitOption)
	container.Provide(InitConfig)
	container.Provide(InitRedisConfig)
	container.Provide(InitMySQLConfig)
	err := container.Invoke(ShowInfo) // 触发
	if err != nil {
		log.Fatal(err)
	}
}

func InitRedisAndMySQLConfig(cfg *ini.File) (Config, error) {
	// 手动构造了Config
	var config Config
	redis, err := InitRedisConfig(cfg)
	if err != nil {
		return config, err
	}
	mysql, err := InitMySQLConfig(cfg)
	if err != nil {
		return config, err
	}
	config.Redis = redis
	config.MySQL = mysql
	return config, nil
}

func PrintInfo3(redis *RedisConfig, mysql *MySQLConfig) {
	fmt.Println("=========== redis section ===========")
	fmt.Println("redis ip:", redis.IP)
	fmt.Println("redis port:", redis.Port)
	fmt.Println("=========== mysql section ===========")
	fmt.Println("mysql ip:", mysql.IP)
	fmt.Println("mysql port:", mysql.Port)
	fmt.Println("mysql user:", mysql.User)
	fmt.Println("mysql password:", mysql.Password)
	fmt.Println("mysql db:", mysql.Database)
}

func test3() {
	container := dig.New()
	container.Provide(InitOption)
	container.Provide(InitConfig)
	container.Provide(InitRedisAndMySQLConfig) // 该函数的output-Config中具有redis和mysql的依赖提供, 所以标记之后可以不需要再提供两者的构造实现
	err := container.Invoke(PrintInfo3)        // 参数要求是redisconfig和mysqlconfig
	if err != nil {
		log.Fatal(err)
	}
}

// -------------------------------------------------------------------------------

type Config4 struct {
	dig.In
	Redis *RedisConfig `optional:"true"` // 可以没有该参数的实现
	MySQL *MySQLConfig
}

func PrintInfo4(config Config4) {
	// 其中redis并没有注入依赖, 所以该值应该是nil
	if config.Redis == nil {
		fmt.Println("no redis config")
	}
}

func test4() {
	container := dig.New()
	container.Provide(InitOption)
	container.Provide(InitConfig)
	container.Provide(InitMySQLConfig)
	container.Invoke(PrintInfo4)
}

// ------------------------------------------------------------------------

type User struct {
	Name string
	Age  int
}

type UserParams struct {
	dig.In
	User1 *User `name:"dj"` // 在token中对多个User对象进行分辨
	User2 *User `name:"dj2"`
}

func NewUser(name string, age int) func() *User {
	return func() *User {
		return &User{name, age}
	}
}

func PrintInfo5(params UserParams) error {
	fmt.Println("User 1 ===========")
	fmt.Println("Name:", params.User1.Name) // 两个相同结构体的成员变量的调用
	fmt.Println("Age:", params.User1.Age)
	fmt.Println("User 2 ===========")
	fmt.Println("Name:", params.User2.Name)
	fmt.Println("Age:", params.User2.Age)
	return nil
}

func test5() {
	container := dig.New()
	container.Provide(NewUser("dj", 18), dig.Name("dj")) // 有名的对象, 相当于给出的构造方法自带参数, 更具象了
	container.Provide(NewUser("dj2", 18), dig.Name("dj2"))
	container.Invoke(PrintInfo5)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	test5()
}
