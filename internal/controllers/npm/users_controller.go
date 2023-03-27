package npm

import (
	"fmt"
	"net/http"

	"github.com/deweppro/go-sdk/random"
	"github.com/deweppro/goppy/plugins/web"
)

func (v *Controller) UserLogin(c web.Context) {
	requestModel := LoginRequest{}
	err := c.BindJSON(&requestModel)
	if err != nil {
		c.Error(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(err, requestModel)

	responseModel := LoginResponse{
		Token: random.String(36),
		Ok:    true,
		ID:    requestModel.ID,
		Rev:   "",
	}
	c.JSON(http.StatusCreated, &responseModel)
	fmt.Println("NEW TOKEN: ", responseModel.Token)
}

func (v *Controller) DeleteToken(c web.Context) {
	token, err := c.Param("token").String()
	if err != nil {
		c.Error(http.StatusInternalServerError, err)
		return
	}
	fmt.Println("DEL TOKEN: ", token)
}
