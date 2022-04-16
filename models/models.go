package models

import "html/template"

type Player struct {
	Enum        int
	FName       string
	LName       string
	Profile     string
	Rating      string
	GamesPlayed string
}

type HomeData struct {
	SrcOne  template.URL
	EnumOne int
	SrcTwo  template.URL
	EnumTwo int
}

type RankingData struct {
	Sources []template.URL
}
