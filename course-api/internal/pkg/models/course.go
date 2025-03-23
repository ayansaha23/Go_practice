package models

// no space between json and fields
type Course struct {
	Id         int      `json:"id"`
	Name       string   `json:"name"`
	Price      float64  `json:"price"`
	Technology []string `json:"technology"`
}

func (c *Course) IsEmpty() bool {

	return c.Name == "" || c.Price == 0 || c.Technology == nil
}
