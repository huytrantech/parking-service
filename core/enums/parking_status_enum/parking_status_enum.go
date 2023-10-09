package parking_status_enum

type ParkingStatusEnum struct {
	Data    int
	Name    string
	Display string
}

func Pending() ParkingStatusEnum {
	return ParkingStatusEnum{
		Data:    1,
		Name:    "Pending",
		Display: "Pending",
	}
}

func Active() ParkingStatusEnum {
	return ParkingStatusEnum{
		Data:    2,
		Name:    "Active",
		Display: "Active",
	}
}

func Block() ParkingStatusEnum {
	return ParkingStatusEnum{
		Data:    3,
		Name:    "Block",
		Display: "Block",
	}
}

func Deleted() ParkingStatusEnum {
	return ParkingStatusEnum{
		Data:    4,
		Name:    "Deleted",
		Display: "Deleted",
	}
}

func Closed() ParkingStatusEnum {
	return ParkingStatusEnum{
		Data:    5,
		Name:    "Closed",
		Display: "Closed",
	}
}

func Denied() ParkingStatusEnum {
	return ParkingStatusEnum{
		Data:    6,
		Name:    "Denied",
		Display: "Denied",
	}
}

func OutStock() ParkingStatusEnum {
	return ParkingStatusEnum{
		Data:    6,
		Name:    "OutStock",
		Display: "OutStock",
	}
}

var dataEnumMap = map[int]ParkingStatusEnum{
	Pending().Data:  Pending(),
	Active().Data:   Active(),
	Block().Data:    Block(),
	Deleted().Data:  Deleted(),
	Closed().Data:   Closed(),
	OutStock().Data: OutStock(),
}

func GetEnumFromData(data int) ParkingStatusEnum {
	return dataEnumMap[data]
}
