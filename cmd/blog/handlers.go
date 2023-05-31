package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type indexPage struct {
	Title           string
	Subtitle        string
	FeaturedPosts   []featuredPostData
	MostRecentPosts []mostRecentPostData
}

type postContentData struct {
	Title    string `db:"title"`
	Subtitle string `db:"subtitle"`
	PostImg  string `db:"image_url"`
	Content  string `db:"content"`
}

type featuredPostData struct {
	PostID            string `db:"post_id"`
	Title             string `db:"title"`
	Subtitle          string `db:"subtitle"`
	Author            string `db:"author"`
	AuthorImg         string `db:"author_url"`
	PublishData       string `db:"publish_date"`
	ImageFeaturedPost string `db:"image_url_featured_post"`
}

type mostRecentPostData struct {
	PostID            string `db:"post_id"`
	Title             string `db:"title"`
	Subtitle          string `db:"subtitle"`
	Author            string `db:"author"`
	AuthorImg         string `db:"author_url"`
	PublishData       string `db:"publish_date"`
	Image             string `db:"image_url"`
	ImageFeaturedPost string `db:"image_url_featured_post"`
}

func index(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := featuredPosts(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500) // В случае ошибки парсинга - возвращаем 500
			log.Println(err)
			return // Не забываем завершить выполнение ф-ии
		}

		recentPosts, err := mostRecent(db)
		if err != nil {
			http.Error(w, "Internal Server Error", 500) // В случае ошибки парсинга - возвращаем 500
			log.Println(err)
			return // Не забываем завершить выполнение ф-ии
		}

		ts, err := template.ParseFiles("pages/index.html") // Главная страница блога
		if err != nil {
			http.Error(w, "Internal Server Error", 500) // В случае ошибки парсинга - возвращаем 500
			log.Println(err)
			return // Не забываем завершить выполнение ф-ии
		}

		data := indexPage{
			Title:           "Let's do it together.",
			Subtitle:        "We travel the world in search of stories. Come along for the ride.",
			FeaturedPosts:   posts,
			MostRecentPosts: recentPosts,
		}

		err = ts.Execute(w, data) // Заставляем шаблонизатор вывести шаблон в тело ответа
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		log.Println("Request completed successfully")
	}
}

func post(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := mux.Vars(r)["postID"] // Получаем postID в виде строки из параметров урла

		postID, err := strconv.Atoi(postIDStr) // Конвертируем строку postID в число
		if err != nil {
			http.Error(w, "Invalid order id", 403)
			log.Println(err)
			return
		}

		post, err := postByID(db, postID)
		if err != nil {
			if err == sql.ErrNoRows {
				// sql.ErrNoRows возвращается, когда в запросе к базе не было ничего найдено
				// В таком случае мы возвращем 404 (not found) и пишем в тело, что ордер не найден
				http.Error(w, "Post not found", 404)
				log.Println(err)
				return
			}

			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		ts, err := template.ParseFiles("pages/post.html") // Главная страница блога
		if err != nil {
			http.Error(w, "Internal Server Error", 500) // В случае ошибки парсинга - возвращаем 500
			log.Println(err.Error())                    // Используем стандартный логгер для вывода ошбики в консоль
			return                                      // Не забываем завершить выполнение ф-ии
		}

		err = ts.Execute(w, post) // Заставляем шаблонизатор вывести шаблон в тело ответа
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err)
			return
		}

		log.Println("Request completed successfully")
	}
}

func featuredPosts(db *sqlx.DB) ([]featuredPostData, error) {
	const query = `
	    SELECT
			post_id,
	        title,
	        subtitle,
	        author,
	        author_url,
	        publish_date,
			image_url_featured_post
	    FROM
	        post
	    WHERE featured = 1
	` //Составляем SQL-запрос для получения записей для секции featured-posts

	var posts []featuredPostData // Заранее объявляем массив с результирующей информацией

	err := db.Select(&posts, query) // Делаем запрос в базу данных
	if err != nil {                 // Проверяем, что запрос в базу данных не завершился с ошибкой
		return nil, err
	}

	return posts, nil
}

func mostRecent(db *sqlx.DB) ([]mostRecentPostData, error) {
	const query = `
	 	SELECT
			post_id,
	  		title,
	  		subtitle,
	  		author,
	  		author_url,
	  		publish_date,
	  		image_url
	 	FROM
	  		post
	 	WHERE featured = 0
	` // Составляем SQL-запрос для получения записей для секции featured-posts

	var posts []mostRecentPostData // Заранее объявляем массив с результирующей информацией

	err := db.Select(&posts, query) // Делаем запрос в базу данных
	if err != nil {                 // Проверяем, что запрос в базу данных не завершился с ошибкой
		return nil, err
	}

	return posts, nil
}

func postByID(db *sqlx.DB, postID int) (postContentData, error) {
	const query = `
		SELECT
			title,
			subtitle,
			image_url,
			content
		FROM
			post
		WHERE
			post_id = ?
	`
	// В SQL-запросе добавились параметры, как в шаблоне. ? означает параметр, который мы передаем в запрос ниже

	var post postContentData

	// Обязательно нужно передать в параметрах postID
	err := db.Get(&post, query, postID)
	if err != nil {
		return postContentData{}, err
	}

	return post, nil
}
