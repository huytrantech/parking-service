package parking_type_car_enum

type ParkingTypeCarEnum struct {
	Data int
	Name string
	Logo string
}

func Car() ParkingTypeCarEnum {
	return ParkingTypeCarEnum{
		Data: 1,
		Name: "Car",
		Logo: "",
	}
}

func Motorbike() ParkingTypeCarEnum {
	return ParkingTypeCarEnum{
		Data: 2,
		Name: "Motorbike",
		Logo: "",
	}
}

var ArrayTypesMap = map[int]ParkingTypeCarEnum{
	Car().Data:       Car(),
	Motorbike().Data: Motorbike(),
}

func GetLogoFromData(typeData int) string {
	return ArrayTypesMap[typeData].Logo
}

func GetTypeNameFromData(typeData int) string {
	return ArrayTypesMap[typeData].Name
}
