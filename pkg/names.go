package gont

import (
	"math/rand"
)

// From: https://en.wikipedia.org/wiki/List_of_Internet_pioneers
//  and  https://www.internethalloffame.org/inductees/all

var (
	Names = []string{
		"akkerhuis",
		"akplogan",
		"allman",
		"andreessen",
		"andres",
		"armour-polly",
		"baker",
		"banks",
		"baran",
		"barlow",
		"berners-lee",
		"bina",
		"brandenburg",
		"bukhalid",
		"bush",
		"cailliau",
		"cerf",
		"chon",
		"cioffi",
		"claffy",
		"clark",
		"cohen",
		"comer",
		"crocker",
		"dalal",
		"davies",
		"dias",
		"elgamal",
		"emtage",
		"engelbart",
		"esterhuysen",
		"estrada",
		"farber",
		"feinler",
		"floyd",
		"flueckiger",
		"frank",
		"fuchs",
		"gerich",
		"getschko",
		"goldstein",
		"gore",
		"goto",
		"hafkin",
		"hagen",
		"heart",
		"herzfeld",
		"hirabaru",
		"holz",
		"hu",
		"huizer",
		"huston",
		"huter",
		"induruwa",
		"irving",
		"ishida",
		"jacobson",
		"jennings",
		"jensen",
		"kahle",
		"kahn",
		"kanchanasut",
		"karrenberg",
		"kent",
		"kirstein",
		"kleinrock",
		"klensin",
		"krol",
		"landweber",
		"laquey-parker",
		"leiner",
		"licklider",
		"loewinder",
		"lynch",
		"mccahill",
		"metcalfe",
		"mills",
		"mockapetris",
		"murai",
		"muthoni",
		"neggers",
		"nelson",
		"newmark",
		"nordhagen",
		"partridge",
		"pellow",
		"perlman",
		"pietrosemoli",
		"postel",
		"pouzin",
		"pun",
		"qian",
		"quaynor",
		"ramani",
		"reynolds",
		"ricart",
		"roberts",
		"sadowsky",
		"schulzrinne",
		"segal",
		"shannon",
		"soriano",
		"stallman",
		"stanton",
		"swartz",
		"takahashi",
		"taylor",
		"tin-wee",
		"tomlinson",
		"torvalds",
		"utreras",
		"van-houweling",
		"vixie",
		"wales",
		"wierenga",
		"wolff",
		"wu",
		"yamaguchi",
		"zimmermann",
		"zorn",
	}
)

func GetRandomName() string {
	index := rand.Intn(len(Names))
	return Names[index]
}
