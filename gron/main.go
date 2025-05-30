package main

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/roylee0704/gron"
	"github.com/roylee0704/gron/xtime"
)

// 设置定时任务
// 1. New一个定时器
// 2. Add添加定时任务
// 3. Start启动
func testEvery() {
	var wg sync.WaitGroup
	wg.Add(1)
	c := gron.New()

	// gron.Every设置定时间隔
	c.AddFunc(gron.Every(5*time.Second), func() {
		fmt.Println("runs every 5 seconds.")
	})
	c.Start()
	wg.Wait()
}

func testAt() {
	// 设置指定时间点
	var wg sync.WaitGroup
	wg.Add(1)
	c := gron.New()
	c.AddFunc(gron.Every(1*xtime.Day).At("15:50"), func() {
		fmt.Println("runs every second.")
	})
	c.Start()
	wg.Wait()
}

type GreetingJob struct {
	Name string
}

// Run -> Job接口
func (g GreetingJob) Run() {
	fmt.Println("Hello ", g.Name)
}

func testJob() {
	var wg sync.WaitGroup
	wg.Add(1)
	g1 := GreetingJob{Name: "dj"}
	g2 := GreetingJob{Name: "dajun"}
	c := gron.New()
	c.Add(gron.Every(5*time.Second), g1) // 自动执行其Run方法
	c.Add(gron.Every(10*time.Second), g2)
	c.Start()
	wg.Wait()
}

// 实现Schedule接口 - Next
type ExponentialBackOffSchedule struct {
	last int
}

func (e *ExponentialBackOffSchedule) Next(t time.Time) time.Time {
	interval := time.Duration(math.Pow(2.0, float64(e.last))) * time.Second
	e.last += 1
	return t.Truncate(time.Second).Add(interval)
}

// 指数退避
func testNext() {
	var wg sync.WaitGroup
	wg.Add(1)
	c := gron.New()
	c.AddFunc(&ExponentialBackOffSchedule{}, func() {
		fmt.Println(time.Now().Local().Format("2006-01-02 15:04:05"), "hello")
	})
	c.Start()
	wg.Wait()
}

type alarmSchedule struct {
	period time.Duration
	hh     int
	mm     int
}

// gron中的时间会强制转换为 UTC时区, 即使输入的t为CST时区, 这样就会导致时间计算有问题
func (as alarmSchedule) reset(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), as.hh, as.mm, 0, 0, t.Location())
}

func (as alarmSchedule) Next(t time.Time) time.Time {
	// t为now
	next := as.reset(t)

	fmt.Printf("t: %v, \n UTC: %v\n", t, t.UTC())
	fmt.Printf("next: %v, \n UTC: %v\n", next, next.UTC())
	fmt.Println("t.After(next):", t.After(next))

	// 判断t是否晚于next
	if (t.UTC()).After(next) {
		fmt.Println("next")
		return next.Add(as.period)
	}
	fmt.Println("early")
	return next
}

type alarm struct {
	period time.Duration
}

func parse(hhmm string) (hh int, mm int, err error) {

	hh = int(hhmm[0]-'0')*10 + int(hhmm[1]-'0')
	mm = int(hhmm[3]-'0')*10 + int(hhmm[4]-'0')

	if hh < 0 || hh > 24 {
		hh, mm = 0, 0
		err = errors.New("invalid hh format")
	}
	if mm < 0 || mm > 59 {
		hh, mm = 0, 0
		err = errors.New("invalid mm format")
	}

	return
}

func (a *alarm) At(t string) gron.Schedule {
	if a.period < xtime.Day {
		panic("period must be at least in days")
	}

	// parse t naively
	h, m, err := parse(t)

	if err != nil {
		panic(err.Error())
	}

	return &alarmSchedule{
		period: a.period,
		hh:     h,
		mm:     m,
	}
}
func testAlarm() {
	// 设置指定时间点
	var wg sync.WaitGroup
	wg.Add(1)
	c := gron.New()
	alarm := &alarm{period: 1 * xtime.Day}
	c.AddFunc(alarm.At("17:20"), func() {
		fmt.Println("runs every second.")
	})
	c.Start()
	wg.Wait()
}

func main() {
	// testJob()
	// testNext()
	testAlarm()
}
