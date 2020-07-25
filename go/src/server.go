package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"runtime/debug"
	"sort"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	jwtware "github.com/gofiber/jwt"
	"github.com/gofiber/utils"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	child bool
	db    *pgxpool.Pool
)

const (
	jwtKey = "secret_jwt_key"
)

func main() {
	initDatabase()

	app := fiber.New(&fiber.Settings{
		CaseSensitive: true,
		StrictRouting: true,
	})

	if utils.GetArgument("-prefork") {
		app.Settings.Prefork = true
	}
	if utils.GetArgument("-prefork-child") {
		child = true
	}
	if utils.GetArgument("-nogc") {
		debug.SetGCPercent(-1)
	}

	app.Get("/json", jsonHandler)
	app.Get("/login", loginHandler)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey:    []byte(jwtKey),
		SigningMethod: "HS512",
	}))

	app.Get("/jwt", jwtHandler)
	app.Get("/flight", flightHandler)
	app.Get("/flights", flightsHandler)
	app.Get("/flights_sorted", flightsSortedHandler)

	app.Listen(3002)
}

// initDatabase :
func initDatabase() {
	maxConn := runtime.NumCPU()
	if maxConn == 0 {
		maxConn = 8
	}
	if child {
		maxConn = maxConn
	} else {
		maxConn = maxConn * 4
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s pool_max_conns=%d", "postgres", 5432, "postgres", "postgres", "demo", maxConn)
	var err error
	db, err = pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		panic(err)
	}
}

func jsonHandler(c *fiber.Ctx) {
	c.JSON(fiber.Map{"hello": "world"})
}

func jwtHandler(c *fiber.Ctx) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	c.JSON(claims)
}

func loginHandler(c *fiber.Ctx) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = fiber.Map{"id": 101, "first_name": "John", "last_name": "Doe"}
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, _ := token.SignedString([]byte(jwtKey))
	c.JSON(fiber.Map{"token": t})
}

type JsonNullTime struct {
	sql.NullTime
}

func (v JsonNullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	}
	return json.Marshal(nil)
}

type Flight struct {
	ID                 uint32       `json:"flight_id"`
	FlightNo           string       `json:"flight_no"`
	ScheduledDeparture JsonNullTime `json:"scheduled_departure"`
	ScheduledArrival   JsonNullTime `json:"scheduled_arrival"`
	DepartureAirport   string       `json:"departure_airport"`
	ArrivalAirport     string       `json:"arrival_airport"`
	Status             string       `json:"status"`
	AircraftCode       string       `json:"aircraft_code"`
	ActualDeparture    JsonNullTime `json:"actual_departure"`
	ActualArrival      JsonNullTime `json:"actual_arrival"`
}

func flightHandler(c *fiber.Ctx) {
	data := Flight{}
	row := db.QueryRow(context.Background(), "SELECT flight_id, flight_no, scheduled_departure, scheduled_arrival, departure_airport, arrival_airport, status, aircraft_code, actual_departure, actual_arrival FROM bookings.flights WHERE aircraft_code = $1 ORDER BY flight_id LIMIT 1", "321")
	row.Scan(&data.ID, &data.FlightNo, &data.ScheduledDeparture, &data.ScheduledArrival, &data.DepartureAirport, &data.ArrivalAirport, &data.Status, &data.AircraftCode, &data.ActualDeparture, &data.ActualArrival)
	c.JSON(&data)
}

func flightsHandler(c *fiber.Ctx) {
	rows, err := db.Query(context.Background(), "SELECT flight_id, flight_no, scheduled_departure, scheduled_arrival, departure_airport, arrival_airport, status, aircraft_code, actual_departure, actual_arrival FROM bookings.flights WHERE aircraft_code = $1 ORDER BY flight_id LIMIT 500", "321")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	data := make([]Flight, 0)

	for rows.Next() {
		row := Flight{}
		err = rows.Scan(&row.ID, &row.FlightNo, &row.ScheduledDeparture, &row.ScheduledArrival, &row.DepartureAirport, &row.ArrivalAirport, &row.Status, &row.AircraftCode, &row.ActualDeparture, &row.ActualArrival)
		if err != nil {
			panic(err)
		}
		data = append(data, row)
	}
	c.JSON(&data)
}

func flightsSortedHandler(c *fiber.Ctx) {
	rows, err := db.Query(context.Background(), "SELECT flight_id, flight_no, scheduled_departure, scheduled_arrival, departure_airport, arrival_airport, status, aircraft_code, actual_departure, actual_arrival FROM bookings.flights WHERE aircraft_code = $1 ORDER BY flight_id LIMIT 500", "321")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	data := make([]Flight, 0)

	for rows.Next() {
		row := Flight{}
		err = rows.Scan(&row.ID, &row.FlightNo, &row.ScheduledDeparture, &row.ScheduledArrival, &row.DepartureAirport, &row.ArrivalAirport, &row.Status, &row.AircraftCode, &row.ActualDeparture, &row.ActualArrival)
		if err != nil {
			panic(err)
		}
		data = append(data, row)
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].ScheduledDeparture.Time.After(data[j].ScheduledDeparture.Time)
	})

	c.JSON(&data)
}
