package structFile

type Car struct {  
    ID      string `json:"id"`  
    Brand   string `json:"brand"`  
    Model   string `json:"model"`  
    Mileage int    `json:"mileage"`  
    Owners  int    `json:"owners"`  
}  

type Furniture struct {  
    ID          string `json:"id"`  
    Name        string `json:"name"`  
    Manufacturer string `json:"manufacturer"`  
    Height      int    `json:"height"`  
    Width       int    `json:"width"`  
    Length      int    `json:"length"`  
}  

type Flower struct {
    ID          string `json:"id"`    
    Name        string `json:"name"`  
    Quantity    int    `json:"quantity"`  
    Price       int    `json:"price"`  
    ArrivalDate string `json:"arrival_date"`  
}  
type Inventory struct {  
    Cars      []Car     `json:"cars"`     
    Furniture []Furniture `json:"furniture"` 
    Flowers   []Flower   `json:"flowers"`             
}