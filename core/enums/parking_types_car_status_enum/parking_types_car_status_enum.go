package parking_types_car_status_enum

type ParkingTypesCarStatus struct {
	Data int
	Name string
	Display string
}

func Available() ParkingTypesCarStatus {
	return ParkingTypesCarStatus{
		Data: 1,
		Name: "Available",
		Display: "Còn chỗ",
	}
}

func AlmostOutStock() ParkingTypesCarStatus {
	return ParkingTypesCarStatus{
		Data: 2,
		Name: "AlmostOutStock",
		Display: "Gần hết chỗ",
	}
}

func OutStock() ParkingTypesCarStatus {
	return ParkingTypesCarStatus{
		Data: 3,
		Name: "OutStock",
		Display: "Hết chỗ",
	}
}

func Closed() ParkingTypesCarStatus {
	return ParkingTypesCarStatus{
		Data: 4,
		Name: "Đóng cửa",
	}
}

var arrayTypesMapDisplay = map[int]string{
	Available().Data: Available().Display,
	AlmostOutStock().Data: AlmostOutStock().Display,
	OutStock().Data: OutStock().Display,
	Closed().Data: Closed().Display,
}

func GetDisplayFromData(typeData int) string {
	return arrayTypesMapDisplay[typeData]
}