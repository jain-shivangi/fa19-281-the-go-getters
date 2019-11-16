/*
		Pizza ordering APIs
*/

package main

import (
	"os"
	"encoding/json"
    	"fmt"
    	"net/http"
    	"github.com/codegangsta/negroni"
    	"github.com/gorilla/handlers"
    	"github.com/gorilla/mux"
    	uuid "github.com/satori/go.uuid"
    	"github.com/unrolled/render"
    	"gopkg.in/mgo.v2"
    	"gopkg.in/mgo.v2/bson"
)

var server_mongo = os.Getenv("Server")
var database = os.Getenv("Database")
var collection = os.Getenv("Collection")
var user = os.Getenv("User")
var password = os.Getenv("Password")


func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
        allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"})
        allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	n.UseHandler(handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins)(mx))
	return n
}


// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/order/ping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/order/{orderId}", getPizzaByOrderIdHandler(formatter)).Methods("GET")
	mx.HandleFunc("/order", orderPizzaHandler(formatter)).Methods("POST")
}


// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Order API is alive"})
	}
}

//GET request to handle and return order details by order id
func getPizzaByOrderIdHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		//Setup
		session, err := mgo.Dial(server_mongo)
		if err != nil {
                formatter.JSON(w, http.StatusInternalServerError, "Error")
                return
         }
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)

		c := session.DB(database).C(collection)
		fmt.Println("Session established")
		parameter := mux.Vars(req)
		var pizza_order PizzaOrder
		var order_id string = parameter["orderId"]

		order_err := c.Find(bson.M{"orderId": order_id}).One(&pizza_order)
		if order_err != nil {
			formatter.JSON(w, http.StatusNotFound, "Please check the order id. This order id doesn't exist.")
			return
		}
		_ = json.NewDecoder(req.Body).Decode(&pizza_order)
		fmt.Println("Order description: ", pizza_order)
		formatter.JSON(w, http.StatusOK, pizza_order)
	}
}


//Post call to place order for pizza
func orderPizzaHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

        //connect to mongodb
		var order_pizza RequiredPayload
		_ = json.NewDecoder(req.Body).Decode(&order_pizza)
		session, error := mgo.Dial(server_mongo)
		if error = session.DB(database).Login(user, password); error != nil {
			formatter.JSON(w, http.StatusInternalServerError, "Error")
			return
		}
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
		c := session.DB(database).C(collection)

		var order PizzaOrder
		var item_list PizzaItem
		item_list.ItemId = order_pizza.ItemId
		item_list.ItemName = order_pizza.ItemName
		item_list.ItemPrice = order_pizza.ItemPrice
		item_list.ItemQuantity = order_pizza.ItemQuantity

		record_error := c.Find(bson.M{"userId": order_pizza.UserId, "orderStatus": "Processing"}).One(&order)

		if record_error == nil {
			order.Items = append(order.Items, item_list)
			order.TotalAmount = (order.TotalAmount + item_list.ItemPrice * item_list.ItemQuantity)
			c.Update(bson.M{"userId": order_pizza.UserId}, bson.M{"$set": bson.M{"items": order.Items, "totalAmount": order.TotalAmount}})
		    fmt.Println("Order already exists")
		} else {
			u, _ := uuid.NewV4()
		    order = PizzaOrder{
		        OrderId: u.String(),
				UserId:  order_pizza.UserId,
				Items: []PizzaItem{
					item_list},
				OrderStatus: "Active",
				TotalAmount: order_pizza.ItemPrice}
			error := c.Insert(order)
			if error != nil {
				formatter.JSON(w, http.StatusInternalServerError, "Sorry, order could not be placed.")
				return
			}
			fmt.Println("Congratulation! Order Placed.")
		}
		formatter.JSON(w, http.StatusOK, order)
	}
}
