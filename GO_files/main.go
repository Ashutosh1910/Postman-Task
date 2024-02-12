package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/xuri/excelize/v2"
)

type Meal struct{
	day string
	items []string
	meal_name string
	date string
}

func get_index(item_list []string,name string)(int){
	for index,item:=range item_list{
		if item==name{ return index}

	}
	return len(item_list)
}

func printMeal(meal Meal){
	fmt.Printf("MEAL DATE: %v\n",meal.date)
	fmt.Printf("MEAL DAY: %v\n",meal.day)
	fmt.Printf("MEAL Name: %v\n",meal.meal_name)
	fmt.Printf("Menu:")
	for _,item:=range meal.items{
		if item!=""{
			fmt.Printf("%v,",item)}}
	fmt.Printf("\n=====================\n")
}

func check_input(to_check string,list []string) bool{
	for _,item:=range list{
		if to_check==item{
			return true
		}}
		return false
}

func main() {
 log.SetFlags(0)
 f,_:=excelize.OpenFile("Sample-Menu.xlsx")
 day_list:=[]string{"MONDAY","TUESDAY","WEDNESDAY","THURSDAY","FRIDAY","SATURDAY","SUNDAY"}
 meal_list:=[]string{"BREAKFAST","LUNCH","DINNER"}
 fmt.Println("1:View Meal")
 fmt.Println("2:Get no. of items in a meal")
 fmt.Println("3:Check item in meal")
 fmt.Println("4:Convert menu to json file")
 fmt.Println("5:Print instances of meals")

 fmt.Println("Enter your choice")
 var c int
 fmt.Scanf("%v\n",&c)
 switch c {
 case 1:{
		fmt.Print("Enter Day and Meal-")
		var day,meal string
		fmt.Scanf("%v %v",&day,&meal)
		
		valid_input:=check_input(day,day_list)&&check_input(meal,meal_list)
		if !valid_input{log.Fatalln("INVALID INPUT")}
		item_list,_:=View_Meal(day,meal,f)
		for _,item:=range item_list{
		fmt.Println(item)}}	
 case 2:{
		fmt.Print("Enter Day and Meal-")
		var day,meal string
		fmt.Scanf("%v %v",&day,&meal)
		valid_input:=check_input(day,day_list)&&check_input(meal,meal_list)
		if !valid_input{log.Fatalln("INVALID INPUT")}
		no_of_item:=No_of_items_in_meal(day,meal,f)
		fmt.Println(no_of_item)}
 case 3:{
		fmt.Println("Enter Day and Meal and Item-")
		var day,meal,i1,i2 string
		fmt.Scanf("%v %v %v %v",&day,&meal,&i1,&i2)
		valid_input:=check_input(day,day_list)&&check_input(meal,meal_list)
		if !valid_input{log.Fatalln("INVALID INPUT")}
		exists:=check_item(day,meal,i1+i2,f)
		if exists {
			fmt.Println("Item found in Meal")
		} else{
			fmt.Println("Item not found in Meal")}
 }
 case 4:
	{
		_=convert_to_json(f,false)
	}
 case 5:{
		create_instances(f)
}
}
}

func View_Meal(day string,meal string, f *excelize.File) ([]string,error){
 full_sheet,err:=f.GetCols("Sheet1")
  item_list:=make([]string,40,40)
  start_index:=0
  end_index:=0
 for _,col:= range full_sheet{
	if col[0]==day {
	 item_list=col
	 start_index=get_index(item_list,meal)+1
	 end_index=get_index(item_list[start_index:],day)+start_index
		break}}
 if err!=nil{ log.Fatalln(err)}
 return item_list[start_index:end_index],nil
}
func No_of_items_in_meal(day string,meal string,f *excelize.File) (int){
		no_of_items:=0
		item_list,_:=View_Meal(day,meal,f)
		no_of_items=len(item_list)
		return no_of_items
}
 
func check_item(day string ,meal string,item_name string,f *excelize.File)(bool){
		exists:=false
		item_list,_:=View_Meal(day,meal,f)
		for _,item:=range item_list{
		item := strings.ReplaceAll(item, " ", "")
		if item==item_name{
			exists=true
		}
	}
	return exists
}

func convert_to_json(f *excelize.File,get_data bool) ([]byte){
 filename:="mess-menu.json"
 json_file,err:=os.Create(filename)
 if err!=nil{
	log.Fatalln(err)
    }
 full_sheet,err:=f.GetCols("Sheet1")
 menu:=make(map[string](map[string][]string))
 for _,col:=range full_sheet{
		
		date:=col[1]
		menu[date]=map[string][]string{}
		breakfast_list,_:=View_Meal(col[0],"BREAKFAST",f)
		day:=[]string{col[0]}
		menu[date]["Day"]=day
		menu[date]["Breakfast"]=breakfast_list
		lunch_list,_:=View_Meal(col[0],"LUNCH",f)
		menu[date]["Lunch"]=lunch_list
		dinner_list,_:=View_Meal(col[0],"DINNER",f)
		
		menu[date]["dinner"]=dinner_list
	}  
	json_data,err:=json.MarshalIndent(menu,"","    ")
	if err != nil {
	log.Fatal(err)
    }
	json_file.Write(json_data)
	defer json_file.Close()
	
	if !get_data {
		fmt.Println(string(json_data))
    }
	return json_data
}

func create_instances(f *excelize.File){

	json_menu:=convert_to_json(f,true)
	menu:=make(map[string](map[string][]string))
	err:=json.Unmarshal(json_menu,&menu)
	if err!=nil{
		log.Fatalln(err)
	}
	var meal_list []Meal
	for key,value:= range menu {

		for meal,item_list:=range value{
			if !check_input(strings.ToUpper(meal),[]string{"BREAKFAST","LUNCH","DINNER"}){continue}
			var m Meal
			m.date=key
			m.meal_name=meal
			m.items=item_list
		    m.day=menu[key]["Day"][0]
			meal_list=append(meal_list, m)
			printMeal(m)


	 	}
	}
}
			
			
				
				
			

	


	
 


