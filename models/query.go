package models

import "redi/constants"

type URLIDQuery struct {
	URLID string `query:"url_id"`
}

type PageQuery struct {
	Limit  int             `query:"limit" default:"20"`
	Offset int             `query:"offset" default:"0"`
	Order  constants.Order `query:"order" default:"asc" validate:"order"`
}
