package model

type Province struct {
	ProvinceID   string `json:"province_id"`
	ProvinceName string `json:"province_name"`
}

type District struct {
	ProvinceID   string `json:"province_id"`
	DistrictID   string `json:"district_id"`
	DistrictName string `json:"district_name"`
}

type SubDistrict struct {
	DistrictID      string `json:"district_id"`
	SubDistrictID   string `json:"sub_district_id"`
	SubDistrictName string `json:"sub_district_name"`
}

type Postcode struct {
	Postcode string `json:"postcode"`
}

type CarBrand struct {
	CarBrandID   string `json:"car_brand_id"`
	CarBrandName string `json:"car_brand_name"`
}

type CarModel struct {
	CarModelID   string `json:"car_model_id"`
	CarModelName string `json:"car_model_name"`
}

type CarYear struct {
	CarYearID string `json:"car_year_id"`
	CarYear   string `json:"car_year"`
}

type CarColor struct {
	CarColorID string `json:"car_color_id"`
	CarColor   string `json:"car_color"`
}
