package endpoint

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/rush-travel/util"

	//"github.com/rs/xid"
	"github.com/rush-travel/config"
	"github.com/rush-travel/model"
)

type Coordinate struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}
type Coordinate1 struct {
	Latitude  float64 `json:"lati"`
	Longitude float64 `json:"long"`
}
type Order struct {
	ID             uint
	UserId         int            `gorm:"column:user_id" json:"user_id"`
	OrderId        string         `gorm:"column:order_id" json:"orderId"`
	Username       string         `gorm:"column:username" json:"username"`
	PickUpLoc      pq.StringArray `gorm:"type:string;column:pick_up_loc"`
	DropOffLoc     pq.StringArray `gorm:"type:string;column:drop_off_loc"`
	Date           string         `gorm:"column:date" json:"date"`
	Days           string         `gorm:"column:days" json:"days"`
	NoOfPassangers int            `gorm:"column:no_off_passangers" json:"no_off_passangers"`
	Price          int            `gorm:"column:price" json:"price"`
	VehicleType    string         `gorm:"column:vehicle_type" json:"vehicle_type"`
	Status         string         `json:"status"`
	Note           string         `json:"note"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
	// DeletedAt      *time.Time     `sql:"index"`
}

func CreateOrder(c *gin.Context) {
	var ord Order
	var buffer bytes.Buffer
	var buf bytes.Buffer
	//var coord []Coordinate
	//cor := Coordinate{}
	var pick []byte
	var drop []byte
	var err error
	//var coord1 []Coordinate1

	//value Data
	//pickUpLoc
	lat, _ := strconv.ParseFloat(c.PostForm("lat"), 64)
	lng, _ := strconv.ParseFloat(c.PostForm("lng"), 64)
	//dropOffLoc
	lati, _ := strconv.ParseFloat(c.PostForm("lati"), 64)
	long, _ := strconv.ParseFloat(c.PostForm("long"), 64)
	// Date, _ := time.Parse("2006-01-02", c.PostForm("Date"))
	NoOffPass, _ := strconv.Atoi(c.PostForm("no_off_passangers"))
	Userid, _ := strconv.Atoi(c.PostForm("UserID"))
	Price, _ := strconv.Atoi(c.PostForm("Price"))

	//
	pickUpLoc := map[string]interface{}{
		"Latitude":  lat,
		"Longitude": lng,
	}

	for _, pickUpLoc := range pickUpLoc {
		pick, err = json.Marshal(pickUpLoc)
		fmt.Println("pick", pick)
		if err != nil {
			fmt.Println(err)
		}
		buffer.WriteString(string(pick) + " ")
	}
	s := strings.TrimSpace(buffer.String())
	fmt.Println("s", s)
	ss := []string{s}
	// json.Unmarshal([]byte(s), &cor)
	fmt.Println("ss", ss)
	//
	dropOffLoc := map[string]interface{}{
		"Latitude":  lati,
		"Longitude": long,
	}
	fmt.Println(" d ", dropOffLoc)
	for _, dropOffLoc := range dropOffLoc {
		drop, err = json.Marshal(dropOffLoc)
		if err != nil {
			fmt.Println(err)
		}
		buf.WriteString(string(drop) + " ")
	}

	d := strings.TrimSpace(buf.String())
	dd := []string{d}

	fmt.Println("dd", d)
	//drop := []string{d}
	n := 5
	unique := make([]byte, n)
	if _, err := rand.Read(unique); err != nil {
		panic(err)
	}
	ID := fmt.Sprintf("%X", unique)
	//fmt.Println(ID)

	createorder := Order{
		OrderId:        ID,
		UserId:         Userid,
		Username:       c.PostForm("Username"),
		PickUpLoc:      ss,
		DropOffLoc:     dd,
		Date:           c.PostForm("Date"),
		Days:           c.PostForm("Days"),
		Price:          Price,
		NoOfPassangers: NoOffPass,
		VehicleType:    c.PostForm("VehicleType"),
		Note:           c.PostForm("Note"),
		Status:         "waiting",
		CreatedAt:      time.Now(),
	}
	err = config.DB.Model(&ord).Save(&createorder).Error
	if err != nil {
		fmt.Println("err", err)
		util.CallServerError(c, "failed", nil)
		return
	}
	fmt.Println("datanya: ", createorder)
	util.CallSuccessOK(c, "Successfully Add Order", createorder)
}

func FetchOrder(c *gin.Context) {
	var order []Order
	var order1 []Order

	config.DB.Model(&order).Find(&order)
	if len(order) < 0 {
		util.CallErrorNotFound(c, "No Order Data", nil)
		return
	}

	for _, item := range order {
		order1 = append(order1, Order{
			ID:             item.ID,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
			DeletedAt:      item.DeletedAt,
			OrderId:        item.OrderId,
			Username:       item.Username,
			PickUpLoc:      item.PickUpLoc,
			DropOffLoc:     item.DropOffLoc,
			Date:           item.Date,
			Days:           item.Days,
			Price:          item.Price,
			NoOfPassangers: item.NoOfPassangers,
			VehicleType:    item.VehicleType,
		})
	}
	util.CallSuccessOK(c, "Fetch All Orders Data ", order1)
}

func UpdateOrderStatus(c *gin.Context) {
	var order model.Order

	//request
	id := c.Param("id")
	orderID := c.Param("orderID")
	note := c.PostForm("note")
	status := c.PostForm("status")

	fmt.Println("notenya: ", note)
	fmt.Println("statusnya: ", status)

	// find user_id and orderid param can look at main.go
	config.DB.Model(&order).Where("user_id = ? AND order_id = ? ", id, orderID).Find(&order)

	update := model.Order{
		Note:   note,
		Status: status,
	}
	// disini pakai where().find(order)
	err := config.DB.Model(&order).Where("user_id = ?  AND order_id = ? ", id, orderID).Update(&update).Error
	if err != nil {
		util.CallServerError(c, "failed when to try update order status", nil)
	}
	util.CallSuccessOK(c, "Successfully Update Order Status", nil)
}

func GetOrderUser(c *gin.Context) {
	var (
		orders []Order
		order  []Order
	)

	// find user_id
	id := c.Param("id")
	// config.DB.Model(&order).Find(&order)
	config.DB.Model(&orders).Where("user_id = ? ", id).Find(&orders)

	if len(orders) <= 0 {
		util.CallErrorNotFound(c, "No Order Found!", nil)
		return
	}

	for _, item := range orders {
		order = append(order, Order{
			ID:             item.ID,
			UserId:         item.UserId,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
			DeletedAt:      item.DeletedAt,
			OrderId:        item.OrderId,
			Username:       item.Username,
			PickUpLoc:      item.PickUpLoc,
			DropOffLoc:     item.DropOffLoc,
			Date:           item.Date,
			Days:           item.Days,
			Price:          item.Price,
			Note:           item.Note,
			Status:         item.Status,
			NoOfPassangers: item.NoOfPassangers,
			VehicleType:    item.VehicleType,
		})
	}

	util.CallSuccessOK(c, "Fetch Order User Success ", order)

}

// haversin(θ) function get price from map
// source : https://gist.github.com/cdipaolo/d3f8db3848278b49db68
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func price(lat1, long1, lat2, long2, priceSeat float64) float64 {
	//must cast radius as float to multiply later
	var (
		la1, lo1, la2, lo2, r float64
	)

	//Pickup la1= 1.146528 lo1 =104.012343
	la1 = lat1 * math.Pi / 180
	lo1 = long1 * math.Pi / 180
	//Destination la2= 1.142307 lo2 = 104.011493
	la2 = lat2 * math.Pi / 180
	lo2 = long2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)
	resulth := 2 * r * math.Asin(math.Sqrt(h))
	//get 3 digit behind comma
	roundDown := math.Floor(resulth * 100 / 100)
	price := (roundDown * 65) + priceSeat

	fmt.Println("dist: ", roundDown)
	fmt.Println("priceORG: ", (roundDown * 65))

	return price
}

func OrderPrice(c *gin.Context) {

	person := c.PostForm("person")
	// pickup convert to float 64
	lat1, _ := strconv.ParseFloat(c.PostForm("lat1"), 64)
	long1, _ := strconv.ParseFloat(c.PostForm("long1"), 64)
	// destionation
	lat2, _ := strconv.ParseFloat(c.PostForm("lat2"), 64)
	long2, _ := strconv.ParseFloat(c.PostForm("long2"), 64)

	// time
	timer := time.Now()
	tFormat := timer.Format("15:04:05")
	//regex compile
	pagi, err := regexp.Compile("^(05|06|07|08|09)")
	siang, err := regexp.Compile("^(10|11|12|13|14)")
	sore, err := regexp.Compile("^(15|16|17)")
	petang, err := regexp.Compile("^(18)")
	malam, err := regexp.Compile("^(19|20|21|22|23)")
	//else
	// diniHari, err := regexp.Compile("^(24|01|02|03|04)")

	if err != nil {
		fmt.Println("Error regex:", err.Error())
	}

	if pagi.MatchString(tFormat) {
		fmt.Println("time PAGI: ", tFormat)
	} else if siang.MatchString(tFormat) {
		fmt.Println("time SIANG: ", tFormat)
	} else if sore.MatchString(tFormat) {
		fmt.Println("time SORE: ", tFormat)
	} else if petang.MatchString(tFormat) {
		fmt.Println("time PETANG: ", tFormat)
	} else if malam.MatchString(tFormat) {
		fmt.Println("time MALAM: ", tFormat)
	} else {
		fmt.Println("time Dini HARI: ", tFormat)
	}

	//price by weekend
	weekday := time.Now().Weekday().String()
	//get price
	switch person {
	case "6":
		if weekday == "Saturday" {
			if pagi.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PAGI: ", tFormat)
			} else if siang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SIANG: ", tFormat)
			} else if sore.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SORE: ", tFormat)
			} else if petang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PETANG: ", tFormat)
			} else if malam.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time MALAM: ", tFormat)
			} else {
				disc := 100000 * 0.2
				totalPrc := (100000 + 50000) - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time Dini HARI: ", tFormat)
			}

		} else if weekday == "Sunday" {
			if pagi.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PAGI: ", tFormat)
			} else if siang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SIANG: ", tFormat)
			} else if sore.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SORE: ", tFormat)
			} else if petang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PETANG: ", tFormat)
			} else if malam.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time MALAM: ", tFormat)
			} else {
				disc := 100000 * 0.2
				totalPrc := (100000 + 50000) - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time Dini HARI: ", tFormat)
			}
		} else {
			util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, 100000))
		}
		break
	case "13":
		if weekday == "Saturday" {
			if pagi.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PAGI: ", tFormat)
			} else if siang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SIANG: ", tFormat)
			} else if sore.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SORE: ", tFormat)
			} else if petang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PETANG: ", tFormat)
			} else if malam.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time MALAM: ", tFormat)
			} else {
				disc := 100000 * 0.2
				totalPrc := (100000 + 50000) - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time Dini HARI: ", tFormat)
			}

		} else if weekday == "Sunday" {
			if pagi.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PAGI: ", tFormat)
			} else if siang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SIANG: ", tFormat)
			} else if sore.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SORE: ", tFormat)
			} else if petang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PETANG: ", tFormat)
			} else if malam.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time MALAM: ", tFormat)
			} else {
				disc := 100000 * 0.2
				totalPrc := (100000 + 50000) - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time Dini HARI: ", tFormat)
			}
		} else {
			util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, 300000))
		}
		break
	case "25":
		if weekday == "Saturday" {
			if pagi.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PAGI: ", tFormat)
			} else if siang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SIANG: ", tFormat)
			} else if sore.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SORE: ", tFormat)
			} else if petang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PETANG: ", tFormat)
			} else if malam.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time MALAM: ", tFormat)
			} else {
				disc := 100000 * 0.2
				totalPrc := (100000 + 50000) - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time Dini HARI: ", tFormat)
			}

		} else if weekday == "Sunday" {
			if pagi.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PAGI: ", tFormat)
			} else if siang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SIANG: ", tFormat)
			} else if sore.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SORE: ", tFormat)
			} else if petang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PETANG: ", tFormat)
			} else if malam.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time MALAM: ", tFormat)
			} else {
				disc := 100000 * 0.2
				totalPrc := (100000 + 50000) - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time Dini HARI: ", tFormat)
			}
		} else {
			util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, 500000))
		}
		break
	case "45":
		if weekday == "Saturday" {
			if pagi.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PAGI: ", tFormat)
			} else if siang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SIANG: ", tFormat)
			} else if sore.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SORE: ", tFormat)
			} else if petang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PETANG: ", tFormat)
			} else if malam.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time MALAM: ", tFormat)
			} else {
				disc := 100000 * 0.2
				totalPrc := (100000 + 50000) - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time Dini HARI: ", tFormat)
			}

		} else if weekday == "Sunday" {
			if pagi.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PAGI: ", tFormat)
			} else if siang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SIANG: ", tFormat)
			} else if sore.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time SORE: ", tFormat)
			} else if petang.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time PETANG: ", tFormat)
			} else if malam.MatchString(tFormat) {
				disc := 100000 * 0.2
				totalPrc := 100000 - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time MALAM: ", tFormat)
			} else {
				disc := 100000 * 0.2
				totalPrc := (100000 + 50000) - disc
				util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, totalPrc))
				fmt.Println("time Dini HARI: ", tFormat)
			}
		} else {
			util.CallSuccessOK(c, "Fetch Price Success ", price(lat1, long1, lat2, long2, 650000))
		}
		break
	default:
		util.CallUserError(c, "Please Insert Passengers around 6,13,25,45", nil)
	}
}
