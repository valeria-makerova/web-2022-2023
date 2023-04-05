package main

import (
	"html/template"
	"log"
	"net/http"
)

type indexPage struct {
	FeaturedPosts   []featuredPostData
	MostRecentPosts []mostRecentPostData
}

type postPage struct{}

type featuredPostData struct {
	Title       string
	Subtitle    string
	Author      string
	AuthorImg   string
	PublishData string
	ImgModifier string
}

type mostRecentPostData struct {
	Title       string
	Subtitle    string
	Author      string
	AuthorImg   string
	PublishData string
	Image       string
}

func index(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/index.html") // Главная страница блога
	if err != nil {
		http.Error(w, "Internal Server Error", 500) // В случае ошибки парсинга - возвращаем 500
		log.Println(err.Error())                    // Используем стандартный логгер для вывода ошбики в консоль
		return                                      // Не забываем завершить выполнение ф-ии
	}

	data := indexPage{
		FeaturedPosts:   featuredPosts(),
		MostRecentPosts: mostRecentPost(),
	}

	err = ts.Execute(w, data) // Заставляем шаблонизатор вывести шаблон в тело ответа
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/post.html") // Главная страница блога
	if err != nil {
		http.Error(w, "Internal Server Error", 500) // В случае ошибки парсинга - возвращаем 500
		log.Println(err.Error())                    // Используем стандартный логгер для вывода ошбики в консоль
		return                                      // Не забываем завершить выполнение ф-ии
	}

	data := postPage{}

	err = ts.Execute(w, data) // Заставляем шаблонизатор вывести шаблон в тело ответа
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func featuredPosts() []featuredPostData {
	return []featuredPostData{
		{
			Title:       "The Road Ahead",
			Subtitle:    "The road ahead might be paved - it might not be.",
			Author:      "Mat Vogels",
			ImgModifier: "featured-post_first",
			PublishData: "September 25, 2015",
			AuthorImg:   "/static/images/mat_vogels.png",
		},
		{
			Title:       "From Top Down",
			Author:      "William Wong",
			Subtitle:    "Once a year.",
			ImgModifier: "featured-post_second",
			PublishData: "September 25, 2015",
			AuthorImg:   "/static/images/william_wong.png",
		},
	}
}

func mostRecentPost() []mostRecentPostData {
	return []mostRecentPostData{
		{
			Title:       "Still Standing Tall",
			Subtitle:    "Life begins at the end of your comfort zone.",
			Author:      "William Wong",
			AuthorImg:   "/static/images/william_wong.png",
			PublishData: "9/25/2015",
			Image:       "/static/images/paraschutes.jpg",
		},
	}
}
