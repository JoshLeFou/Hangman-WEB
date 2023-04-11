package main

import (
	"fmt"
	"hangman-web/classic"
	"html/template"
	"net/http"
)

var temp *template.Template
var isstarted bool
var stats classic.Variables

func home(w http.ResponseWriter, r *http.Request) {
	temp = template.Must(template.ParseGlob("home.gohtml"))
	temp.Execute(w, "home.gohtml")
	if r.Method == "POST" { //si le joueur appuie sur le bouton pour relancer une partie
		http.Redirect(w, r, "/home", http.StatusSeeOther) //redirige vers la page index
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if stats.Win {
		http.Redirect(w, r, "/win", http.StatusSeeOther) //est ce que la partie est gagnée
	}
	if stats.Lose {
		http.Redirect(w, r, "/lost", http.StatusSeeOther) //est ce que la partie est perdue
	}
	temp = template.Must(template.ParseGlob("index.gohtml")) //instancie le fichier index.html
	if !isstarted {
		stats = classic.Start() //initialise les valeures des statistiques de la partie
		isstarted = true        //permet d'initialiser qu'une fois la partie
	}
	if r.Method == "POST" {
		stats.Currentword = r.FormValue("Word") //get le mot sur le site
		stats = classic.Check(stats)            //teste le mot rentré par le joueur
		var Hwordsent string
		for _, v := range stats.Hiddenword { //permet de rajouter un espace entre les "_" du mot caché
			Hwordsent += string(v) + " "
		}
		data := struct { //création d'un objet contenant une structure avec des variables
			Hword      string
			Wordtested [50]string
			Hp         int
			Hangman    [8]string
		}{
			Hword:      Hwordsent, //Word ==> string  word ==> contenue du
			Wordtested: stats.Tested,
			Hp:         stats.Hp,
			Hangman:    stats.Hangman,
		}
		if stats.Lose {
			http.Redirect(w, r, "/lost", http.StatusSeeOther) //redirrection avant affichage
		}
		if stats.Win {
			http.Redirect(w, r, "/win", http.StatusSeeOther)
		}
		temp.ExecuteTemplate(w, "index.gohtml", data)
	}
	if r.Method == "GET" {
		var Hwordsent string
		for _, v := range stats.Hiddenword {
			Hwordsent += string(v) + " "
		}
		data := struct { //création d'un objet contenant une structure avec des variables
			Hword      string
			Wordtested [50]string
			Hp         int
			Hangman    [8]string
		}{
			Hword:      Hwordsent, //Word ==> string  word ==> contenue du string
			Wordtested: stats.Tested,
			Hp:         stats.Hp,
			Hangman:    stats.Hangman,
		}
		temp.ExecuteTemplate(w, "index.gohtml", data)
	}
}

func lost(w http.ResponseWriter, r *http.Request) {
	if !stats.Lose {
		http.Redirect(w, r, "/", http.StatusSeeOther) //cas ou le joueur essaie d'accéder à la page
	} //de défaite via l'url avant d'avoir perdu
	temp = template.Must(template.ParseGlob("lost.gohtml"))
	var Hwordsent string
	for _, v := range stats.Hiddenword {
		Hwordsent += string(v) + " "
	}
	data := struct { //création d'un objet contenant une structure avec des variables
		Word       string
		Hword      string
		Wordtested [50]string
		Hp         int
		Hangman    [8]string
	}{
		Word:       stats.Realword,
		Hword:      Hwordsent,
		Wordtested: stats.Tested,
		Hp:         stats.Hp,
		Hangman:    stats.Hangman,
	}
	if r.Method == "POST" { //si le joueur appuie sur le bouton pour rejouer
		stats = classic.Start()                       //réinitialisation de la partie
		http.Redirect(w, r, "/", http.StatusSeeOther) //redirige vers la page de jeu
	}
	temp.ExecuteTemplate(w, "lost.gohtml", data)
}

func win(w http.ResponseWriter, r *http.Request) {
	if !stats.Win {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	temp = template.Must(template.ParseGlob("win.gohtml"))
	data := struct { //création d'un objet contenant une structure avec des variables
		Word       string
		Hword      string
		Wordtested [50]string
		Hp         int
		Hangman    [8]string
	}{
		Word:       stats.Realword,
		Hword:      stats.Hiddenword, //Word ==> string  word ==> contenue du string
		Wordtested: stats.Tested,
		Hp:         stats.Hp,
		Hangman:    stats.Hangman,
	}
	if r.Method == "POST" { //si le joueur appuie sur le bouton pour relancer une partie
		stats = classic.Start()                       //réinitialise la partie
		http.Redirect(w, r, "/", http.StatusSeeOther) //redirige vers la page de jeu
	}
	temp.ExecuteTemplate(w, "win.gohtml", data)
}

func main() {
	http.HandleFunc("/index", index) //préparer la page internet localhost:8080/
	http.HandleFunc("/", home)
	http.HandleFunc("/lost", lost)
	http.HandleFunc("/win", win)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	fmt.Println("http://localhost:8080/ - Server is on")
	http.ListenAndServe(":8080", nil) //définie le port 8080
}
