package main

import ("fmt"
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"time"
	"encoding/json"
	"github.com/go-redis/redis"
)

func main(){
	fmt.Println("Starting rest service")

	r := mux.NewRouter()

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	Init(r)

	log.Fatal(srv.ListenAndServe())
}


var client *redis.Client

func Init(r *mux.Router) {
	client = NewDb()

	r .HandleFunc("/order", GetOrders).Methods("GET")
	r .HandleFunc("/order/{name}", GetOrder).Methods("GET")
	r .HandleFunc("/order/{name}", PlaceOrder).Methods("POST")
	r .HandleFunc("/order/{name}", DeleteOrder).Methods("DELETE")
}

type Order struct {
	Name string   `json:"name,omitempty"`
	Date time.Time `json:"date,omitempty"`
	Coffee string `json:"coffee,omitempty"`
}

var orders[] Order

func GetOrders(w http.ResponseWriter, r *http.Request) {
	val, err := client.Keys("*").Result()

	if err != nil {
		panic(err)
	}

	for _,key:=range val {
		res,err := client.Get(key).Result()

		if err != nil{
			panic(err)
		}
		var order Order
		r := json.Unmarshal([]byte(res),&order)

		if r != nil {
			panic(err)
		}

		orders = append(orders,order)
	}

	json.NewEncoder(w).Encode(orders)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	val, err := client.Get(params["name"]).Result()
	if err != nil {
		json.NewEncoder(w).Encode("No record found at index " + params["name"])
		panic(err)
	}else{
		json.NewEncoder(w).Encode(val)
	}
}

func PlaceOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var order Order
	json.NewDecoder(r.Body).Decode(&order)
	order.Name = params["name"]
	order.Date = time.Now()

	err := client.Set(order.Name, order,0).Err()
	if err != nil {
		json.NewEncoder(w).Encode("Error placing order for " + order.Name)
		panic(err)
	}

	json.NewEncoder(w).Encode("Order place for " + order.Name)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	client.Del(params["name"])

	json.NewEncoder(w).Encode("Deleted: "+params["name"])
}

func (order Order) MarshalBinary() ([]byte, error) {
	return json.Marshal(order)
}

func NewDb() *redis.Client{
	client := redis.NewClient(&redis.Options{
		Addr:         ":6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
	client.FlushDB()
	return client
}