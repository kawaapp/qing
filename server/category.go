package server

import (
	"github.com/labstack/echo"
	"github.com/kawaapp/kawaqing/store"
	"fmt"
	"net/http"
	"strconv"
	"github.com/kawaapp/kawaqing/model"
)

func GetCategoryList(c echo.Context) error {
	arr, err := store.FromContext(c).GetCategoryList()
	if err != nil {
		return fmt.Errorf("GetCategoryList err, %v", err)
	}
	return jsonResp(c, 0, arr)
}

func GetCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	ctg, err := store.FromContext(c).GetCategory(int64(id))
	if err != nil {
		return fmt.Errorf("GetCategory err, %v", err)
	}
	return jsonResp(c, 0, ctg)
}

func CreateCategory(c echo.Context) error {
	ctg := new(model.Category)
	if err := c.Bind(ctg); err != nil {
		return c.String(400, err.Error())
	}
	err := store.FromContext(c).CreateCategory(ctg)
	if err != nil {
		return fmt.Errorf("CreateCategory err, %v", err)
	}
	return jsonResp(c, 0, ctg)
}

func UpdateCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	db := store.FromContext(c)
	ctg, err := db.GetCategory(int64(id))
	if err != nil {
		return fmt.Errorf("UpdateCategory, %v", err)
	}

	m := echo.Map{}
	if err := c.Bind(&m); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	in := (map[string]interface{})(m)

	// partial update
	if v, ok := in["name"]; ok {
		ctg.Name = v.(string)
	}
	if v, ok := in["summary"]; ok {
		ctg.Summary = v.(string)
	}
	if v, ok := in["sort"]; ok {
		ctg.Sort = int(v.(float64))
	}
	if v, ok := in["parent_id"]; ok {
		ctg.ParentId = int64(v.(float64))
	}
	if err := db.UpdateCategory(ctg); err != nil {
		return fmt.Errorf("UpdateCategory err, %v", err)
	}
	return jsonResp(c, 0, ctg)
}

func DeleteCategory(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = store.FromContext(c).DeleteCategory(int64(id))
	if err != nil {
		return fmt.Errorf("DeleteCategory err, %v", err)
	}
	return c.NoContent(200)
}